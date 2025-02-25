package service

import (
	"database/sql"
	"dc_honest/src/internal/core"
	"dc_honest/src/internal/core/domain"
	"dc_honest/src/internal/core/ports/output"
)

type QuestionsService struct {
	db   *sql.DB
	repo output.QuestionRepoPort
}

func NewQuestionsService(db *sql.DB, repo output.QuestionRepoPort) *QuestionsService {
	return &QuestionsService{
		db:   db,
		repo: repo,
	}
}

func (q *QuestionsService) GetRandomQuestion(levelID string, clientID string) (question domain.Question, b bool, err error) {
	tx, err := q.db.Begin()
	if err != nil {
		return domain.Question{}, false, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// понять все ли взято
	allQuestionsGot, err := q.repo.AreAllQuestionsGot(tx, clientID, levelID)
	if err != nil {
		return domain.Question{}, false, err
	}
	if !allQuestionsGot {
		question, err = q.repo.GetRandomQuestion(tx, clientID, levelID)
		if err != nil {
			return question, false, err
		}
		err = q.repo.AddQuestionToHistory(tx, clientID, question)
		if err != nil {
			return question, false, err
		}
		err = q.repo.AddQuestionToUsedQuestions(tx, clientID, question)
		if err != nil {
			return question, false, err
		}
	} else {
		question = domain.Question{
			ID:             "-1",
			Level:          domain.Level{ID: levelID},
			Text:           core.GetConfig().LastCardText,
			AdditionalText: nil,
		}
		err = q.repo.ClearUsedQuestions(tx, clientID, levelID)
		if err != nil {
			return question, false, err
		}
	}

	err = tx.Commit()
	if err != nil {
		return question, false, err
	}
	return question, allQuestionsGot, nil
}
