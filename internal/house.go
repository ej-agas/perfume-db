package internal

import (
	"time"
)

type House struct {
	ID          int
	Slug        string
	Name        string
	Country     string
	Description string
	YearFounded time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func NewHouse(name string, country string, description string, yearFounded time.Time) *House {
	now := time.Now()
	return &House{
		Slug:        CreateSlug(name),
		Name:        name,
		Country:     country,
		Description: description,
		YearFounded: yearFounded,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

type HouseService interface {
	Save(house House) error
	Find(id int) (*House, error)
	FindBySlug(s string) (*House, error)
}
