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
		ID:          id,
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

	return &NoteGroup{
		Name:        Name,
		Slug:        CreateSlug(Name),
		Description: Description,
		ImageURL:    ImageURL,
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}

func (factory Factory) NewHouse(name, country, description string, foundedAt time.Time) (*House, error) {
	now := time.Now()

	return &House{
		Name:        name,
		Slug:        CreateSlug(name),
		Country:     country,
		Description: description,
		FoundedAt:   foundedAt,
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}

func (factory Factory) NewPerfumer(name, nationality, imageUrl string, birthDate time.Time) (*Perfumer, error) {
	now := time.Now()

	return &Perfumer{
		Slug:        CreateSlug(name),
		Name:        name,
		Nationality: nationality,
		ImageURL:    imageUrl,
		BirthDate:   birthDate,
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}

func (factory Factory) NewPerfume(opts ...PerfumeOption) (*Perfume, error) {
	now := time.Now()
	id, err := factory.IdGenerator.Generate()
	if err != nil {
		return &Perfume{}, err
	}

	p := &Perfume{
		ID:        id,
		CreatedAt: now,
		UpdatedAt: now,
	}

	for _, opt := range opts {
		opt(p)
	}

	p.Slug = CreateSlug(p.Name + "-" + p.Concentration.String())

	return p, nil
}
