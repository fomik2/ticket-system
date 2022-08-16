package app

import (
	"log"
	"net/http"

	"github.com/fomik2/ticket-system/internal/handlers"
	"github.com/gorilla/mux"
)

func Run(cfg map[string]string, repo handlers.Tickets) {

	r := mux.NewRouter()
	handler := handlers.New(cfg)

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		handler.WelcomeHandler(w, r, repo)
	}).Methods("GET")

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		handler.CreateTicket(w, r, repo)
	}).Methods("POST")

	r.HandleFunc("/tickets/{id:[0-9]+}", func(w http.ResponseWriter, r *http.Request) {
		handler.EditHandler(w, r, repo)
	}).Methods("POST")

	r.HandleFunc("/tickets/{id:[0-9]+}", func(w http.ResponseWriter, r *http.Request) {
		handler.GetTicketForEdit(w, r, repo)
	}).Methods("GET")

	http.Handle("/", r)

	fs := http.FileServer(http.Dir(cfg["css_path"]))
	http.Handle("/css/", http.StripPrefix("/css/", fs))

	err := http.ListenAndServe(cfg["http_port"], nil)
	if err != nil {
		log.Fatal("Problem related to starting HTTP server", err)
	}
}
