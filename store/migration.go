package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
	_ "github.com/lib/pq"
)

var DB *sql.DB

func init() {
	var err error
	DB, err = sql.Open("postgres", "postgresql://postgres:123456@172.28.40.251/?sslmode=disable")
	driver, err := postgres.WithInstance(DB, &postgres.Config{})
	checkerr(err, "migration get driver")

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres", driver)
	checkerr(err, "migration NewWithDatabaseInstance")

	m.Steps(2)
	fmt.Println("init done.")
}

func checkerr(err error, phase string) {
	if err != nil {
		fmt.Println(phase+" ", err)
		os.Exit(1)
	}
}
