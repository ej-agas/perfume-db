package internal

import "time"

type NoteGroup struct {
	ID          int       `json:"-"`
	PublicId    string    `json:"id"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	Description string    `json:"description"`
	ImageURL    string    `json:"image_url"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
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
