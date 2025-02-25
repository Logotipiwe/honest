package main

import (
	"database/sql"
	"dc_honest/docs"
	_ "dc_honest/docs"
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
	"strconv"
)

// @title           Swagger honest API
// @version         1.0
// @description     This is a honest service api.

// @BasedPath  /

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {
	LoadDcConfig()

	config := core.GetConfig()

	db, err := sql.Open("mysql", config.GetMysqlConnectionStr())
	if err != nil {
		log.Fatal(err)
	}

	fw := flyway.NewFlyway(db, "./data/migrations/")
	err = fw.Migrate()
	if err != nil {
		log.Fatal(err)
	}

	decksStorage := ms.NewDecksMsStorage(db)
	shuffleRepo := ms.NewShuffleRepoMs(db)
	questionsRepo := ms.NewQuestionMsRepo(db)
	levelsRepo := ms.NewLevelsMsRepo(db)

	decksService := service.NewDecksService(decksStorage)
	shuffleService := service.NewShuffleService(shuffleRepo)
	questionsService := service.NewQuestionsService(db, questionsRepo)

	_ = application.App{
		DecksPort:   decksService,
		ShufflePort: shuffleService,
	}

	router := gin.Default()
	_ = adapters.NewDecksAdapterHttp(router, decksService)
	_ = adapters.NewShuffleHttpAdapter(router, shuffleService)
	_ = adapters.NewQuestionsAdapterHttp(router, questionsService, levelsRepo)
	adapters.HandlerSwaggerRoute(router)

	setupSwagger(config)

	err = router.Run(":" + strconv.Itoa(config.Port))
	if err != nil {
		log.Fatal(err)
	}
}

func setupSwagger(config *core.Config) {
	if config.SwaggerHost == "" {
		panic("swagger host not set")
	}
	if config.SwaggerBasePath == "" {
		panic("swagger base path not set")
	}
	docs.SwaggerInfo.Host = config.SwaggerHost
	docs.SwaggerInfo.BasePath = config.SwaggerBasePath
}
