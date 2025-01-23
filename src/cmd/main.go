package main

import (
	"database/sql"
	"dc_honest/src/internal/adapters/input"
	"dc_honest/src/internal/adapters/output"
	"dc_honest/src/internal/core/application"
	"dc_honest/src/internal/core/service"
	"dc_honest/src/internal/infrastructure"
	"dc_honest/src/internal/infrastructure/flyway"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	. "github.com/logotipiwe/dc_go_config_lib"
	"log"
)

func main() {
	LoadDcConfig()

	config := infrastructure.NewConfig()

	db, err := sql.Open("mysql", config.GetMysqlConnectionStr())
	if err != nil {
		log.Fatal(err)
	}

	fw := flyway.NewFlyway(db, "./data/migrations/")
	err = fw.Migrate()
	if err != nil {
		log.Fatal(err)
	}

	decksStorage := output.NewDecksStorageMs(db)

	decksService := service.NewDecksService(decksStorage)

	_ = application.App{
		DecksPort: decksService,
	}

	router := gin.Default()

	_ = input.NewDecksAdapterHttp(router, decksService)

	err = router.Run(":82")
	if err != nil {
		log.Fatal(err)
	}
}
