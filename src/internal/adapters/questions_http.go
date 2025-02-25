package adapters

import (
	"dc_honest/src/internal/core/ports/input"
	"dc_honest/src/internal/infrastructure/ms"
	"dc_honest/src/pkg"
	"fmt"
	"github.com/gin-gonic/gin"
)

type QuestionsAdapterHttp struct {
	service    input.QuestionPort
	levelsRepo *ms.LevelsMsRepo
}

func NewQuestionsAdapterHttp(router *gin.Engine, service input.QuestionPort, repo *ms.LevelsMsRepo) *QuestionsAdapterHttp {
	d := &QuestionsAdapterHttp{
		service:    service,
		levelsRepo: repo,
	}

	router.GET("/api/v1/question", pkg.WithError(d.GetRandQuestion))

	return d
}

// GetRandQuestion godoc
// @Summary      Получить рандомный вопрос
// @Accept       json
// @Produce      json
// @Param		 clientId query string true "client id"
// @Param		 levelId query string true "level id"
// @Success      200  {object}  QuestionOutput
// @Router       /api/v1/question [get]
func (q *QuestionsAdapterHttp) GetRandQuestion(c *gin.Context) error {
	clientID := c.Query("clientId")
	levelID := c.Query("levelId")

	exists, err := q.levelsRepo.LevelExists(levelID)
	if err != nil {
		return err
	}
	if !exists {
		c.JSON(400, gin.H{
			"error": fmt.Sprintf("level %v does not exist", levelID),
		})
		return nil
	}

	question, isLast, err := q.service.GetRandomQuestion(levelID, clientID)
	if err != nil {
		return err
	}
	out := ToOutputDto(question, isLast)
	c.JSON(200, out)
	return nil
}
