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
	House            House
	Perfumers        []Perfumer
	Notes            []Note
	YearReleased     time.Time
	YearDiscontinued time.Time
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

func NewPerfume(
	Name, Description, ImageURL string,
	Concentration Concentration,
	House House,
	YearReleased, YearDiscontinued time.Time,
) *Perfume {
	now := time.Now()
	return &Perfume{
		Slug:             CreateSlug(Name + "-" + Concentration.String()),
		Name:             Name,
		ImageURL:         ImageURL,
		Description:      Description,
		Concentration:    Concentration,
		House:            House,
		YearReleased:     YearReleased,
		YearDiscontinued: YearDiscontinued,
		CreatedAt:        now,
		UpdatedAt:        now,
	}
}
