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
			public.Post("/login", srv.loginWithEmailPassword)
			public.Get("/country-state", srv.getCountryAndState)
			//public.Get("/user", srv.fetchUser)
			public.Route("/", func(user chi.Router) {
				user.Use(srv.MiddlewareProvider.Middleware())
				user.Post("/register", srv.register)
				user.Post("/change-password", srv.changePassword)
				user.Post("/logout", srv.logout)
				user.Get("/dashboard", srv.dashboard)
				user.Get("/users", srv.Users)
				user.Route("/profile", func(profile chi.Router) {
					profile.Get("/", srv.fetchUser)
					profile.Put("/", srv.editProfile)
				})
			})
			public.Route("/order", func(order chi.Router) {
				order.Use(srv.MiddlewareProvider.Middleware())
				order.Post("/", srv.Order)
				order.Post("/scan", srv.ScanQR)
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
