package server

import (
	"context"
	"github.com/shivendra195/supplyChainManagement/providers"
	"github.com/shivendra195/supplyChainManagement/providers/dbhelperprovider"
	"github.com/shivendra195/supplyChainManagement/providers/dbprovider"
	"github.com/shivendra195/supplyChainManagement/providers/middlewareprovider"
	"log"
	"net/http"
	"time"
	// "example.com/supplyChainManagement/providers/dbhelperprovider"
	//"example.com/supplyChainManagement/providers/dbprovider"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

const (

	// jdbc:postgresql://localhost:5432/scmdb
	defaultPSQLURL    = "postgresql://admin:root@localhost:5432/scmdb?sslmode=require"
	defaultPortNumber = "80"
)

type Server struct {
	//AdminQueries *admin.Queries
	//PSQLC        providers.PSQLCProvider
	//PSQLC     	providers.PSQLProvider
	//MiddleProvider providers.
	//AuthProvider providers.AuthProvider
	MiddlewareProvider providers.MiddlewareProvider
	DBHelper           providers.DBHelperProvider
	PSQL               providers.PSQLProvider
	httpServer         *http.Server
}

func SrvInit() *Server {

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Some error occured. Err: %s", err)
	}

	// connection to database
	//PSQLC := dbprovider.NewSQLCProvider(os.Getenv("PSQL_DB_URL"))

	//getting queries database
	//AdminQueries := admin.New(PSQLC.DB())
	//dbQuries  := newdbprovider.QueriesProviders(dbsqlc)

	//PSQL connection
	//db := dbprovider.NewPSQLProvider(os.Getenv("PSQL_DB_URL"))
	//db := dbprovider.NewPSQLProvider(PSQL_DB_URL)
	db := dbprovider.NewPSQLProvider(defaultPSQLURL)

	// database helper functions
	dbHelper := dbhelperprovider.NewDBHepler(db.DB())

	middleware := middlewareprovider.NewMiddleware(dbHelper)

	return &Server{
		PSQL:               db,
		DBHelper:           dbHelper,
		MiddlewareProvider: middleware,
	}
}

func (srv *Server) Start() {
	addr := ":" + defaultPortNumber

	httpSrv := &http.Server{
		Addr:    addr,
		Handler: srv.InjectRoutes(),
	}
	srv.httpServer = httpSrv

	logrus.Info("Server running at PORT ", addr)
	if err := httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logrus.Fatalf("Start %v", err)
		return
	}
}

func (srv *Server) Stop() {
	logrus.Info("closing Postgres...")
	_ = srv.PSQL.DB().Close()
	//_ = srv.PSQLC.DB().Close()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	logrus.Info("closing server...")
	_ = srv.httpServer.Shutdown(ctx)
	logrus.Info("Done")
}
