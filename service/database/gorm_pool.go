package database

import (
	"fmt"
	"log"
	"time"

	"github.com/zlbenjamin/gotextgin/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var db *gorm.DB

// Initialize var db with #config.DatabaseSetting
func init() {
	log.Println("Init database pool: type=", config.DatabaseSetting.Type)
	switch config.DatabaseSetting.Type {
	case "mysql":
		initMySql(config.DatabaseSetting)
	default:
		log.Fatalln("Unsupported database type. type=", config.DatabaseSetting.Type)
	}
}

// Init mysql pool
func initMySql(setting *config.Database) {
	uri := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=true",
		setting.User,
		setting.Password,
		setting.Host,
		setting.Port,
		setting.Name,
	)
	dialector := mysql.New(mysql.Config{
		DSN:                       uri,
		DefaultStringSize:         256,
		DisableDatetimePrecision:  true,
		DontSupportRenameIndex:    true,
		DontSupportRenameColumn:   true,
		SkipInitializeWithVersion: false,
	})

	conn, err := gorm.Open(dialector, &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   config.DatabaseSetting.TablePrefix, // set table prefix
			SingularTable: true,                               // set table singular
		},
	})
	if err != nil {
		log.Fatalln("gorm.Open err=", err.Error())
	}

	sqlDB, err := conn.DB()
	if err != nil {
		log.Fatalln("conn.DB err=", err.Error())
	}

	// set pool
	sqlDB.SetMaxOpenConns(20)
	sqlDB.SetConnMaxIdleTime(10)
	sqlDB.SetConnMaxLifetime(10 * time.Minute)
	// sqlDB.SetConnMaxIdleTime(10 * time.Minute)

	db = conn
}

// Get DB
func GetDB() *gorm.DB {
	sqlDB, err := db.DB()
	if err != nil {
		log.Println("Connect db server failed. err=", err.Error())
	}

	if err = sqlDB.Ping(); err != nil {
		log.Println("Ping db servr failed. err=", err.Error())
		sqlDB.Close()
	}

	return db
}
