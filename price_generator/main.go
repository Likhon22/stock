package main

import (
	"log"
	"price_generator/internal/bootstrap"
)

func main() {
	app := bootstrap.NewApp()
	if err := app.Run(); err != nil {
		log.Fatalf("Application failed: %v", err)
	}
}
