package ms

import "database/sql"

type LevelsMsRepo struct {
	db *sql.DB
}

func NewLevelsMsRepo(db *sql.DB) *LevelsMsRepo {
	return &LevelsMsRepo{db: db}
}

func (l *LevelsMsRepo) LevelExists(levelID string) (bool, error) {
	var count int
	err := l.db.QueryRow("SELECT count(*) FROM levels where id = ?", levelID).Scan(&count)
	return count > 0, err
}
