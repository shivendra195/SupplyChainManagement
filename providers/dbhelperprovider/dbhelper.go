package dbhelperprovider

import (
	"github.com/jmoiron/sqlx"
	"github.com/shivendra195/supplyChainManagement/providers"
)

type DBHelper struct {
	DB *sqlx.DB
}

func NewDBHepler(db *sqlx.DB) providers.DBHelperProvider {
	return &DBHelper{
		DB: db,
	}
}
