package infrastructure

import (
	"dc_honest/src/internal/core"
	"dc_honest/src/internal/core/domain"
	. "dc_honest/src/internal/core/service"
	"dc_honest/src/internal/infrastructure/ms"
	. "dc_honest/src/tests"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestQuestionMsRepo_GetRandomQuestion(t *testing.T) {
	db := SetupDb(t)
	fw := SetupFlyway(db, t)
	fw.CleanMigrate()
	config := core.GetConfig()
	t.Run("Get rand question works", func(t *testing.T) {
		q := ms.NewQuestionMsRepo(db)
		service := NewQuestionsService(db, q)
		var givenQuestions []domain.Question
		neededLevelQuestions := []domain.Question{
			{ID: "1", Level: domain.Level{ID: "1"}, Text: "text1"},
			{ID: "2", Level: domain.Level{ID: "1"}, Text: "text2"},
			{ID: "3", Level: domain.Level{ID: "1"}, Text: "text3"},
		}
		givenQuestions = append(givenQuestions, neededLevelQuestions...)
		givenQuestions = append(givenQuestions, domain.Question{ID: "4", Level: domain.Level{ID: "2"}, Text: "text4"})
		err := q.SaveQuestions(givenQuestions)
		if err != nil {
			t.Fatal(err)
			return
		}

		var expected []domain.Question
		expected = append(expected, neededLevelQuestions...)
		expected = append(expected,
			domain.Question{ID: "-1", Level: domain.Level{ID: "1"}, Text: config.LastCardText},
		)

		for i := 0; i < 5; i++ {
			var got []domain.Question
			for j := 0; j < len(expected); j++ {
				item, _, err := service.GetRandomQuestion("1", "1")
				if err != nil {
					t.Errorf("GetRandomQuestion() error = %v", err)
					return
				}
				got = append(got, item)
			}
			assert.True(t, SameElements(got, expected), "Got wrong questions. Expected:\n%v\nGot:\n%v", givenQuestions, got)
		}
	})
}

func SameElements[T any](elements []T, required []T) bool {
	for _, req := range required {
		found := false
		for _, el := range elements {
			if reflect.DeepEqual(req, el) {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}
