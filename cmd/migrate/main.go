// cmd/migrate/main.go
package main

import (
	"fmt"
	"log"
	"movie-project/config"
	"os"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)

	// Получаем текущую рабочую директорию
	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get current directory: %v", err)
	}

	// Формируем путь к директории с миграциями
	migrationsDir := filepath.Join(currentDir, "migrations")

	// Проверяем существование директории и файла миграции
	migrationFile := filepath.Join(migrationsDir, "001_create_movies_table.sql")
	if _, err := os.Stat(migrationFile); os.IsNotExist(err) {
		log.Fatalf("Migration file does not exist: %s", migrationFile)
	}

	// Создаем URL для источника миграций
	migrationsURL := fmt.Sprintf("file://%s", filepath.ToSlash(migrationsDir))

	log.Printf("Migrations URL: %s", migrationsURL)
	log.Printf("Database URL: %s", dbURL)

	m, err := migrate.New("C:\\Users\\Admin\\GolandProjects\\movie-project\\migrations", dbURL)
	if err != nil {
		log.Fatalf("Failed to create migrate instance: %v", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Failed to apply migrations: %v", err)
	}

	log.Println("Migrations applied successfully")
}
