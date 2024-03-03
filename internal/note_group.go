package internal

import "time"

type NoteGroup struct {
	ID          int
	Name        string
	Slug        string
	Description string
	ImageURL    string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func NewNoteGroup(Name, Description, ImageURL string) *NoteGroup {
	now := time.Now()
	return &NoteGroup{
		Name:        Name,
		Slug:        CreateSlug(Name),
		Description: Description,
		ImageURL:    ImageURL,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}
