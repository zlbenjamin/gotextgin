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
)

// Params of adding text
type AddTextParams struct {
	Content string `json:"content" binding:"required,max=10000"`
	Type    string `json:"type" binding:"max=10"`
}

func (params AddTextParams) ConvertToText() sttext.Text {
	var ret sttext.Text
	ret.Content = params.Content
	ret.Type = params.Type
	return ret
}

// Add a text
func AddText(c *gin.Context) {
	var params AddTextParams
	if err := c.ShouldBind(&params); err != nil {
		resp := pkg.ApiResponse{
			Code:    400,
			Message: err.Error(),
		}
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	// trim
	params.Content = strings.Trim(params.Content, " ")
	params.Type = strings.Trim(params.Type, " ")

	dml := "INSERT INTO " + sttext.Table_Text + " (content, type) VALUES (?, ?)"
	id, err := database.AddRecordToTable(dml, params.Content, params.Type)
	if err != nil {
		log.Println("Add record to table failed. dml=", dml, err.Error())
		resp := pkg.ApiResponse{
			Code:    500,
			Message: "Add text failed",
		}
		c.JSON(http.StatusInternalServerError, resp)
		return
	}

	// Success.
	resp := pkg.ApiResponse{
		Code:    200,
		Message: "OK",
		Data:    id,
	}
	c.JSON(http.StatusOK, resp)
}

// Get a text by id
func GetTextById(c *gin.Context) {
	id := c.Param("id")
	idi, err := strconv.Atoi(id)
	if err != nil {
		log.Println("todo 1 err=", err)
		resp := pkg.ApiResponse{
			Code:    400,
			Message: "Bad request",
		}
		c.JSON(http.StatusBadRequest, resp)
		return
	}
	if idi < 1 {
		log.Println("todo 2 err=", err)
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
		log.Println("todo 1 err=", err)
		resp := pkg.ApiResponse{
			Code:    400,
			Message: "Bad request",
		}
		c.JSON(http.StatusBadRequest, resp)
		return
	}
	if idi < 1 {
		log.Println("todo 2 err=", err)
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

// Params of page find
type pageFindParams struct {
	PageNo    int32  `json:"pageNo" binding:"required"`
	PageSize  int32  `json:"pageSize" binding:"required"`
	KwContent string `json:"kwContent"`
	Type      string `json:"type"`
}

// Page find texts.
// ORDER BY create_time DESC
func PageFindText(c *gin.Context) {
	var params pageFindParams
	if err := c.ShouldBind(&params); err != nil {
		log.Println("Bind err=", err)
		resp := pkg.ApiResponse{
			Code:    400,
			Message: "Bad request",
		}
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	// more check
	if params.PageNo < 1 {
		resp := pkg.ApiResponse{
			Code:    400,
			Message: "pageNo < 1",
		}
		c.JSON(http.StatusBadRequest, resp)
		return
	}
	maxSize := int32(500)
	if params.PageSize < 1 || params.PageSize > maxSize {
		resp := pkg.ApiResponse{
			Code:    400,
			Message: "pageSize: [1, 500]",
		}
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	// todo

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
