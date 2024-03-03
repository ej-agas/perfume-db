package internal

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNewPerfumer(t *testing.T) {
	Slug := "john-mcdoe-doe-iii"
	Name := "John McDoe Doe III"
	BirthDate := time.Date(1999, time.January, 20, 0, 0, 0, 0, time.UTC)
	PhotoURL := "/images/john-mcdoe-doe-iii.png"

	perfumer := NewPerfumer(Name, BirthDate, PhotoURL)

	assert.Equal(t, Slug, perfumer.Slug)
	assert.Equal(t, Name, perfumer.Name)
	assert.Equal(t, PhotoURL, perfumer.PhotoURL)
	assert.Equal(t, BirthDate, perfumer.BirthDate)
}
