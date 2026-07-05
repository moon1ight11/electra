package main

import (
	"electra/internal/app"
	"electra/internal/config"
	"log"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("fatal error in config load: %v", err)
	}

	application := app.New(*cfg)

	if err := application.Init(); err != nil {
		log.Fatalf("fatal error in app init: %v", err)
	}

	if err := application.Run(); err != nil {
		log.Fatalf("fatal error in app run: %v", err)
	}
}
