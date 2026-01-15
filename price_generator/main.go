package main

import (
	"log"
	"price_generator/internal/bootstrap"

	"github.com/joho/godotenv"
)

func main() {
	// Load .env file (silently fails if not found - that's OK for Docker)
	_ = godotenv.Load()
	
	app := bootstrap.NewApp()
	if err := app.Run(); err != nil {
		log.Fatalf("Application failed: %v", err)
	}
}
