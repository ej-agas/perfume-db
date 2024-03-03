package internal

import "time"

type Note struct {
	ID          int
	Name        string
	Slug        string
	Description string
	ImageURL    string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func NewNote(Name, Description, ImageURL string) *Note {
	now := time.Now()
	return &Note{
		Name:        Name,
		Slug:        CreateSlug(Name),
		Description: Description,
		ImageURL:    ImageURL,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}
