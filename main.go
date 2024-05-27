package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func init() {
	// Connect to Database
}

func init() {
	// Create table
}

func main() {
	r := gin.Default()

	s := &http.Server{
		Addr:           ":40000",
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 2 << 20,
	}

	s.ListenAndServe()
}
