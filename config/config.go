package config

import (
	"fmt"
	"log"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
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

var (
	cfg  = pflag.StringP("config", "c", "", "Configuration file")
	help = pflag.BoolP("help", "h", false, "Show this help message")
)

func init() {
	fmt.Println(1)
	pflag.Parse()
	if *help {
		pflag.Usage()
		return
	}

	// read config
	if *cfg != "" {
		viper.SetConfigFile(*cfg)
		viper.SetConfigType("yaml")
	} else {
		viper.AddConfigPath(".")
		viper.AddConfigPath("$HOME/.gotextgin")
		viper.SetConfigName("gotextgin")
	}

	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("fatal error config file: %s", err))
	}

	log.Printf("[viper]Used configuation file is: %s\n", viper.ConfigFileUsed())
}

func init() {
	// initialize first
	GlobalConfigSetting.Database = &Database{}

	GlobalConfigSetting.Database.Type = viper.GetString("database.type")
	GlobalConfigSetting.Database.Name = viper.GetString("database.name")
	GlobalConfigSetting.Database.Host = viper.GetString("database.host")
	GlobalConfigSetting.Database.Port = viper.GetString("database.port")
	GlobalConfigSetting.Database.User = viper.GetString("database.user")
	GlobalConfigSetting.Database.Password = viper.GetString("database.password")
	GlobalConfigSetting.Database.TablePrefix = viper.GetString("database.table_prefix")

	fmt.Println("todo delete:")
	fmt.Println(GlobalConfigSetting)
	fmt.Println(GlobalConfigSetting.Database)

	log.Println("Get config DB success.")
	DatabaseSetting = GlobalConfigSetting.Database
}
