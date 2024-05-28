package api

import (
	"github.com/gin-gonic/gin"
	"github.com/zlbenjamin/gotextgin/service"
)

const apiGroupText = "/api/text"

// Initialize APIs
func InitTextApis(r *gin.Engine) {
	textGroup := r.Group(apiGroupText)
	textGroup.POST("add", AddText)
	textGroup.GET(":id", GetTextById)
	textGroup.DELETE(":id", DeleteTextById)
	textGroup.POST("page", PageFindText)
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
