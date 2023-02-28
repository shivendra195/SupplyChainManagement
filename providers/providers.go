package providers

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
)

// PSQLProvider provides the database client from the Database after connecting to database.
type PSQLProvider interface {
	// DB returns the database client.
	DB() *sqlx.DB
}

type DBProvider interface {
	// Ping verifies the connection with the database.
	Ping() error
	PSQLProvider
}

// PSQLProvider provides the database client from the Database after connecting to database.
type PSQLCProvider interface {
	// DB returns the database client.
	DB() *sql.DB
}

type NewDBProvider interface {
	// Ping verifies the connection with the database.
	Ping() error
	PSQLCProvider
}
