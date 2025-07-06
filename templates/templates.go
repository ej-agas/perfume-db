package templates

import (
	_ "embed"
	"html/template"
	"io/fs"
	"os"
	"path/filepath"
)

//go:embed base.html
var Base string

//go:embed components/404.html
var NotFoundTemplate string

//go:embed components/house_card.html
var HouseCardTemplate string

//go:embed components/house.html
var HouseTemplate string

//go:embed components/perfume.html
var PerfumeTemplate string

//go:embed components/perfumer.html
var PerfumerTemplate string

func NotFound() *template.Template {
	tmpl, _ := template.New("404").Option("missingkey=zero").Parse(NotFoundTemplate)
	return tmpl
}

// HouseCard returns a template for rendering a single house card
func HouseCard() *template.Template {
	tmpl, _ := template.New("house-card").Option("missingkey=zero").Parse(HouseCardTemplate)
	return tmpl
}

// House returns a template for rendering a house detail page
func House() *template.Template {
	tmpl, _ := template.New("house").Option("missingkey=zero").Funcs(template.FuncMap{
		"safeHTML": func(s string) template.HTML { return template.HTML(s) },
	}).Parse(HouseTemplate)

	return tmpl
}

func Perfume() *template.Template {
	tmpl, _ := template.New("perfume").Option("missingkey=zero").Funcs(template.FuncMap{
		"safeHTML": func(s string) template.HTML { return template.HTML(s) },
	}).Parse(PerfumeTemplate)

	return tmpl
}

func Pefumer() *template.Template {
	tmpl, _ := template.New("perfumer").Option("missingkey=zero").Funcs(template.FuncMap{
		"safeHTML": func(s string) template.HTML { return template.HTML(s) },
	}).Parse(PerfumerTemplate)

	return tmpl
}

//go:embed index.html
var Index string

// makeSlice creates a slice of integers from 0 to n-1
func makeSlice(n int) []struct{} {
	return make([]struct{}, n)
}

// NewTemplates parses all the templates and returns a template.Template
func NewTemplates() (*template.Template, error) {
	// Create a new template set with the base template first
	tmpl := template.New("").Option("missingkey=zero")

	// Add custom functions
	tmpl = tmpl.Funcs(template.FuncMap{
		"makeSlice": makeSlice,
		"safeHTML":  func(s string) template.HTML { return template.HTML(s) },
	})

	var err error
	tmpl, err = tmpl.Parse(Base)
	if err != nil {
		return nil, err
	}

	// Parse the index template into the same template set
	_, err = tmpl.Parse(Index)
	if err != nil {
		return nil, err
	}

	// Parse all HTML files in the components directory
	componentsDir := "templates/components"
	if _, err := os.Stat(componentsDir); err == nil {
		err = filepath.Walk(componentsDir, func(path string, info fs.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() && filepath.Ext(path) == ".html" {
				content, err := os.ReadFile(path)
				if err != nil {
					return err
				}
				_, err = tmpl.New(filepath.Base(path)).Parse(string(content))
				if err != nil {
					return err
				}
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
	}

	return tmpl, nil
}
