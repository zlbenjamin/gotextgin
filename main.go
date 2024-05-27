package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/zlbenjamin/gotextgin/api"
	"github.com/zlbenjamin/gotextgin/service/database"
)

func init() {
	// Connect to Database
	database.InitMySqlPool()
}

func init() {
	// Create table
}

func main() {
	r := gin.Default()

	api.InitTextApis(r)

	s := &http.Server{
		Addr:           ":40000",
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 2 << 20,
	}

	s.ListenAndServe()
}
