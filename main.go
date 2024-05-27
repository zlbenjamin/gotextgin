package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/zlbenjamin/gotextgin/api"
	pkg "github.com/zlbenjamin/gotextgin/pkg/text"
	"github.com/zlbenjamin/gotextgin/service/database"
)

func init() {
	// Connect to Database
	database.InitMySqlPool()
}

func init() {
	// Create table
	if database.CheckTableExists(pkg.Table_Text) {
		return
	}

	log.Println("Start create table: ", pkg.Table_Text)

	ddl := `
	CREATE TABLE text (
		id INT NOT NULL AUTO_INCREMENT COMMENT 'PK',
		content MEDIUMTEXT NOT NULL COMMENT 'text content',
		type VARCHAR(100) NULL COMMENT 'type, such markdown, golang, c++, java, python, html, javascript etc.',
		create_time DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'create time',
		update_time DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'the lastest update time',
		PRIMARY KEY (id))
	  COMMENT = 'text';
	`
	database.CreateMySqlTable(ddl)

	log.Println("Create table: ", pkg.Table_Text, "success.")
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
