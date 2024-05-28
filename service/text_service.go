package service

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
	"github.com/zlbenjamin/gotextgin/pkg"
	sttext "github.com/zlbenjamin/gotextgin/pkg/text"
	"github.com/zlbenjamin/gotextgin/service/database"
)

type AddTextParams struct {
	Content string `json:"content" binding:"required"`
	Type    string `json:"type"`
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
			Message: "Bad request",
		}
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	// trim
	params.Content = strings.Trim(params.Content, " ")
	params.Type = strings.Trim(params.Type, " ")

	if utf8.RuneCountInString(params.Content) > 10_000 {
		resp := pkg.ApiResponse{
			Code:    400,
			Message: "content: length exceed",
		}
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	if utf8.RuneCountInString(params.Type) > 10 {
		resp := pkg.ApiResponse{
			Code:    400,
			Message: "type: length exceed",
		}
		c.JSON(http.StatusBadRequest, resp)
		return
	}

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

	// TODO

	resp := pkg.ApiResponse{
		Code:    200,
		Message: "OK",
		Data:    true,
	}
	c.JSON(http.StatusOK, resp)
}
