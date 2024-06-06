package service

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"

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
		c.JSON(http.StatusBadRequest, resp)
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
			log.Println("Add text failed")
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
	id := c.Param("id")
	idi, err := strconv.Atoi(id)
	if err != nil {
		resp := pkg.ApiResponse{
			Code:    400,
			Message: "Bad request",
		}
		c.JSON(http.StatusBadRequest, resp)
		return
	}
	if idi < 1 {
		resp := pkg.ApiResponse{
			Code:    400,
			Message: "Bad request < 1",
		}
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	var textData sttext.Text

	dql := "SELECT * FROM " + sttext.Table_Text + " WHERE id = ?"
	row := database.GetOneRecrod(dql, idi)
	if err := row.Scan(
		&textData.Id,
		&textData.Content,
		&textData.Type,
		&textData.CreateTime,
		&textData.UpdateTime,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			resp := pkg.ApiResponse{
				Code:    400,
				Message: "No data for id",
			}
			c.JSON(http.StatusBadRequest, resp)
			return
		}

		log.Println("Query failed: id=", id, "err=", err.Error())
		resp := pkg.ApiResponse{
			Code:    500,
			Message: "Query failed",
		}
		c.JSON(http.StatusBadRequest, resp)
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
	id := c.Param("id")
	idi, err := strconv.Atoi(id)
	if err != nil {
		resp := pkg.ApiResponse{
			Code:    400,
			Message: "Bad request",
		}
		c.JSON(http.StatusBadRequest, resp)
		return
	}
	if idi < 1 {
		resp := pkg.ApiResponse{
			Code:    400,
			Message: "Bad request < 1",
		}
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	// Exec
	dml := `DELETE FROM ` + sttext.Table_Text + ` WHERE id = ?`
	var num int64
	if num, err = database.DeleteRecordByPk(dml, idi); err != nil {
		log.Println("Delete text failed: id=", id, "err=", err)
		resp := pkg.ApiResponse{
			Code:    500,
			Message: "Delete failed",
		}
		c.JSON(http.StatusOK, resp)
	}

	resp := pkg.ApiResponse{
		Code:    200,
		Message: "OK",
		Data:    true,
	}

	if num > 0 {
		log.Println("Delete text success. id=", id)
	} else {
		resp.Message = "Do nothing"
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
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	// 1/2 Query data
	dql, dqlParams := makeDql(params)
	rows, err := database.QueryDataList(dql, dqlParams...)
	if err != nil {
		log.Println("Page find texts failed: err=", err.Error(), "params=", params)

		resp := pkg.ApiResponse{
			Code:    500,
			Message: "Page find failed",
		}
		c.JSON(http.StatusInternalServerError, resp)
		return
	}

	var textList []sttext.Text
	for rows.Next() {
		var record sttext.Text

		err := rows.Scan(&record.Id, &record.Content, &record.Type, &record.CreateTime, &record.UpdateTime)
		if err != nil {
			log.Println(err)
			resp := pkg.ApiResponse{
				Code:    500,
				Message: "Page find failed",
			}
			c.JSON(http.StatusInternalServerError, resp)
			return
		}

		textList = append(textList, record)
	}
	if len(textList) < 1 {
		// Set "list":[], or "list":null
		textList = make([]sttext.Text, 0)
	}

	// 2/2 Query total
	pfData := pkg.ApiPageFindResponse{
		PageNo:    params.PageNo,
		PageSize:  params.PageSize,
		Total:     0,
		TotalPage: 0,
		List:      textList,
	}
	resp := pkg.ApiResponse{
		Code:    200,
		Message: "OK",
		Data:    pfData,
	}

	var totalPage int64
	queryTotal := pageFindTotal(params)
	if queryTotal < 0 {
		// failed
		resp.Message = "Get total failed"
		goto labelend
	}

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

labelend:
	c.JSON(http.StatusOK, resp)
}

// Make DQL for paging find texts
func makeDql(params pageFindParams) (dql string, dqlParams []any) {
	dql = "SELECT * FROM " + sttext.Table_Text

	// where clause
	whereBody := ""

	// KwContent
	if params.KwContent != "" {
		whereBody += " content like concat('%', ?, '%')"
		dqlParams = append(dqlParams, params.KwContent)
	}

	// Type
	if params.Type != "" {
		if whereBody == "" {
			whereBody += " type = ?"
		} else {
			whereBody += " and type = ?"
		}
		dqlParams = append(dqlParams, params.Type)
	}
	if whereBody != "" {
		dql += " WHERE " + whereBody
	}

	// ordery by
	dql += " ORDER BY create_time DESC"

	// paging
	offset := (params.PageNo - 1) * params.PageSize
	dql += " LIMIT ?,?"
	dqlParams = append(dqlParams, offset)
	dqlParams = append(dqlParams, params.PageSize)

	return
}

// Query the text satisfying the conditions
// Return -1 means query failed.
func pageFindTotal(params pageFindParams) (total int64) {
	dql := "SELECT count(1) FROM " + sttext.Table_Text
	var dqlParams []any

	// where clause
	whereBody := ""

	// KwContent
	if params.KwContent != "" {
		whereBody += " content like concat('%', ?, '%')"
		dqlParams = append(dqlParams, params.KwContent)
	}

	// Type
	if params.Type != "" {
		if whereBody == "" {
			whereBody += " type = ?"
		} else {
			whereBody += " and type = ?"
		}
		dqlParams = append(dqlParams, params.Type)
	}
	if whereBody != "" {
		dql += " WHERE " + whereBody
	}

	// Exec
	row := database.GetOneRecrod(dql, dqlParams...)
	if err := row.Scan(&total); err != nil {
		log.Println("Query total of texts failed. err=", err.Error())
		return -1
	}

	return
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
		c.JSON(http.StatusBadRequest, resp)
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
		c.JSON(http.StatusBadRequest, resp)
		return
	} else if err != nil {
		log.Println("Add comment failed: err=", err.Error(), ", params=", params)
		resp := pkg.ApiResponse{
			Code:    500,
			Message: "Server error",
			Data:    err.Error(),
		}
		c.JSON(http.StatusInternalServerError, resp)
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
			Message: "Invalid id",
			Data:    err.Error(),
		}
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	var tc sttext.TextComment
	ra, err := database.DeleteOneRecordByPk(&tc, params1.Id)
	if err != nil {
		log.Println("Delete comment failed. id=", params1.Id, "err=", err.Error())
		resp := pkg.ApiResponse{
			Code:    500,
			Message: "Delete error",
			Data:    err.Error(),
		}
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	if ra > 0 {
		log.Println("Delete comment success. id=", params1.Id, "RowsAffected=", ra)
	}

	// Success.
	resp := pkg.ApiResponse{
		Code:    200,
		Message: "OK",
		Data:    true,
	}
	c.JSON(http.StatusOK, resp)
}
