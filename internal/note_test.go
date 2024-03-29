package internal

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewNote(t *testing.T) {
	Name := "Citrus"
	Slug := "citrus"
	Description := "Citrus Note"
	ImageURL := "/images/citrus-note.png"

	note := NewNote(Name, Description, ImageURL, "123")

	assert.Equal(t, Name, note.Name)
	assert.Equal(t, Slug, note.Slug)
	assert.Equal(t, Description, note.Description)
	assert.Equal(t, ImageURL, note.ImageURL)
}
