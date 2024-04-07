package internal

import (
	"time"
)

type Perfume struct {
	ID               int
	Slug             string
	Name             string
	Description      string
	Concentration    Concentration
	ImageURL         string
	House            *House
	Perfumers        []*Perfumer
	Notes            map[NoteCategory][]*Note
	YearReleased     time.Time
	YearDiscontinued time.Time
	CreatedAt        time.Time
	UpdatedAt        time.Time
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
