package config

import (
	"encoding/json"
	"log"
	"os"
)

type Database struct {
	Type        string `json:"type"`
	Name        string `json:"name"`
	Host        string `json:"host"`
	Port        string `json:"port"`
	User        string `json:"user"`
	Password    string `json:"password"`
	TablePrefix string `json:"tablePrefix"`
}

// public
var DatabaseSetting = &Database{}

type Config struct {
	Database *Database `json:"database"`
}

var GlobalConfigSetting = &Config{}

func init() {
	configPath := os.Getenv("gomysqlConfig")
	if configPath == "" {
		log.Println("Warning: There is no env named 'gomysqlConfig'.")
		return
	}

	log.Println("Start to config DB by file:", configPath)

	cfile, err := os.Open(configPath)
	if err != nil {
		log.Fatalln("Invalid gomysqlConfig. err=", err.Error())
	}
	defer cfile.Close()

	// json decoder
	decoder := json.NewDecoder(cfile)
	err = decoder.Decode(GlobalConfigSetting)
	if err != nil {
		log.Fatalln("gomysqlConfig can't be decoded. err=", err.Error())
	}

	log.Println("Get config DB success.")
	DatabaseSetting = GlobalConfigSetting.Database
}
