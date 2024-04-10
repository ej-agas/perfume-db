package internal

import (
	"time"
)

type Perfume struct {
	ID               int                      `json:"-"`
	PublicId         string                   `json:"id"`
	Slug             string                   `json:"slug"`
	Name             string                   `json:"name"`
	Description      string                   `json:"description"`
	Concentration    Concentration            `json:"concentration"`
	ImageURL         string                   `json:"image_url"`
	House            *House                   `json:"house"`
	Perfumers        []*Perfumer              `json:"perfumers"`
	Notes            map[NoteCategory][]*Note `json:"notes"`
	YearReleased     time.Time                `json:"year_released"`
	YearDiscontinued time.Time                `json:"year_discontinued"`
	CreatedAt        time.Time                `json:"created_at"`
	UpdatedAt        time.Time                `json:"updated_at"`
}

type PerfumeOption func(*Perfume)

func WithName(name string) PerfumeOption {
	return func(p *Perfume) {
		p.Name = name
	}
}

func WithDescription(description string) PerfumeOption {
	return func(p *Perfume) {
		p.Description = description
	}
}

func WithImageURL(imageURL string) PerfumeOption {
	return func(p *Perfume) {
		p.ImageURL = imageURL
	}
}

func WithConcentration(concentration Concentration) PerfumeOption {
	return func(p *Perfume) {
		p.Concentration = concentration
	}
}

func WithHouse(house *House) PerfumeOption {
	return func(p *Perfume) {
		p.House = house
	}
}

func WithPerfumers(perfumer ...*Perfumer) PerfumeOption {
	return func(p *Perfume) {
		p.Perfumers = perfumer
	}
}

func WithNotes(notes map[NoteCategory][]*Note) PerfumeOption {
	return func(p *Perfume) {
		p.Notes = notes
	}
}

func WithYearReleased(yearReleased time.Time) PerfumeOption {
	return func(p *Perfume) {
		p.YearReleased = yearReleased
	}
}

func WithYearDiscontinued(yearDiscontinued time.Time) PerfumeOption {
	return func(p *Perfume) {
		p.YearDiscontinued = yearDiscontinued
	}
}

type PerfumeService interface {
	List(cursor, perPage int) ([]*Perfume, error)
	Save(note *Perfume) error
	Find(publicId string) (*Perfume, error)
	FindBySlug(s string) (*Perfume, error)
	FindMany(publicIds []string) ([]*Perfume, error)
}
