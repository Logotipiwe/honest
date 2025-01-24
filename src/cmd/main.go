package main

import (
	"database/sql"
	"dc_honest/src/internal/adapters"
	"dc_honest/src/internal/core"
	"dc_honest/src/internal/core/application"
	"dc_honest/src/internal/core/service"
	"dc_honest/src/internal/infrastructure/ms"
	"dc_honest/src/internal/infrastructure/ms/flyway"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	. "github.com/logotipiwe/dc_go_config_lib"
	"log"
)

func main() {
	LoadDcConfig()

	config := core.NewConfig()

	db, err := sql.Open("mysql", config.GetMysqlConnectionStr())
	if err != nil {
		log.Fatal(err)
	}

	fw := flyway.NewFlyway(db, "./data/migrations/")
	err = fw.Migrate()
	if err != nil {
		log.Fatal(err)
	}

	decksStorage := ms.NewDecksStorageMs(db)

	decksService := service.NewDecksService(decksStorage)

	_ = application.App{
		DecksPort: decksService,
	}

	router := gin.Default()

	_ = adapters.NewDecksAdapterHttp(router, decksService)

	err = router.Run(":82")
	if err != nil {
		log.Fatal(err)
	}
}
