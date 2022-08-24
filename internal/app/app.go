package app

import (
	"net/http"

	"github.com/fomik2/ticket-system/internal/handlers"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func Run(index, layout, editor, auth, user_create, http_port, css_path, database, secret string, repo handlers.RepoInterface) error {
	r := mux.NewRouter()
	handler, err := handlers.New(index, layout, editor, auth, user_create, secret, repo)
	if err != nil {
		return err
	}

	// var pingCounter = prometheus.NewCounter(
	// 	prometheus.CounterOpts{
	// 		Name: "ping_request_count",
	// 		Help: "No of request handled by Ping handler",
	// 	},
	// )

	//Prometheus metrics
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":2112", nil)

	r.HandleFunc("/login", handler.Login).Methods("GET")
	r.HandleFunc("/login", handler.LoginHandler).Methods("POST")
	r.HandleFunc("/logout", handler.LogoutHandler).Methods("GET")

	r.HandleFunc("/api/create", handler.JWTAuthMiddleWare(handler.APICreateTicket)).Methods("POST")
	r.HandleFunc("/api/tickets/{id:[0-9]+}", handler.JWTAuthMiddleWare(handler.APIGetTicket)).Methods("GET")
	r.HandleFunc("/api/tickets", handler.JWTAuthMiddleWare(handler.APIGetListTickets)).Methods("GET")
	r.HandleFunc("/api/tickets/byuser", handler.JWTAuthMiddleWare(handler.APIGetListTicketsByUser)).Methods("GET")
	r.HandleFunc("/api/tickets/{id:[0-9]+}", handler.JWTAuthMiddleWare(handler.APIUpdateTicket)).Methods("PUT")
	r.HandleFunc("/api/signin", handler.APISignin).Methods("POST")

	r.HandleFunc("/", handler.Authentication(handler.WelcomeHandler)).Methods("GET")
	r.HandleFunc("/", handler.Authentication(handler.CreateTicket)).Methods("POST")
	r.HandleFunc("/tickets/{id:[0-9]+}", handler.Authentication(handler.EditHandler)).Methods("POST")
	r.HandleFunc("/tickets/{id:[0-9]+}/delete/", handler.Authentication(handler.DeleteHandler)).Methods("POST")
	r.HandleFunc("/tickets/{id:[0-9]+}", handler.Authentication(handler.GetTicketForEdit)).Methods("GET")
	r.HandleFunc("/user_create/", handler.CreateUserGet).Methods("GET")
	r.HandleFunc("/user_create/", handler.CreateUser).Methods("POST")
	http.Handle("/", r)
	fs := http.FileServer(http.Dir(css_path))
	http.Handle("/css/", http.StripPrefix("/css/", fs))
	err = http.ListenAndServe(http_port, nil)
	if err != nil {
		return err
	}
	return nil
}
