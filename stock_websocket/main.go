package main

import (
	"log"
	"stock_websocket/internal/bootstrap"
)

func main() {
	app := bootstrap.NewApp()
	if err := app.Run(); err != nil {
		log.Fatalf("Application failed: %v", err)
	}
}
