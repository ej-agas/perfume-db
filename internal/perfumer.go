package internal

import (
	"time"
)

type Perfumer struct {
	ID          int       `json:"id"`
	Slug        string    `json:"slug"`
	Name        string    `json:"name"`
	Nationality string    `json:"nationality"`
	BirthDate   time.Time `json:"birth_date"`
	PhotoURL    string    `json:"photo_url"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func NewPerfumer(Name, Nationality, photoURL string, birthDate time.Time) *Perfumer {
	now := time.Now()
	return &Perfumer{
		Slug:        CreateSlug(Name),
		Name:        Name,
		Nationality: Nationality,
		BirthDate:   birthDate,
		PhotoURL:    photoURL,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}
