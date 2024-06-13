package api

import (
	"github.com/gin-gonic/gin"
	"github.com/zlbenjamin/gotextgin/service"
)

const apiGroupText = "/api/text"

// Initialize APIs
func InitTextApis(r *gin.Engine) {
	textGroup := r.Group(apiGroupText)

	// ---text---
	textGroup.POST("add", AddText)
	textGroup.GET(":id", GetTextById)
	textGroup.DELETE(":id", DeleteTextById)
	textGroup.POST("page", PageFindText)

	//---comment---
	textGroup.POST("comment/add", AddTextComment)
	textGroup.DELETE("comment/:textId/:id", DeleteTextCommentById)

	//---tag---

}

// @Summary Add a text
// @Accept json
// @Produce json
// @Param data body service.AddTextParams true "request body"
// @Router /api/text/add [post]
// @Success 200 {object} pkg.ApiResponse{code=int} "code=200, 400, 500 or self-defined codes"
// @Failure 500 {object} pkg.ApiResponse "other errors, such as network errors"
func AddText(c *gin.Context) {
	service.AddText(c)
}

// @Summary Get a text by its primary key
// @Produce json
// @Param id path int true "PK of text"
// @Router /api/text/{id} [get]
// @Success 200 {object} pkg.ApiResponse{code=int} "code=200, 400, 500 or self-defined codes"
// @Failure 500 {object} pkg.ApiResponse "other errors, such as network errors"
func GetTextById(c *gin.Context) {
	service.GetTextById(c)
}

// @Summary Delete a text by its primary key
// @Produce json
// @Param id path int true "PK of text"
// @Router /api/text/{id} [delete]
// @Success 200 {object} pkg.ApiResponse{code=int} "code=200, 400, 500 or self-defined codes"
// @Failure 500 {object} pkg.ApiResponse "other errors, such as network errors"
func DeleteTextById(c *gin.Context) {
	service.DeleteTextById(c)
}

// @Summary Paging text
// @Accept json
// @Produce json
// @Param data body service.pageFindParams true "request body"
// @Router /api/text/page [post]
// @Success 200 {object} pkg.ApiResponse{code=int,data=pkg.ApiPageFindResponse{list=[]service.PageFindVO}}
// "code=200, 400, 500 or self-defined codes"
// @Failure 500 {object} pkg.ApiResponse "other errors, such as network errors"
func PageFindText(c *gin.Context) {
	service.PageFindText(c)
}

// @Summary Add a comment for a text
// @Accept json
// @Produce json
// @Param data body service.AddTextCommentParams true "request body"
// @Router /api/text/comment/add [post]
// @Success 200 {object} pkg.ApiResponse{code=int} "code=200, 400, 500 or self-defined codes"
// @Failure 500 {object} pkg.ApiResponse "other errors, such as network errors"
func AddTextComment(c *gin.Context) {
	service.AddTextComment(c)
}

// @Summary Delete a comment
// @Produce json
// @Param data path service.DelTextCommentParams true "path variables"
// @Router /api/text/comment/{textId}/{id} [delete]
// @Success 200 {object} pkg.ApiResponse{code=int} "code=200, 400, 500 or self-defined codes"
// @Failure 500 {object} pkg.ApiResponse "other errors, such as network errors"
func DeleteTextCommentById(c *gin.Context) {
	service.DeleteTextCommentById(c)
}

func GetComments(c *gin.Context) {
	service.GetComments(c)
}
