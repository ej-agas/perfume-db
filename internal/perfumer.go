package internal

import (
	"encoding/json"
	"time"
)

type Perfumer struct {
	ID          int       `json:"-"`
	PublicId    string    `json:"id"`
	Slug        string    `json:"slug"`
	Name        string    `json:"name"`
	Nationality string    `json:"nationality"`
	BirthDate   time.Time `json:"birth_date"`
	ImageURL    string    `json:"image_url"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (p Perfumer) GetID() int {
	return p.ID
}

func (p Perfumer) MarshalJSON() ([]byte, error) {
	type Alias Perfumer

	return json.Marshal(&struct {
		*Alias
		BirthDate string `json:"birth_date"`
	}{
		Alias:     (*Alias)(&p),
		BirthDate: p.BirthDate.Format("January 2, 2006"),
	})
}

func NewPerfumer(Name, Nationality, photoURL string, birthDate time.Time) *Perfumer {
	now := time.Now()
	return &Perfumer{
		Slug:        CreateSlug(Name),
		Name:        Name,
		Nationality: Nationality,
		BirthDate:   birthDate,
		ImageURL:    photoURL,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

type PerfumerService interface {
	List(cursor, perPage int) ([]Perfumer, error)
	Save(note *Perfumer) error
	Find(publicId string) (*Perfumer, error)
	FindBySlug(s string) (*Perfumer, error)
}
