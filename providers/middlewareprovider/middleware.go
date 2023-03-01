package middlewareprovider

import (
	"github.com/jmoiron/sqlx"
	"net/http"
)

func authMiddleware(db *sqlx.DB) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		})
	}
}

func CustomTokenAuthWithClaims(devClaims map[string]interface{}) {

}
