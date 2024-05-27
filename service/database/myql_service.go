package database

import (
	"database/sql"
	"errors"
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

// Check if the table exists
func CheckTableExists(tableName string) bool {
	rows, err := gsqlDb.Query("show tables")
	if err != nil {
		log.Panicln("Show tables. err=", err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		var item string
		errq := rows.Scan(&item)
		if errq != nil {
			log.Panicln("rows.Scan. err=", err.Error())
		}

		if tableName == item {
			return true
		}
	}

	return false
}

// Create a table with the DDL
func CreateMySqlTable(ddl string) {
	_, err := gsqlDb.Exec(ddl)
	if err != nil {
		log.Fatalln("Create table failed: ddl=", ddl, err.Error())
	}
	// fmt.Println(result, reflect.TypeOf(result)) // {0xc00040a000 0xc00040c090} sql.driverResult
	// fmt.Println(result.RowsAffected())          // 0 <nil>
	// fmt.Println(result.LastInsertId())          // 0 <nil>
}

// Add a record to table
func AddRecordToTable(dml string, params ...any) (int32, error) {
	result, err := gsqlDb.Exec(dml, params...)
	if err != nil {
		return -1, errors.New("Add record failed")
	}

	id, err := result.LastInsertId()
	if err != nil {
		return -2, errors.New("Add record failed")
	}
	return int32(id), nil
}
