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

	router.HandleFunc("POST /houses", app.createHouseHandler)
	router.HandleFunc("GET /houses", app.listHouses)
	router.HandleFunc("GET /houses/{slug}", app.showHouseBySlug)
	router.HandleFunc("PATCH /houses/{publicId}", app.updateHouseByPublicId)

	router.HandleFunc("POST /note-groups", app.createNoteGroupHandler)
	router.HandleFunc("GET /note-groups", app.listNoteGroups)
	router.HandleFunc("GET /note-groups/{slug}", app.showNoteGroupBySlug)
	router.HandleFunc("PATCH /note-groups/{publicId}", app.updateNoteGroupByPublicId)

	router.HandleFunc("POST /notes", app.createNoteHandler)
	router.HandleFunc("GET /notes", app.listNotes)
	router.HandleFunc("GET /notes/{slug}", app.showNoteBySlug)
	router.HandleFunc("PATCH /notes/{publicId}", app.updateNoteByPublicId)

	router.HandleFunc("POST /perfumers", app.createPerfumerHandler)
	router.HandleFunc("PATCH /perfumers/{publicId}", app.updatePerfumerByPublicIdHandler)
	router.HandleFunc("GET /perfumers", app.listPerfumersHandler)
	router.HandleFunc("GET /perfumers/{slug}", app.showPerfumerBySlugHandler)

	router.HandleFunc("POST /perfumes", app.createPerfumeHandler)
	router.HandleFunc("PATCH /perfumes/{publicId}", app.updatePerfumeHandler)
	router.HandleFunc("GET /perfumes/{slug}", app.showPerfumeBySlug)

	return router
}
