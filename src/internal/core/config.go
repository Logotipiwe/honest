package core

import (
	"dc_honest/src/pkg"
	"fmt"
)

type Config struct {
	DbLogin    string
	DbPassword string
	DbName     string
	DbHost     string
}

func NewConfig() *Config {
	return &Config{
		DbLogin:    pkg.OsGetNonEmpty("DB_LOGIN"),
		DbPassword: pkg.OsGetNonEmpty("DB_PASS"),
		DbName:     pkg.OsGetNonEmpty("DB_NAME"),
		DbHost:     pkg.OsGetNonEmpty("DB_HOST"),
	}
}

func (c Config) GetMysqlConnectionStr() string {
	return fmt.Sprintf("%v:%v@tcp(%v)/%v", c.DbLogin, c.DbPassword, c.DbHost, c.DbName)
}
