package handlers

import (
	"github.com/ej-agas/perfume-db/internal"
	"net/http"
)

type HouseHandler struct {
	Service internal.HouseService
}

func (h HouseHandler) CreateHouse(w http.ResponseWriter, r *http.Request) {

	data := struct {
		Status  int
		Message string
	}{
		Status:  200,
		Message: "Create House Handler",
	}

	JSONResponse(w, data, 200)
}
