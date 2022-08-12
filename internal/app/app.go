package app

import (
	"log"
	"net/http"

	"github.com/fomik2/ticket-system/internal/filerw"
	"github.com/fomik2/ticket-system/internal/handlers"
	"github.com/gorilla/mux"
)

func Run(cfg map[string]string) {
	filerw.ReadTicketsFromFiles(cfg["tickets"], cfg["counter"])
	r := mux.NewRouter()
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		handlers.WelcomeHandler(w, r, cfg)
	})
	r.HandleFunc("/tickets/{id:[0-9]+}", func(w http.ResponseWriter, r *http.Request) {
		handlers.EditHandler(w, r, cfg)
	})
	http.Handle("/", r)
	fs := http.FileServer(http.Dir(cfg["css_path"]))
	http.Handle("/css/", http.StripPrefix("/css/", fs))

	err := http.ListenAndServe(cfg["http_port"], nil)
	if err != nil {
		log.Fatal("Problem related to starting HTTP server", err)
	}
}
