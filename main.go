package main

import (
	"log"
	"net/http"
	"time"

	_ "github.com/zlbenjamin/gotextgin/prjswagger"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	_ "github.com/go-sql-driver/mysql"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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

// create or update tables with gorm
func init() {
	db := database.GetDB()
	err := db.AutoMigrate(&pkg.Text{})
	if err != nil {
		log.Panicln("create or update table text_comment failed:", err.Error())
	}
	log.Println("create or update table text success.")

	err = db.AutoMigrate(&pkg.TextComment{})
	if err != nil {
		log.Panicln("create or update table text_comment failed:", err.Error())
	}
	log.Println("create or update table text_comment success.")

	err = db.AutoMigrate(&pkg.TextTag{})
	if err != nil {
		log.Panicln("create or update table text_tag failed:", err.Error())
	}
	log.Println("create or update table text_tag success.")
}

func main() {
	r := gin.New()

	// r.Use(middlewares.LoggerApi(), middlewares.CustomRecovery())
	r.Use(middlewares.LoggerApi(), gin.Recovery())

	r.NoRoute(middlewares.Handle404())

	// url := ginSwagger.URL("http://localhost:40000/swagger/doc.json")
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// custom validators
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		log.Println("Register validators...")
		v.RegisterValidation("checktags", service.CheckTags)
	}

	api.InitTextApis(r)

	// addr := ":40000" // ok
	// addr := "0.0.0.0:40000" // same with the up addr
	// Warning, Can't deploy on container!
	addr := "localhost:40000"
	s := &http.Server{
		Addr:           addr,
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 2 << 20,
	}

	log.Println("Start web server at " + addr)
	s.ListenAndServe()
}
