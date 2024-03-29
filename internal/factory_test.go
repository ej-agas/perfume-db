package internal

import (
	"github.com/ej-agas/perfume-db/nanoid"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestFactory_NewNoteGroup(t *testing.T) {
	alphabet := "0123456789abcdefghijklmnopqrstuvwxyz"
	length := 12
	factory := Factory{IdGenerator: nanoid.NewNanoIdGenerator(alphabet, length)}

	Name := "Fruits, Vegetables And Nuts"
	Slug := "fruits-vegetables-and-nuts"
	Description := "Foo"
	ImageURL := "/images/fruits-vegtables-and-nuts.png"

	noteGroup, err := factory.NewNoteGroup(Name, Description, ImageURL)
	assert.Nil(t, err)
	assert.Equal(t, length, len(noteGroup.PublicId))
	assert.Equal(t, Name, noteGroup.Name)
	assert.Equal(t, Slug, noteGroup.Slug)
	assert.Equal(t, Description, noteGroup.Description)
	assert.Equal(t, ImageURL, noteGroup.ImageURL)

	// Assert characters in PublicId are within the alphabet
	for _, char := range noteGroup.PublicId {
		assert.Contains(t, alphabet, string(char))
	}
}

func TestFactory_NewNote(t *testing.T) {
	alphabet := "0123456789abcdefghijklmnopqrstuvwxyz"
	length := 12
	factory := Factory{IdGenerator: nanoid.NewNanoIdGenerator(alphabet, length)}

	Name := "Citrus"
	Slug := "citrus"
	Description := "Citrus Note"
	ImageURL := "/images/citrus-note.png"
	NoteGroupId, _ := factory.IdGenerator.Generate()

	note, err := factory.NewNote(Name, Description, ImageURL, NoteGroupId)

	assert.Nil(t, err)
	assert.Equal(t, length, len(note.PublicId))
	assert.Equal(t, Name, note.Name)
	assert.Equal(t, Slug, note.Slug)
	assert.Equal(t, Description, note.Description)
	assert.Equal(t, ImageURL, note.ImageURL)
	assert.Equal(t, NoteGroupId, note.NoteGroupId)

	// Assert characters in PublicId are within the alphabet
	for _, char := range note.PublicId {
		assert.Contains(t, alphabet, string(char))
	}
}

func TestFactory_NewHouse(t *testing.T) {
	alphabet := "0123456789abcdefghijklmnopqrstuvwxyz"
	length := 12
	factory := Factory{IdGenerator: nanoid.NewNanoIdGenerator(alphabet, length)}

	slug := "perfume-company-123"
	name := "Perfume Company 123"
	country := "Philippines"
	description := "Niche Perfume House"
	yearFounded := time.Now()

	house, err := factory.NewHouse(name, country, description, yearFounded)

	assert.Nil(t, err)
	assert.Equal(t, length, len(house.PublicId))
	assert.Equal(t, name, house.Name)
	assert.Equal(t, country, house.Country)
	assert.Equal(t, description, house.Description)
	assert.Equal(t, yearFounded, house.YearFounded)
	assert.Equal(t, slug, house.Slug)

	// Assert characters in PublicId are within the alphabet
	for _, char := range house.PublicId {
		assert.Contains(t, alphabet, string(char))
	}
}
