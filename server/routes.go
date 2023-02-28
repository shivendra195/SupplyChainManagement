package server

import (
	"github.com/gorilla/mux"
	"net/http"
)

func (srv *Server) InjectRoutes() *mux.Router {
	r := mux.NewRouter()
	//r.HandleFunc("/health-check", srv.HealthCheck).Methods("GET")
	//r.HandleFunc("/api", srv.home).Methods("POST")
	r.HandleFunc("/register", srv.CreateUser).Methods("POST")
	r.HandleFunc("/users", srv.FetchAllUser).Methods("GET")

	http.Handle("/", r)
	return r
}
