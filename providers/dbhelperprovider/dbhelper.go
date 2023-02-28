package dbhelperprovider

import (
	"example.com/supplyChainManagement/providers"
	"github.com/jmoiron/sqlx"
)

type DBHelper struct {
	DB *sqlx.DB
}

func NewDBHepler(db *sqlx.DB) providers.DBHelperProvider {
	return &DBHelper{
		DB: db,
	}
}
