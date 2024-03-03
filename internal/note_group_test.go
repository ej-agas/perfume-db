package internal

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewNoteGroup(t *testing.T) {
	Name := "Fruits, Vegetables And Nuts"
	Slug := "fruits-vegetables-and-nuts"
	Description := "Foo"
	ImageURL := "/images/fruits-vegtables-and-nuts.png"

	noteGroup := NewNoteGroup(Name, Description, ImageURL)

	assert.Equal(t, Name, noteGroup.Name)
	assert.Equal(t, Slug, noteGroup.Slug)
	assert.Equal(t, Description, noteGroup.Description)
	assert.Equal(t, ImageURL, noteGroup.ImageURL)
}
