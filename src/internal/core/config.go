package core

import (
	"dc_honest/src/pkg"
	"fmt"
	"os"
)

type Config struct {
	DbLogin         string
	DbPassword      string
	DbName          string
	DbHost          string
	Port            int
	SwaggerHost     string
	SwaggerBasePath string
	LastCardText    string
}

var (
	config *Config
)

func GetConfig() *Config {
	if config == nil {
		config = NewConfig()
	}
	return config
}

func NewConfig() *Config {

	return &Config{
		DbLogin:         pkg.OsGetNonEmpty("DB_LOGIN"),
		DbPassword:      pkg.OsGetNonEmpty("DB_PASS"),
		DbName:          pkg.OsGetNonEmpty("DB_NAME"),
		DbHost:          pkg.OsGetNonEmpty("DB_HOST"),
		Port:            getPort(),
		SwaggerHost:     os.Getenv("SWAGGER_HOST"),
		SwaggerBasePath: os.Getenv("SWAGGER_BASE_PATH"),
		LastCardText:    pkg.OsGetNonEmpty("LAST_CARD_TEXT"),
	}
}

func getPort() int {
	var port = 80
	osPort := pkg.OsGetEnvInt("PORT")
	if osPort != nil {
		port = *osPort
	}
	return port
}

func (c Config) GetMysqlConnectionStr() string {
	return fmt.Sprintf("%v:%v@tcp(%v)/%v", c.DbLogin, c.DbPassword, c.DbHost, c.DbName)
}
