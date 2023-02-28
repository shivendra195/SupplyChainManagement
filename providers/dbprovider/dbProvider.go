package dbprovider

import (
	"time"

	"example.com/supplyChainManagement/providers"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type psqlProvider struct {
	db *sqlx.DB
}

func NewPSQLProvider(connectionString string) providers.DBProvider {
	var (
		db          *sqlx.DB
		err         error
		maxAttempts = 3
	)

	for i := 0; i < maxAttempts; i++ {
		db, err = sqlx.Connect("postgres", connectionString)
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

	return &psqlProvider{
		db: db,
	}
}

func (pp *psqlProvider) Ping() error {
	return pp.db.Ping()
}

func (pp *psqlProvider) DB() *sqlx.DB {
	return pp.db
}

//
//func (db *	DB) Migration() error {
//
//	m, err := migrate.New("file:///database/migration/", "postgres://postgres:postgres@127.0.0.1:5432/go_graphql?sslmode=disable"))
//	println(m)
//	if err != nil {
//		// **I get error here!!**
//		return fmt.Errorf("error happened when migration")
//	}
//	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
//		return fmt.Errorf("error when migration up: %v", err)
//	}
//
//	log.Println("migration completed!")
//	return err
//}
