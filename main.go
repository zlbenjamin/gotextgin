package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	_ "github.com/go-sql-driver/mysql"
	"github.com/zlbenjamin/gotextgin/api"
	"github.com/zlbenjamin/gotextgin/pkg/middlewares"
	pkg "github.com/zlbenjamin/gotextgin/pkg/text"
	"github.com/zlbenjamin/gotextgin/service"
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

// create or update tables with gorm
func init() {
	db := database.GetDB()
	err := db.AutoMigrate(&pkg.TextComment{})
	if err != nil {
		log.Panicln("create table text_comment failed:", err.Error())
	}
	log.Println("create table text_comment success.")

	err = db.AutoMigrate(&pkg.TextTag{})
	if err != nil {
		log.Panicln("create table text_tag failed:", err.Error())
	}
	log.Println("create table text_tag success.")
}

func main() {
	r := gin.New()

	// r.Use(middlewares.LoggerApi(), middlewares.CustomRecovery())
	r.Use(middlewares.LoggerApi(), gin.Recovery())

	r.NoRoute(middlewares.Handle404())

	// custom validators
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		log.Println("Register validators...")
		v.RegisterValidation("checktags", service.CheckTags)
	}

	api.InitTextApis(r)

	addr := ":40000"
	s := &http.Server{
		Addr:           addr,
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 2 << 20,
	}

	log.Println("Start web server at 0.0.0.0" + addr)
	s.ListenAndServe()
}
