package internal

import "time"

type Note struct {
	ID          string    `json:"id"`
	Slug        string    `json:"slug"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	ImageURL    string    `json:"image_url"`
	NoteGroupId string    `json:"note_group_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (n Note) GetID() string {
	return n.ID
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
	Find(publicId string) (*Note, error)
	FindBySlug(s string) (*Note, error)
	FindMany(publicIds []string) ([]*Note, error)
}
