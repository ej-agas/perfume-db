package main

import (
	"encoding/json"
	"net/http"
	"time"
)

func Home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(200)

	data := struct {
		Status  int    `json:"status"`
		Time    string `json:"server_time"`
		Message string `json:"message"`
	}{
		Status:  200,
		Time:    time.Now().Format("2006-01-02T15:04:05Z07:00"),
		Message: "Perfume DB API",
	}

	json.NewEncoder(w).Encode(data)
}

func (app *application) routes() http.Handler {
	router := http.NewServeMux()

	router.HandleFunc("GET /", Home)
	router.HandleFunc("GET /houses", app.listHouses)
	router.HandleFunc("POST /houses", app.createHouseHandler)
	router.HandleFunc("GET /houses/{slug}", app.showHouseBySlug)

	return router
}
