package main

import (
	"database/sql"
	"log"

	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
	_ "github.com/lib/pq"
)

func init() {
	var err error
	DB, err = sql.Open("postgres", "postgresql://postgres:123456@localhost/?sslmode=disable")
	driver, err := postgres.WithInstance(DB, &postgres.Config{})
	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres", driver)
	if err != nil {
		log.Fatal("migration fatal", err)
	}
	m.Steps(2)
}
