package app

import (
	"fmt"

	"log"
	"net/http"

	"github.com/fomik2/ticket-system/config"
	"github.com/fomik2/ticket-system/internal/filerw"
	"github.com/fomik2/ticket-system/internal/handlers"
)

func Run(cfg *config.Config) {
	filerw.ReadTicketsFromFiles()
	http.HandleFunc("/", handlers.WelcomeHandler)
	http.HandleFunc("/tickets/", handlers.EditHandler)
	fs := http.FileServer(http.Dir(cfg.CSS.Path))
	http.Handle("/css/", http.StripPrefix("/css/", fs))
	fmt.Println(cfg.HTTP.Port)
	err := http.ListenAndServe(cfg.HTTP.Port, nil)
	if err != nil {
		log.Fatal("Problem related to starting HTTP server", err)
	}
}
