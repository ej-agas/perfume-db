package internal

import "time"

type Note struct {
	ID          int       `json:"-"`
	PublicId    string    `json:"id"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	Description string    `json:"description"`
	ImageURL    string    `json:"image_url"`
	NoteGroupId string    `json:"note_group_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func NewNote(Name, Description, ImageURL, NoteGroupId string) *Note {
	now := time.Now()
	return &Note{
		Name:        Name,
		Slug:        CreateSlug(Name),
		Description: Description,
		ImageURL:    ImageURL,
		NoteGroupId: NoteGroupId,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

type NoteService interface {
	List(cursor, perPage int) ([]Note, error)
	Save(note *Note) error
	Find(id int) (*Note, error)
	FindBySlug(s string) (*Note, error)
}
