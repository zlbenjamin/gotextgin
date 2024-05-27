package service

import (
	"fmt"
	"net/http"
	"strings"

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
		fmt.Println("todo err=", err)
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

	fmt.Println("todo params=", params)

	dml := "INSERT INTO " + sttext.Table_Text + " (content, type) VALUES (?, ?)"
	id, err := database.AddRecordToTable(dml, params.Content, params.Type)
	if err != nil {
		fmt.Println("todo err=", err)
		resp := pkg.ApiResponse{
			Code:    500,
			Message: "Add failed",
		}
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	// Success.
	resp := pkg.ApiResponse{
		Code:    200,
		Message: "OK",
		Data:    id,
	}
	c.JSON(http.StatusBadRequest, resp)
}
