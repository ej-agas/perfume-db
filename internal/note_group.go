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

func (n NoteGroup) GetID() int {
	return n.ID
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

type NoteGroupService interface {
	List(cursor, perPage int) ([]NoteGroup, error)
	Save(noteGroup *NoteGroup) error
	Find(publicId string) (*NoteGroup, error)
	FindBySlug(s string) (*NoteGroup, error)
}
