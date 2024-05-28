package middlewares

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Custom logger
func LoggerApi() gin.HandlerFunc {
	return func(c *gin.Context) {
		t0 := time.Now()

		c.Next()

		// After request
		latency := time.Since(t0)

		status := c.Writer.Status()
		uri := c.Request.URL
		mtd := c.Request.Method
		raddr := c.Request.RemoteAddr
		var uagent string
		if status != http.StatusOK {
			uagent = ", user_agent=" + c.Request.UserAgent()
		}
		log.Printf("Request: uri=%s %s, status=%d, latency=%v ms, remote=%s%s\n",
			mtd, uri, status, latency.Milliseconds(), raddr, uagent)
	}
}
