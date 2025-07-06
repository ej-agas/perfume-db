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

	// Handle the root path
	mux.HandleFunc("GET /", app.homeUI)
	mux.HandleFunc("GET /houses/{slug}", app.showHouseUI)
	mux.HandleFunc("GET /perfumes/{slug}", app.showPerfumeUI)
	mux.HandleFunc("GET /perfumers/{slug}", app.showPerfumerUI)

	// Serve static files from the local static directory
	staticDir := "./static"
	app.logger.Info("Serving static files from", "directory", staticDir)

	// Create a file server that serves files from the static directory
	fs := http.FileServer(http.Dir(staticDir))

	// Handle requests to /static/ by stripping the /static/ prefix
	mux.Handle("GET /static/", http.StripPrefix("/static/", fs))

	// Log the contents of the static directory for debugging
	if entries, err := os.ReadDir(staticDir); err == nil {
		app.logger.Info("Static directory contents:", "path", staticDir)
		for _, entry := range entries {
			app.logger.Info("  " + entry.Name())
		}
	} else {
		app.logger.Error("Failed to read static directory", "error", err, "path", staticDir)
	}

	// Create a subrouter for UI API endpoints
	uiApiRouter := http.NewServeMux()
	uiApiRouter.HandleFunc("GET /houses", app.listHousesUI)

	uiApiRouter.HandleFunc("GET /houses/new", app.newHouseFormUI)

	mux.HandleFunc("GET /ui/", func(w http.ResponseWriter, r *http.Request) {
		originalPath := r.URL.Path
		r.URL.Path = originalPath[len("/ui"):] // Strip "/ui"

		if r.URL.Path == "" { // Handle /api/ base path
			w.WriteHeader(http.StatusNotFound)
			return
		}

		uiApiRouter.ServeHTTP(w, r)
		// Restore original path for any downstream middleware or logging
		r.URL.Path = originalPath
	})

	return mux
}
