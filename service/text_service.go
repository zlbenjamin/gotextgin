package service

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/zlbenjamin/gotextgin/pkg"
	sttext "github.com/zlbenjamin/gotextgin/pkg/text"
	"github.com/zlbenjamin/gotextgin/service/database"
	"gorm.io/gorm"
)

// Add a text
func AddText(c *gin.Context) {
	var params AddTextParams
	if err := c.ShouldBind(&params); err != nil {
		resp := pkg.ApiResponse{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}
		c.JSON(http.StatusOK, resp)
		return
	}

	// trim
	params.Content = strings.Trim(params.Content, " ")
	params.Type = strings.Trim(params.Type, " ")
	for i, item := range params.Tags {
		params.Tags[i] = strings.Trim(item, " ")
	}

	var retId int32

	// operations in a transaction
	db := database.GetDB()
	db.Transaction(func(tx *gorm.DB) error {
		// 1/2 add text
		record := &sttext.Text{
			Content: params.Content,
			Type:    params.Type,
		}
		result := tx.Create(record)
		if result.Error != nil {
			log.Panicln("Add text failed：", result.Error.Error())
			return result.Error
		}
		log.Println("Add text success. id=", record.Id)

		// 2/2 add tags
		addTagsForText(tx, record.Id, params.Tags)

		retId = record.Id
		return nil
	})

	// Success.
	resp := pkg.ApiResponse{
		Code:    200,
		Message: "OK",
		Data:    retId,
	}
	c.JSON(http.StatusOK, resp)
}

// func addTagsForText(id int32, tags []string) {
func addTagsForText(tx *gorm.DB, id int32, tags []string) error {
	if len(tags) < 1 {
		return nil
	}

	var records []*sttext.TextTag
	for _, item := range tags {
		records = append(records, &sttext.TextTag{
			TextId: id,
			Name:   item,
		})
	}

	// db := database.GetDB()
	// result := db.Create(records)
	result := tx.Create(records)
	if result.Error != nil {
		log.Panicln("Add tags for text failed, id=", id, ", tags=", tags, "err=", result.Error.Error())
		return result.Error
	}

	log.Println("Add tags for text success, id=", id)
	return nil
}

// Get a text by id
func GetTextById(c *gin.Context) {
	var dto GetTextDTO
	if err := c.ShouldBindUri(&dto); err != nil {
		resp := pkg.ApiResponse{
			Code:    400,
			Message: "Invalid params: " + err.Error(),
		}
		c.JSON(http.StatusOK, resp)
		return
	}
	idi := dto.Id

	// return data
	var textData TextFullVO

	// 1/3 text
	var record sttext.Text
	db := database.GetDB()
	result := db.First(&record, idi)
	if errors.Is(result.Error, sql.ErrNoRows) {
		resp := pkg.ApiResponse{
			Code:    400,
			Message: "No data for id",
		}
		c.JSON(http.StatusOK, resp)
		return
	} else if result.Error != nil {
		log.Println("Query failed: id=", idi, "err=", result.Error)
		resp := pkg.ApiResponse{
			Code:    500,
			Message: "Query text failed",
		}
		c.JSON(http.StatusOK, resp)
		return
	}

	textData.Text = record

	// 2/3 tags
	result = db.Where("text_id = ?", idi).Order("create_time ASC").Find(&textData.Tags)
	if errors.Is(result.Error, sql.ErrNoRows) {
		textData.Tags = make([]sttext.TextTag, 0)
	} else if result.Error != nil {
		log.Println("Query tags failed: id=", idi, "err=", result.Error)
		resp := pkg.ApiResponse{
			Code:    500,
			Message: "Query tags failed",
		}
		c.JSON(http.StatusOK, resp)
		return
	}

	// 3/3 comments
	result = db.Where("text_id = ?", idi).Order("create_time DESC").Find(&textData.Comments)
	if errors.Is(result.Error, sql.ErrNoRows) {
		textData.Comments = make([]sttext.TextComment, 0)
	} else if result.Error != nil {
		log.Println("Query comments failed: id=", idi, "err=", result.Error)
		resp := pkg.ApiResponse{
			Code:    500,
			Message: "Query comments failed",
		}
		c.JSON(http.StatusOK, resp)
		return
	}

	resp := pkg.ApiResponse{
		Code:    200,
		Message: "OK",
		Data:    textData,
	}
	c.JSON(http.StatusOK, resp)
}

// Delete a text by id
func DeleteTextById(c *gin.Context) {
	var dto TextIdDTO
	if err := c.ShouldBindUri(&dto); err != nil {
		resp := pkg.ApiResponse{
			Code:    400,
			Message: "Invalid params: " + err.Error(),
		}
		c.JSON(http.StatusOK, resp)
		return
	}
	idi := dto.Id

	// Exec in a transaction
	// delete text, tags, comments
	rowsAffecteds := struct {
		R1 int64
		R2 int64
		R3 int64
	}{}
	db := database.GetDB()
	err := db.Transaction(func(tx *gorm.DB) error {
		result := tx.Delete(&sttext.Text{}, idi)
		if result.Error != nil {
			log.Printf("delete text failed: id=%d, err=%s\n", idi, result.Error.Error())
			return result.Error
		}
		if result.RowsAffected < 1 {
			// record has been deleted.
			return nil
		}
		rowsAffecteds.R1 = result.RowsAffected

		result = tx.Where("text_id = ?", idi).Delete(&sttext.TextTag{})
		if result.Error != nil {
			log.Printf("delete text failed#tags: id=%d, err=%s\n", idi, result.Error.Error())
			return result.Error
		}
		rowsAffecteds.R2 = result.RowsAffected

		result = tx.Where("text_id = ?", idi).Delete(&sttext.TextComment{})
		if result.Error != nil {
			log.Printf("delete text failed#comments: id=%d, err=%s\n", idi, result.Error.Error())
			return result.Error
		}
		rowsAffecteds.R3 = result.RowsAffected

		// stop delete type#1
		// log.Panicln("start a error...")

		// stop delete type#2
		// if 1 == 1 {
		// 	return errors.New("mock a error")
		// }

		// success
		// return nil will commit the whole transaction
		return nil
	})

	var resp pkg.ApiResponse
	if err == nil {
		if rowsAffecteds.R1 > 0 {
			log.Printf("Delete text success. id=%d, deleted rowsAffecteds=[%d, %d, %d]",
				idi, rowsAffecteds.R1, rowsAffecteds.R2, rowsAffecteds.R3)
		}
		resp = pkg.ApiResponse{
			Code:    200,
			Message: "OK",
			Data:    true,
		}
	} else {
		log.Printf("Delete text failed. id=%d, err=%s\n", idi, err.Error())
		resp = pkg.ApiResponse{
			Code:    500,
			Message: "DELETE FAILED",
			Data:    false,
		}
	}

	c.JSON(http.StatusOK, resp)
}

// Page find texts.
// ORDER BY create_time DESC
func PageFindText(c *gin.Context) {
	var params pageFindParams
	if err := c.ShouldBind(&params); err != nil {
		resp := pkg.ApiResponse{
			Code:    400,
			Message: err.Error(),
		}
		c.JSON(http.StatusOK, resp)
		return
	}

	textList2, total := pageFindByGorm(params)

	if len(textList2) < 1 {
		// Set "list":[], or "list":null
		textList2 = make([]sttext.Text, 0)
	}

	// Response data
	retList := queryTagsForTextList(textList2)

	queryCommentsForText(&retList)

	pfData := pkg.ApiPageFindResponse{
		PageNo:    params.PageNo,
		PageSize:  params.PageSize,
		Total:     0,
		TotalPage: 0,
		List:      retList,
	}
	resp := pkg.ApiResponse{
		Code:    200,
		Message: "OK",
		Data:    pfData,
	}

	var totalPage int64
	queryTotal := total

	if queryTotal > 0 {
		// Calc
		totalPage = queryTotal / int64(params.PageSize)
		if queryTotal%int64(params.PageSize) != 0 {
			totalPage++
		}
		pfData.Total = queryTotal
		pfData.TotalPage = totalPage
		// set Data again!
		resp.Data = pfData
	}

	c.JSON(http.StatusOK, resp)
}

// Page find text based on gorm
// db.Raw function
func pageFindByGorm(params pageFindParams) ([]sttext.Text, int64) {
	// query data
	sql1 := "SELECT id,content,type,create_time,update_time" +
		" FROM text"
	whreBody, dqlParams := makePageFidWhereBody(params)
	sql1 += whreBody

	// ordery by
	sql1 += " ORDER BY create_time DESC"

	// paging
	offset := (params.PageNo - 1) * params.PageSize
	sql1 += " LIMIT ?,?"
	dqlParams = append(dqlParams, offset)
	dqlParams = append(dqlParams, params.PageSize)

	db := database.GetDB()

	var textList []sttext.Text
	db.Raw(sql1, dqlParams...).Scan(&textList)

	// query total
	var total int64
	sql2 := "SELECT count(1) FROM text"
	sql2 += whreBody
	countParams := dqlParams[:len(dqlParams)-2]
	db.Raw(sql2, countParams...).Scan(&total)

	return textList, total
}

func makePageFidWhereBody(params pageFindParams) (dql string, dqlParams []any) {
	// where clause
	whereBody := ""

	// Tags
	tags := params.Tags
	tagsLen := len(tags)
	if tagsLen > 0 {
		whereBody += " id IN (SELECT t1.text_id FROM text_tag t1"

		if tagsLen == 1 {
			whereBody += " WHERE t1.name = ?" + ")"
			dqlParams = append(dqlParams, tags[0])
			goto moreCond
		}

		if tagsLen > 1 {
			whereBody += " JOIN text_tag t2 ON t1.text_id = t2.text_id" +
				" AND t1.name = ?" +
				" AND t2.name = ?"
			dqlParams = append(dqlParams, tags[0])
			dqlParams = append(dqlParams, tags[1])
		}
		if tagsLen > 2 {
			whereBody += " JOIN text_tag t3 ON t1.text_id = t3.text_id" +
				" AND t3.name = ?"
			dqlParams = append(dqlParams, tags[2])
		}
		if tagsLen > 3 {
			whereBody += " JOIN text_tag t4 ON t1.text_id = t4.text_id" +
				" AND t4.name = ?"
			dqlParams = append(dqlParams, tags[3])
		}
		if tagsLen > 4 {
			whereBody += " JOIN text_tag t5 ON t1.text_id = t5.text_id" +
				" AND t5.name = ?"
			dqlParams = append(dqlParams, tags[4])
		}

		whereBody += ")"
	}

moreCond:

	// KwContent
	if params.KwContent != "" {
		if whereBody == "" {
			whereBody += " content LIKE concat('%', ?, '%')"
		} else {
			whereBody += " AND content LIKE concat('%', ?, '%')"
		}
		dqlParams = append(dqlParams, params.KwContent)
	}

	// Type
	if params.Type != "" {
		if whereBody == "" {
			whereBody += " type = ?"
		} else {
			whereBody += " AND type = ?"
		}
		dqlParams = append(dqlParams, params.Type)
	}

	if whereBody != "" {
		dql += " WHERE " + whereBody
	}

	return
}

func queryTagsForTextList(records []sttext.Text) (retList []PageFindVO) {
	if len(records) < 1 {
		return make([]PageFindVO, 0)
	}

	tids := make([]int32, len(records))
	for i, rd := range records {
		var vo PageFindVO
		vo.Id = rd.Id
		vo.Content = rd.Content
		vo.Type = rd.Type
		vo.CreateTime = rd.CreateTime
		vo.UpdateTime = rd.UpdateTime
		vo.Tags = make([]sttext.TextTag, 0)

		retList = append(retList, vo)

		// for query
		tids[i] = vo.Id
	}

	// Query tags of a records
	var tagsAll []sttext.TextTag
	db := database.GetDB()
	db.Where("text_id IN ?", tids).Find(&tagsAll)

	for _, tag := range tagsAll {
		tid := tag.TextId
		for j, vo := range retList {
			void := vo.Id
			if void == tid {
				vo.Tags = append(vo.Tags, tag)
				retList[j] = vo
				break
			}
		}
	}

	return
}

func queryCommentsForText(textList *[]PageFindVO) {
	const maxComments = 5

	var wg sync.WaitGroup

	for i, v := range *textList {
		wg.Add(1)

		go func(textId int32) {
			defer wg.Done()

			var comments []sttext.TextComment
			db := database.GetDB()
			db.Where("text_id = ?", textId).
				Order("create_time desc").
				Limit(maxComments).
				Find(&comments)

			(*textList)[i].Comments = comments

			var count int64
			db.Model(&sttext.TextComment{}).Where("text_id =?", textId).Count(&count)
			(*textList)[i].TotalOfComments = count
		}(v.Id)
	}

	wg.Wait()
}

// ---comment---

// Add a comment for a text
func AddTextComment(c *gin.Context) {
	var params AddTextCommentParams
	if err := c.ShouldBind(&params); err != nil {
		resp := pkg.ApiResponse{
			Code:    400,
			Message: err.Error(),
		}
		c.JSON(http.StatusOK, resp)
		return
	}

	// Check if the text exists.
	var trec sttext.Text
	_, err := database.GetRecordByPk(&trec, params.TextId)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		resp := pkg.ApiResponse{
			Code:    400,
			Message: "No text with this id",
		}
		c.JSON(http.StatusOK, resp)
		return
	} else if err != nil {
		log.Println("Add comment failed: err=", err.Error(), ", params=", params)
		resp := pkg.ApiResponse{
			Code:    500,
			Message: "Server error",
			Data:    err.Error(),
		}
		c.JSON(http.StatusOK, resp)
		return
	}

	record := &sttext.TextComment{
		TextId:  params.TextId,
		Comment: params.Comment,
	}

	database.AddOneRecord(record)
	id := record.Id
	log.Println("Add comment success. id=", id, ", textId=", params.TextId)

	// Success.
	resp := pkg.ApiResponse{
		Code:    200,
		Message: "OK",
		Data:    id,
	}
	c.JSON(http.StatusOK, resp)
}

// Delete a comment
func DeleteTextCommentById(c *gin.Context) {
	var params1 DelTextCommentParams
	if err := c.ShouldBindUri(&params1); err != nil {
		resp := pkg.ApiResponse{
			Code:    400,
			Message: "Invalid params: " + err.Error(),
		}
		c.JSON(http.StatusOK, resp)
		return
	}

	// find comment with textId and id first
	var cmt sttext.TextComment
	db := database.GetDB()
	result := db.First(&cmt, params1.Id)
	err := result.Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		resp := pkg.ApiResponse{
			Code:    200,
			Message: "No such comment or deleted",
			Data:    true,
		}
		c.JSON(http.StatusOK, resp)
		return
	} else if err != nil {
		log.Println("Query comment failed: err=", err.Error(), ", params=", params1)
		resp := pkg.ApiResponse{
			Code:    500,
			Message: "Query comment failed: " + err.Error(),
		}
		c.JSON(http.StatusOK, resp)
		return
	}

	if params1.TextId != cmt.TextId {
		log.Println("Delete comment failed: Comment doesn't belong to the text. params=", params1)
		resp := pkg.ApiResponse{
			Code:    http.StatusBadRequest,
			Message: "Comment doesn't belong to the text",
		}
		c.JSON(http.StatusOK, resp)
		return
	}

	// exec delete
	result = db.Delete(&sttext.TextComment{}, cmt.Id)
	err = result.Error
	if err != nil {
		log.Println("Delete comment failed: err=", err.Error(), ", params=", params1)
		resp := pkg.ApiResponse{
			Code:    500,
			Message: "Delete comment failed: " + err.Error(),
		}
		c.JSON(http.StatusOK, resp)
		return
	}

	log.Printf("Delete comment success: textId=%d, id=%d\n", params1.TextId, params1.Id)

	// Success.
	resp := pkg.ApiResponse{
		Code:    200,
		Message: "OK",
		Data:    true,
	}
	c.JSON(http.StatusOK, resp)
}

func GetComments(c *gin.Context) {
	panic("unimplemented")
}
