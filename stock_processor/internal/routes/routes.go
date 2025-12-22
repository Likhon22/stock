package routes

import (
	"stock-processor/internal/handler"

	"github.com/go-chi/chi/v5"
)

func SetUpRoutes(processorHandler *handler.Handler)*chi.Mux  {
	r := chi.NewRouter()
	r.Get("/price/{symbol}",processorHandler.GetPrice)
	r.Get("/prices",processorHandler.GetAllPrices)
	r.Get("/price/history/{symbol}/{limit}",processorHandler.GetHistory)
	return r
}

