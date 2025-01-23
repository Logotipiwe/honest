package flyway_test

import (
	"crypto/sha256"
	"database/sql"
	"dc_honest/src/internal/infrastructure"
	. "dc_honest/src/internal/infrastructure/flyway"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func setupRealDB(t *testing.T) *sql.DB {
	godotenv.Load("../.env")
	config := infrastructure.NewConfig()
	db, err := sql.Open("mysql", config.GetMysqlConnectionStr())
	if err != nil {
		t.Fatalf("Failed to connect to real database: %v", err)
	}
	return db
}

func setupMigrationsDir(t *testing.T) string {
	dir, err := ioutil.TempDir("", "migrations")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	return dir
}

func writeMigrationFile(t *testing.T, dir, filename, content string) {
	filePath := filepath.Join(dir, filename)
	if err := ioutil.WriteFile(filePath, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write migration file: %v", err)
	}
}

func calculateHash(content string) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(content)))
}

func TestMigrate(t *testing.T) {
	db := setupRealDB(t)
	defer db.Close()

	dir := setupMigrationsDir(t)
	defer os.RemoveAll(dir)

	writeMigrationFile(t, dir, "V1__init.sql", "CREATE TABLE test (id INT PRIMARY KEY);")
	writeMigrationFile(t, dir, "V2__add_column.sql", "ALTER TABLE test ADD COLUMN name VARCHAR(255);")

	fw := NewFlyway(db, dir)
	if err := fw.Migrate(); err != nil {
		t.Fatalf("Migration failed: %v", err)
	}
	// Check if the table exists
	var tableName string
	err := db.QueryRow("SHOW TABLES LIKE 'test'").Scan(&tableName)
	if err != nil {
		t.Fatalf("Expected table 'test' to exist, but it doesn't: %v", err)
	}
}

func TestMigrate_DuplicateMigration(t *testing.T) {
	db := setupRealDB(t)
	defer db.Close()

	dir := setupMigrationsDir(t)
	defer os.RemoveAll(dir)

	fw := NewFlyway(db, dir)

	writeMigrationFile(t, dir, "V1__init.sql", "CREATE TABLE test (id INT PRIMARY KEY);")

	_ = fw.Migrate()
	if err := fw.Migrate(); err != nil {
		t.Fatal("Expected skipped migration, but got err", err)
	}
}

func TestMigrate_AlteredMigration(t *testing.T) {
	db := setupRealDB(t)
	defer db.Close()

	dir := setupMigrationsDir(t)
	defer os.RemoveAll(dir)

	writeMigrationFile(t, dir, "V1__init.sql", "CREATE TABLE test (id INT PRIMARY KEY);")
	fw := NewFlyway(db, dir)
	_ = fw.Migrate()

	writeMigrationFile(t, dir, "V1__init.sql", "CREATE TABLE test (id INT PRIMARY KEY, name VARCHAR(255));")
	if err := fw.Migrate(); err == nil {
		t.Fatal("Expected migration error due to hash mismatch, but got none")
	}
}
func TestMigrate_Incremental(t *testing.T) {
	db := setupRealDB(t)
	defer db.Close()

	dir := setupMigrationsDir(t)
	defer os.RemoveAll(dir)

	writeMigrationFile(t, dir, "V1__init.sql", "CREATE TABLE test (id INT PRIMARY KEY);")
	fw := NewFlyway(db, dir)
	_ = fw.Migrate()

	writeMigrationFile(t, dir, "V2__add_column.sql", "ALTER TABLE test ADD COLUMN name VARCHAR(255);")
	if err := fw.Migrate(); err != nil {
		t.Fatalf("Incremental migration failed: %v", err)
	}

	// Check if the column exists
	var field, fieldType, isNull, key, defaultValue, extra *string
	err := db.QueryRow("SHOW COLUMNS FROM test WHERE Field = 'name'").Scan(&field, &fieldType, &isNull, &key, &defaultValue, &extra)
	if err != nil {
		t.Fatalf("Expected column 'name' to exist in 'test' table, but it doesn't: %v", err)
	}
}

func TestClean(t *testing.T) {
	db := setupRealDB(t)
	defer db.Close()

	fw := NewFlyway(db, "")
	if err := fw.Clean(); err != nil {
		t.Fatalf("Clean failed: %v", err)
	}
}
