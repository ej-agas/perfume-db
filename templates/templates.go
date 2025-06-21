package templates

import (
	_ "embed"
	"fmt"
	"html/template"
)

//go:embed base.html
var Base string

//go:embed index.html
var Index string

// NewTemplates parses all the templates and returns a template.Template
func NewTemplates() (*template.Template, error) {
	// Create a new template set with the base template first
	tmpl, err := template.New("").Parse(Base)
	if err != nil {
		return nil, err
	}

	// Parse the index template into the same template set
	_, err = tmpl.Parse(Index)
	if err != nil {
		return nil, err
	}

	fmt.Printf("%#v\n", tmpl)

	return tmpl, nil
}
