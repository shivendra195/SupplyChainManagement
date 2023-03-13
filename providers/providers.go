package providers

import (
	"context"
	"database/sql"
	"github.com/go-chi/chi"
	"github.com/jmoiron/sqlx"
	"github.com/shivendra195/supplyChainManagement/models"
	"net/http"
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

type AuthProvider interface {
	generateJWT(userToken string) (string, error)
}

type MiddlewareProvider interface {
	Middleware() func(next http.Handler) http.Handler
	UserFromContext(ctx context.Context) models.UserContextData
	Default() chi.Middlewares
}
