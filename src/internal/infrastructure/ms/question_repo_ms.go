package ms

import (
	"database/sql"
	"dc_honest/src/internal/core/domain"
)

var (
	GetRandQuestionSql = `SELECT q.id, q.level_id, q.text, q.additional_text 
		FROM questions q 
		where q.level_id = ? 
			AND id not in (
				select question_id from used_questions where client_id = ?
		   ) 
		ORDER BY RAND() LIMIT 1`
)

type QuestionMsRepo struct {
	db *sql.DB
}

func NewQuestionMsRepo(
	db *sql.DB,
) *QuestionMsRepo {
	return &QuestionMsRepo{
		db: db,
	}
}

func (q *QuestionMsRepo) ClearUsedQuestions(tx *sql.Tx, clientID, levelID string) error {
	_, err := tx.Exec(`DELETE FROM used_questions 
       		WHERE question_id IN (SELECT id from questions WHERE level_id = ?)
			AND client_id = ?`, levelID, clientID)
	return err
}

func (q *QuestionMsRepo) AreAllQuestionsGot(tx *sql.Tx, clientID string, levelID string) (bool, error) {
	var allQuestionsGot bool
	err := tx.QueryRow(`WITH q as (
			SELECT * FROM questions WHERE level_id = ?
		),
		u as (SELECT * FROM used_questions WHERE client_id = ? AND question_id IN (SELECT id FROM q))
		SELECT (select count(*) from q) = (select count(*) from u)
    `, levelID, clientID).Scan(&allQuestionsGot)
	if err != nil {
		return false, err
	}
	return allQuestionsGot, nil
}

func (q *QuestionMsRepo) AddQuestionToHistory(tx *sql.Tx, clientID string, question domain.Question) error {
	_, err := tx.Exec("INSERT INTO questions_history (level_id, question_id,client_id) VALUES (?, ?, ?)",
		question.Level.ID, question.ID, clientID)
	return err
}

func (q *QuestionMsRepo) AddQuestionToUsedQuestions(tx *sql.Tx, clientID string, question domain.Question) (err error) {
	_, err = tx.Exec("INSERT INTO used_questions (client_id, question_id) VALUES (?, ?)",
		clientID, question.ID)
	if err != nil {
		return err
	}
	return nil
}

func (q *QuestionMsRepo) GetRandomQuestion(tx *sql.Tx, clientID string, levelID string) (domain.Question, error) {
	var question domain.Question
	row := tx.QueryRow(GetRandQuestionSql, levelID, clientID)
	if row.Err() != nil {
		return question, row.Err()
	}
	err := row.Scan(&question.ID, &question.Level.ID, &question.Text, &question.AdditionalText)
	if err != nil {
		return question, err
	}
	return question, nil
}

func (q *QuestionMsRepo) SaveQuestion(question domain.Question) error {
	return q.SaveQuestions([]domain.Question{question})
}

func (q *QuestionMsRepo) SaveQuestions(questions []domain.Question) error {
	tx, err := q.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	for _, question := range questions {
		_, err := q.db.Exec("INSERT INTO questions (id, level_id, text, additional_text) VALUES (?,?,?,?)",
			question.ID, question.Level.ID, question.Text, question.AdditionalText)
		if err != nil {
			return err
		}
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}
