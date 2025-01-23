package adapters

import (
	"database/sql"
	"dc_honest/src/internal/adapters/output"
	"dc_honest/src/internal/infrastructure"
	"github.com/joho/godotenv"
	"testing"
)

func TestDecksAdapter(t *testing.T) {
	err := godotenv.Load("../.env")
	if err != nil {
		t.Fatal("Error loading .env file")
	}
	config := infrastructure.NewConfig()
	db, err := sql.Open("mysql", config.GetMysqlConnectionStr())
	if err != nil {
		t.Fatal(err)
	}
	output.NewDecksStorageMs(db)

	t.Run("GetDecksForClient gets", func(t *testing.T) {

	})
}
