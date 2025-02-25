package tests

import (
	"database/sql"
	"dc_honest/src/internal/core"
	"dc_honest/src/internal/infrastructure/ms/flyway"
	"github.com/joho/godotenv"
	"os"
	"path/filepath"
	"testing"
)

func SetupDb(t *testing.T) *sql.DB {
	LoadTestEnv(t)
	config := core.NewConfig()
	db, err := sql.Open("mysql", config.GetMysqlConnectionStr())
	if err != nil {
		t.Fatalf("Failed to connect to real database: %v", err)
	}
	return db
}

func SetupFlyway(db *sql.DB, t *testing.T) *flyway.Flyway {
	LoadTestEnv(t)
	fw := flyway.NewFlyway(db, filepath.Join(GetRootDir(t), "data/migrations"))
	return fw
}

// GetRootDir определяет корневую директорию проекта
func GetRootDir(t *testing.T) string {
	// Получаем путь к файлу, в котором вызвана эта функция
	exePath, err := os.Getwd()
	if err != nil {
		t.Fatalf("Ошибка при получении пути к исполняемому файлу: %v", err)
	}

	// Проходим вверх по директориям, пока не найдем .git или go.mod
	for {
		if _, err := os.Stat(filepath.Join(exePath, "go.mod")); err == nil {
			return exePath
		}
		if _, err := os.Stat(filepath.Join(exePath, ".git")); err == nil {
			return exePath
		}
		parent := filepath.Dir(exePath)
		if parent == exePath {
			break
		}
		exePath = parent
	}

	t.Fatal("Не удалось определить корневую директорию проекта")
	return ""
}

func LoadTestEnv(t *testing.T) {
	rootDir := GetRootDir(t)
	envPath := filepath.Join(rootDir, "src/tests/.env")

	if err := godotenv.Load(envPath); err != nil {
		t.Fatalf("Ошибка загрузки .env: %v", err)
	}
}
