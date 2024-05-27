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
}

func AddText(c *gin.Context) {
	service.AddText(c)
}

func DeleteTextById(c *gin.Context) {
	// TODO
}

func GetTextById(c *gin.Context) {
	// TODO
}

func PageFindText(c *gin.Context) {
	// TODO
}
