package handlers

import (
	"github.com/ej-agas/perfume-db/internal"
	"net/http"
)

type HouseHandler struct {
	Service internal.HouseService
}

func (h HouseHandler) CreateHouse(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(202)
}
