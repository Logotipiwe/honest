package core

import (
	"dc_honest/src/pkg"
	"fmt"
)

type Config struct {
	DbLogin     string
	DbPassword  string
	DbName      string
	DbHost      string
	Port        int
	SwaggerHost string
}

func NewConfig() *Config {

	return &Config{
		DbLogin:     pkg.OsGetNonEmpty("DB_LOGIN"),
		DbPassword:  pkg.OsGetNonEmpty("DB_PASS"),
		DbName:      pkg.OsGetNonEmpty("DB_NAME"),
		DbHost:      pkg.OsGetNonEmpty("DB_HOST"),
		Port:        getPort(),
		SwaggerHost: pkg.OsGetNonEmpty("SWAGGER_HOST"),
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
