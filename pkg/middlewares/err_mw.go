package middlewares

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/zlbenjamin/gotextgin/pkg"
)

// Handle no route with status=404
func Handle404() gin.HandlerFunc {
	return func(c *gin.Context) {
		resp := pkg.ApiResponse{
			Code:    http.StatusInternalServerError,
			Message: "404 Not found",
		}

		c.JSON(c.Writer.Status(), resp)
	}
}

// Handle server error
func CustomRecovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		// [GIN-debug] [WARNING] Headers were already written. Wanted to override status code 200 with 500
		defer func() {
			if err := recover(); err != nil {
				resp := pkg.ApiResponse{
					Code:    http.StatusInternalServerError,
					Message: fmt.Sprintf("err=%v", err),
				}
				c.AbortWithStatusJSON(http.StatusInternalServerError, resp)
			}
		}()

		c.Next()
	}
}
