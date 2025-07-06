package main

import (
	"fmt"
	"github.com/ej-agas/perfume-db/internal"
	"github.com/ej-agas/perfume-db/templates"
	"html/template"
	"net/http"
	"time"
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
	time.Sleep(250 * time.Millisecond)
	if err != nil {
		app.logger.Error("failed to list houses", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	for _, house := range houses {
		tmpl := templates.HouseCard()
		if err := tmpl.Execute(w, house); err != nil {
			app.logger.Error("error executing house card template", "error", err)
		}
	}
}

// HouseWithPerfumes represents the data needed for the house detail page
type HouseWithPerfumes struct {
	*internal.House
	Perfumes *[]*internal.Perfume
}

func (app *application) showHouseUI(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")

	house, err := app.services.House.FindBySlug(slug)
	if err != nil {
		app.logger.Error("failed to find house", "error", err, "slug", slug)
		app.renderErrorTemplate(w, "404", http.StatusNotFound)
		return
	}

	// Fetch perfumes for the house
	perfumes, err := app.services.House.FindPerfumesByHouse(*house)
	if err != nil {
		app.logger.Error("failed to fetch perfumes", "error", err, "house_id", house.ID)
		// Continue with empty perfumes slice on error
		emptyPerfumes := make([]*internal.Perfume, 0)
		perfumes = &emptyPerfumes
	}

	data := HouseWithPerfumes{
		House:    house,
		Perfumes: perfumes,
	}

	// Check if this is an HTMX request
	if r.Header.Get("HX-Request") == "true" {
		tmpl := templates.House()
		if tmpl == nil {
			app.logger.Error("house template is nil")
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		if err := tmpl.Execute(w, data); err != nil {
			app.logger.Error("error executing house template", "error", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	// Full page load
	app.renderTemplate(w, r, "house", data)
}

// renderTemplate is a helper to render templates with the base layout
func (app *application) renderTemplate(w http.ResponseWriter, r *http.Request, templateName string, data interface{}) {
	tmpl, err := templates.NewTemplates()
	if err != nil {
		app.logger.Error("error creating templates", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Create a temporary template that defines "content" as our template
	tmpl = template.Must(tmpl.Parse(fmt.Sprintf(`{{define "content"}}{{template "%s" .}}{{end}}`, templateName)))
	if err := tmpl.ExecuteTemplate(w, "base", data); err != nil {
		app.logger.Error("error executing base template", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// renderErrorTemplate is a helper to render error templates
func (app *application) renderErrorTemplate(w http.ResponseWriter, templateName string, statusCode int) {
	w.WriteHeader(statusCode)
	app.renderTemplate(w, nil, templateName, nil)
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
