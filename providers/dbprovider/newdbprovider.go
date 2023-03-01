package dbprovider

import (
	"database/sql"
	"github.com/shivendra195/supplyChainManagement/providers"
	"github.com/sirupsen/logrus"
	"time"

	_ "github.com/lib/pq"
)

type pSQLCProvider struct {
	db *sql.DB
}

func NewSQLCProvider(connection string) providers.NewDBProvider {

	var (
		db          *sql.DB
		err         error
		maxAttempts = 3
	)

	for i := 0; i < maxAttempts; i++ {
		db, err = sql.Open("postgres", connection)
		if err != nil {
			logrus.Errorf("unable to connect to postgres PSQL %v", err)
			time.Sleep(3 * time.Second)
			continue
		}
		break
	}

	if err != nil {
		logrus.Fatalf("Failed to initialize PSQL: %v", err)
	} else {
		logrus.Info("connected to postgresql database")
	}

	return &pSQLCProvider{
		db: db,
	}
}

func (pp *pSQLCProvider) Ping() error {
	return pp.db.Ping()
}

func (pp *pSQLCProvider) DB() *sql.DB {
	return pp.db
}
