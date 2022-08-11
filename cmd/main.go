package main

import (
	"log"

	"github.com/fomik2/ticket-system/config"
	"github.com/fomik2/ticket-system/internal/app"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	app.Run(cfg)
}
