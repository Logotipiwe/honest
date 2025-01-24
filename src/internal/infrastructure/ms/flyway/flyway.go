package flyway

import (
	"crypto/sha256"
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"sort"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

type Flyway struct {
	DB         *sql.DB
	Migrations string
}

// NewFlyway creates a new Flyway instance
func NewFlyway(db *sql.DB, migrationsPath string) *Flyway {
	return &Flyway{DB: db, Migrations: migrationsPath}
}

// Migrate applies only new SQL files in the migrations directory
func (f *Flyway) Migrate() error {
	_, err := f.DB.Exec("CREATE TABLE IF NOT EXISTS flyway_schema_history (version VARCHAR(255) PRIMARY KEY, hash VARCHAR(64))")
	if err != nil {
		return err
	}

	rows, err := f.DB.Query("SELECT version, hash FROM flyway_schema_history")
	if err != nil {
		return err
	}
	defer rows.Close()

	appliedMigrations := map[string]string{}
	for rows.Next() {
		var version, hash string
		if err := rows.Scan(&version, &hash); err != nil {
			return err
		}
		appliedMigrations[version] = hash
	}

	files, err := ioutil.ReadDir(f.Migrations)
	if err != nil {
		return err
	}

	sqlFiles := []string{}
	versionRegex := regexp.MustCompile(`V(\d+)__.*\.sql`)
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".sql") {
			sqlFiles = append(sqlFiles, file.Name())
		}
	}

	sort.Strings(sqlFiles) // Ensure migrations run in order

	for _, file := range sqlFiles {
		matches := versionRegex.FindStringSubmatch(file)
		if len(matches) < 2 {
			log.Printf("Skipping invalid migration filename: %s", file)
			continue
		}
		version := matches[1]

		path := fmt.Sprintf("%s/%s", f.Migrations, file)
		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		hash := fmt.Sprintf("%x", sha256.Sum256(content))

		if appliedHash, exists := appliedMigrations[version]; exists {
			if appliedHash == hash {
				continue // Skip already applied and unchanged migrations
			} else {
				return fmt.Errorf("migration file %s has changed after being applied", file)
			}
		}

		log.Printf("Applying migration: %s", file)
		tx, err := f.DB.Begin()
		if err != nil {
			return err
		}

		statements := strings.Split(string(content), ";")
		for _, stmt := range statements {
			stmt = strings.TrimSpace(stmt)
			if stmt == "" {
				continue
			}

			_, err = tx.Exec(stmt)
			if err != nil {
				tx.Rollback()
				return fmt.Errorf("error applying statement in migration %s: %v", file, err)
			}
		}

		_, err = tx.Exec("INSERT INTO flyway_schema_history (version, hash) VALUES (?, ?)", version, hash)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("error recording migration %s: %v", file, err)
		}

		err = tx.Commit()
		if err != nil {
			return err
		}
	}

	log.Println("Migrations applied successfully.")
	return nil
}

// Clean drops all tables from the database
func (f *Flyway) Clean() error {
	rows, err := f.DB.Query("SHOW TABLES")
	if err != nil {
		return err
	}
	defer rows.Close()

	tables := []string{}
	for rows.Next() {
		var table string
		if err := rows.Scan(&table); err != nil {
			return err
		}
		tables = append(tables, table)
	}

	if len(tables) == 0 {
		log.Println("No tables to clean.")
		return nil
	}

	for _, table := range tables {
		_, err := f.DB.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s", table))
		if err != nil {
			return fmt.Errorf("error dropping table %s: %v", table, err)
		}
	}

	_, err = f.DB.Exec("DROP TABLE IF EXISTS flyway_schema_history")
	if err != nil {
		return err
	}

	log.Println("Database cleaned successfully.")
	return nil
}

// Close closes the database connection
func (f *Flyway) Close() error {
	return f.DB.Close()
}
