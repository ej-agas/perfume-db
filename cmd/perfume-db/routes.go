package main

import (
	"encoding/json"
	"net/http"
	"os"
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
		Message: "Perfume DB API foo bar baz",
	}

	json.NewEncoder(w).Encode(data)
}

func (app *application) routes() http.Handler {
	// Create a new mux that will be our main router
	mux := http.NewServeMux()

	// Use fixed path for static files in the container
	staticDir := "/home/perfume_db_user/static"

	// Check if static directory exists
	if _, err := os.Stat(staticDir); os.IsNotExist(err) {
		app.logger.Error("static directory does not exist", "path", staticDir)
	} else {
		app.logger.Info("Serving static files from", "directory", staticDir)
		fs := http.FileServer(http.Dir(staticDir))
		mux.Handle("GET /static/", http.StripPrefix("/static/", fs))
	}

	// Create a subrouter for UI API endpoints
	uiApiRouter := http.NewServeMux()
	uiApiRouter.HandleFunc("GET /houses", app.listHousesUI)
	uiApiRouter.HandleFunc("GET /houses/new", app.newHouseFormUI)

	// Mount UI API under /ui
	mux.Handle("/ui/", http.StripPrefix("/ui", uiApiRouter))

	// Create a subrouter for API endpoints
	apiRouter := http.NewServeMux()
	apiRouter.HandleFunc("GET /health", Home)

	// House routes
	apiRouter.HandleFunc("POST /houses", app.createHouseHandler)
	apiRouter.HandleFunc("GET /houses", app.listHouses)
	apiRouter.HandleFunc("GET /houses/{slug}", app.showHouseBySlug)
	apiRouter.HandleFunc("PATCH /houses/{publicId}", app.updateHouseByPublicId)

	// Note Group routes
	apiRouter.HandleFunc("POST /note-groups", app.createNoteGroupHandler)
	apiRouter.HandleFunc("GET /note-groups", app.listNoteGroups)
	apiRouter.HandleFunc("GET /note-groups/{slug}", app.showNoteGroupBySlug)
	apiRouter.HandleFunc("PATCH /note-groups/{publicId}", app.updateNoteGroupByPublicId)

	// Note routes
	apiRouter.HandleFunc("POST /notes", app.createNoteHandler)
	apiRouter.HandleFunc("GET /notes", app.listNotes)
	apiRouter.HandleFunc("GET /notes/{slug}", app.showNoteBySlug)
	apiRouter.HandleFunc("PATCH /notes/{publicId}", app.updateNoteByPublicId)

	// Perfumer routes
	apiRouter.HandleFunc("POST /perfumers", app.createPerfumerHandler)
	apiRouter.HandleFunc("PATCH /perfumers/{publicId}", app.updatePerfumerByPublicIdHandler)
	apiRouter.HandleFunc("GET /perfumers", app.listPerfumersHandler)
	apiRouter.HandleFunc("GET /perfumers/{slug}", app.showPerfumerBySlugHandler)

	// Perfume routes
	apiRouter.HandleFunc("POST /perfumes", app.createPerfumeHandler)

	// Mount API under /api
	mux.Handle("/api/", http.StripPrefix("/api", apiRouter))

	// Handle the root path
	mux.HandleFunc("GET /home", app.homeUI)

	return mux
}
