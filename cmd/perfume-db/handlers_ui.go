package main

import (
	"github.com/ej-agas/perfume-db/templates"
	"net/http"
)

func (app *application) render(w http.ResponseWriter, r *http.Request, name string, data interface{}) {
	// Get the parsed templates
	ts, err := templates.NewTemplates()
	if err != nil {
		app.logger.Error("template parsing error", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Execute the template with the provided data
	// The name parameter should match the template name (e.g., "base" for base.html)
	if err := ts.ExecuteTemplate(w, name, data); err != nil {
		app.logger.Error("template execution error", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func (app *application) homeUI(w http.ResponseWriter, r *http.Request) {
	// The template name should match the template definition in base.html
	app.render(w, r, "base", nil)
}

func (app *application) listHousesUI(w http.ResponseWriter, r *http.Request) {
	houses, err := app.services.House.List(0, 10)
	if err != nil {
		app.logger.Error("failed to list houses", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// For now, just return a simple list - we'll enhance this later
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	for _, house := range houses {
		_, _ = w.Write([]byte(`<div class="mb-2">` + house.Name + `</div>`))
	}
}

func (app *application) newHouseFormUI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write([]byte(`
        <h3 class="text-lg font-semibold mb-4">Add New Perfume House</h3>
        <form hx-post="/houses" hx-target="#houses-list" @submit="modalOpen = false">
            <div class="mb-4">
                <label class="block text-gray-700 text-sm font-bold mb-2" for="name">
                    House Name
                </label>
                <input class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline" 
                       id="name" name="name" type="text" required>
            </div>
            <div class="flex justify-end space-x-2">
                <button type="button" @click="modalOpen = false" 
                        class="bg-gray-300 hover:bg-gray-400 text-gray-800 font-bold py-2 px-4 rounded">
                    Cancel
                </button>
                <button type="submit" 
                        class="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded">
                    Save
                </button>
            </div>
        </form>
    `))
}
