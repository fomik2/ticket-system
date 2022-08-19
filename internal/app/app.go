package app

import (
	"net/http"

	"github.com/fomik2/ticket-system/internal/handlers"
	"github.com/gorilla/mux"
)

func Run(index, editor, tickets, counter, http_port, css_path string, repo handlers.Tickets) error {

	r := mux.NewRouter()
	handler, err := handlers.New(index, editor, tickets, counter, repo)
	if err != nil {
		return err
	}
	r.HandleFunc("/", handler.WelcomeHandler).Methods("GET")
	r.HandleFunc("/", handler.CreateTicket).Methods("POST")
	r.HandleFunc("/tickets/{id:[0-9]+}", handler.EditHandler).Methods("POST")
	r.HandleFunc("/tickets/{id:[0-9]+}", handler.GetTicketForEdit).Methods("GET")
	http.Handle("/", r)

	fs := http.FileServer(http.Dir(css_path))
	http.Handle("/css/", http.StripPrefix("/css/", fs))

	err = http.ListenAndServe(http_port, nil)
	if err != nil {
		return err
	}
	return nil
}
