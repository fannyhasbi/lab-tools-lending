package config

import (
	"database/sql"
	"fmt"
	"log"

	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func InitPostgresDB() *sql.DB {
	DBHost := os.Getenv("DB_HOST")
	DBDriver := os.Getenv("DB_DRIVER")
	DBUser := os.Getenv("DB_USER")
	DBPass := os.Getenv("DB_PASSWORD")
	DBPort := os.Getenv("DB_PORT")
	DBName := os.Getenv("DB_NAME")

	connStr := fmt.Sprintf("%s://%s:%s@%s:%s/%s?sslmode=disable", DBDriver, DBUser, DBPass, DBHost, DBPort, DBName)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	m, err := migrate.New(
		"file://database/migration",
		connStr,
	)
	log.Println("checkpoint 1")
	if err != nil {
		log.Fatal(err)
	}

	log.Println("checkpoint 2")

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal(err)
	}

	log.Println("checkpoint 3")

	return db
}
