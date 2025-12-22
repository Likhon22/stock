package handler

import (
	"log"
	"net/http"
	"stock-processor/internal/service"
	"stock-processor/utils"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type Handler struct{
	service service.ProcessorService
}

func NewHandler(service service.ProcessorService)*Handler  {
	return &Handler{
		service: service,
	}
	
}

func (h *Handler) GetPrice(w http.ResponseWriter, r *http.Request) {
    
	  symbol:=chi.URLParam(r,"symbol")
    price, err := h.service.GetCache(r.Context(), symbol)
    if err != nil {
        http.Error(w, "Not found", 404)
        return
    }
    
		data:=map[string]interface{}{
			"symbol":symbol,
			"price":price,
		}
    if err:=utils.SendResponse(w,"successfully send price for the symbol",http.StatusOK,&data);err!=nil {
			log.Fatal("err happen",err)
		}
}


func (h *Handler) GetAllPrices(w http.ResponseWriter, r *http.Request) {
    prices, err := h.service.GetAll(r.Context())
    if err != nil {
        http.Error(w, "Error", 500)
        return
    }
    
    if err:=utils.SendResponse(w,"successfully send all the prices",http.StatusOK,&prices);err!=nil {
			log.Fatal("err happen",err)
		}
}


func (h *Handler) GetHistory(w http.ResponseWriter, r *http.Request) {
	symbol:=chi.URLParam(r,"symbol")
	limit:=chi.URLParam(r,"limit")
	value, err := strconv.ParseInt(limit, 10, 64)
	 if err != nil {
        http.Error(w, "Error", 500)
        return
    }
    prices, err := h.service.GetHistory(r.Context(),symbol,int(value))
    if err != nil {
        http.Error(w, "Error", 500)
        return
    }
    
    if err:=utils.SendResponse(w,"successfully send price history",http.StatusOK,&prices);err!=nil {
			log.Fatal("err happen",err)
		}
}