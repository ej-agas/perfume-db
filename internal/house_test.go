package internal

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNewHouse(t *testing.T) {
	slug := "perfume-company-123"
	name := "Perfume Company 123"
	country := "Philippines"
	description := "Niche Perfume House"
	yearFounded := time.Now()

	company := NewHouse(name, country, description, yearFounded)

	assert.Equal(t, name, company.Name)
	assert.Equal(t, country, company.Country)
	assert.Equal(t, description, company.Description)
	assert.Equal(t, yearFounded, company.YearFounded)
	assert.Equal(t, slug, company.Slug)
}
