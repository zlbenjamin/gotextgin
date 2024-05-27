package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

const driverName = "mysql"

var gsqlDb *sql.DB

// MySql connect info
type DbInfo struct {
	User     string
	Password string
	Host     string
	Port     int
	DbName   string
}

// Initialize db pool
func InitMySqlPool() *sql.DB {
	if gsqlDb != nil {
		return gsqlDb
	}

	dbInfo := getMySqlInfoFromEnvironment()
	// 1.?parseTime=true
	// 2.?parseTime=True&loc=Local&charset=utf8
	dataSourceName := fmt.Sprintf("%s:%s@(%s:%d)/%s?parseTime=true",
		dbInfo.User, dbInfo.Password, dbInfo.Host, dbInfo.Port, dbInfo.DbName)
	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		log.Panicln("sql.Open err=", err.Error())
	}

	// set pool params
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(60 * time.Minute)

	// !
	gsqlDb = db

	err = db.Ping()
	if err != nil {
		log.Panicln("sql.Open err=", err.Error())
	}

	log.Println("Init mysql pool success. MaxOpenConnections=", db.Stats().MaxOpenConnections)
	return db
}

// Get MySql info from Environment
// see /docs/addEnvironments.md
// TODO more check
func getMySqlInfoFromEnvironment() DbInfo {
	var ret DbInfo

	// os.Getenv
	ret.User = os.Getenv("godbuser")
	ret.Password = os.Getenv("godbpassword")
	ret.Host = os.Getenv("godbhost")

	iport, err := strconv.Atoi(os.Getenv("godbport"))
	if err != nil {
		// strconv.Atoi: parsing "3306  ": invalid syntax
		log.Panicln("Convert failed: env godbport to int.", os.Getenv("godbport"), err.Error())
	}
	ret.Port = iport

	ret.DbName = os.Getenv("godbname")

	return ret

}
