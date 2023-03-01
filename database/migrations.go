package database

import (
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"

	// source/file import is required for migration files to read
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"

	//"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func Migrations(db *sqlx.DB) {
	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		logrus.Fatalf("Migrations : Failed to postgres.WithInstance : %v", err)
	}
	m, err := migrate.NewWithDatabaseInstance("file://database", "postgres", driver)
	if err != nil {
		logrus.Fatalf("Migrations : Failed to postgres.NewWithDatabaseInstance : %v", err)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		logrus.Fatalf("Migrations : Failed to up and down migrations  : %v", err)
	}
}
