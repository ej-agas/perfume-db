package internal

import (
	"time"
)

type Factory struct {
	IdGenerator IdGenerator
}

func (factory Factory) NewNote(Name, Description, ImageURL, NoteGroupId string) (*Note, error) {
	now := time.Now()
	id, err := factory.IdGenerator.Generate()
	if err != nil {
		return &Note{}, err
	}

	return &Note{
		PublicId:    id,
		Name:        Name,
		Slug:        CreateSlug(Name),
		Description: Description,
		ImageURL:    ImageURL,
		NoteGroupId: NoteGroupId,
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}

func (factory Factory) NewNoteGroup(Name, Description, ImageURL string) (*NoteGroup, error) {
	now := time.Now()
	id, err := factory.IdGenerator.Generate()
	if err != nil {
		return &NoteGroup{}, err
	}

	return &NoteGroup{
		PublicId:    id,
		Name:        Name,
		Slug:        CreateSlug(Name),
		Description: Description,
		ImageURL:    ImageURL,
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}

func (factory Factory) NewHouse(name string, country string, description string, yearFounded time.Time) (*House, error) {
	now := time.Now()
	id, err := factory.IdGenerator.Generate()
	if err != nil {
		return &House{}, err
	}

	return &House{
		PublicId:    id,
		Name:        name,
		Slug:        CreateSlug(name),
		Country:     country,
		Description: description,
		YearFounded: yearFounded,
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}
