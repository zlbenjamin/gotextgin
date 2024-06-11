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

func AddText(c *gin.Context) {
	service.AddText(c)
}

func GetTextById(c *gin.Context) {
	service.GetTextById(c)
}

func DeleteTextById(c *gin.Context) {
	service.DeleteTextById(c)
}

func PageFindText(c *gin.Context) {
	service.PageFindText(c)
}

// Add a comment for a text
func AddTextComment(c *gin.Context) {
	service.AddTextComment(c)
}

// Delete a comment
func DeleteTextCommentById(c *gin.Context) {
	service.DeleteTextCommentById(c)
}

func GetComments(c *gin.Context) {
	service.GetComments(c)
}
