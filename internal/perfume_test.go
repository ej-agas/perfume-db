package internal

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNewPerfume(t *testing.T) {
	Slug := "perfume-abc-eau-de-parfum"
	Name := "Perfume ABC"
	concentration := EauDeParfum
	Description := "Foo"
	ImageURL := "/images/perfume-abc-edt.png"
	YearReleased := time.Now()
	YearDiscontinued := time.Now()

	house := House{ID: 1000}
	perfume := NewPerfume(Name, Description, ImageURL, concentration, house, YearReleased, YearDiscontinued)

	assert.Equal(t, Slug, perfume.Slug)
	assert.Equal(t, Name, perfume.Name)
	assert.Equal(t, Description, perfume.Description)
	assert.Equal(t, concentration, perfume.Concentration)
	assert.Equal(t, ImageURL, perfume.ImageURL)
	assert.Equal(t, house, perfume.House)
	assert.Equal(t, YearReleased, perfume.YearReleased)
	assert.Equal(t, YearDiscontinued, perfume.YearDiscontinued)
}
