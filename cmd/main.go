package main

import (
	"electra/internal/app"
	"electra/internal/config"
	"github.com/joho/godotenv"
	"log"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, using system env")
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	application := app.New(*cfg)

	if err := application.Init(); err != nil {
		log.Fatalf("failed to init app: %v", err)
	}

	if err := application.Run(); err != nil {
		log.Fatalf("failed to run app: %v", err)
	}
}
