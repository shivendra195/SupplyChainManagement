package server

import (
	//"github.com/gorilla/mux"
	"github.com/go-chi/chi"
)

func (srv *Server) InjectRoutes() *chi.Mux {
	r := chi.NewRouter()

	r.Get(`/health`, srv.HealthCheck)

	r.Route("/api", func(api chi.Router) {
		api.Route("/public", func(public chi.Router) {
			public.Post("/register", srv.register)
			public.Post("/login", srv.loginWithEmailPassword)
			public.Get("/user", srv.fetchUser)
			public.Route("/", func(user chi.Router) {
				user.Use(srv.MiddlewareProvider.Middleware())
				user.Get("/profile", srv.fetchUser)
				user.Post("/logout", srv.logout)
			})
		})
	})
	return r
}

//func (srv *Server) InjectRoutes() *mux.Router {
//	r := mux.NewRouter()
//	//r.HandleFunc("/health-check", srv.HealthCheck).Methods("GET")
//	//r.HandleFunc("/api", srv.home).Methods("POST")
//	r.HandleFunc("/register", srv.CreateUser).Methods("POST")
//	r.HandleFunc("/users", srv.FetchAllUser).Methods("GET")
//
//	http.Handle("/", r)
//	return r
//}
