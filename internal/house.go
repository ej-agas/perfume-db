package internal

import (
	"encoding/json"
	"strconv"
	"time"
)

type House struct {
	ID          int       `json:"-"`
	PublicId    string    `json:"id"`
	Slug        string    `json:"slug"`
	Name        string    `json:"name"`
	Country     string    `json:"country"`
	Description string    `json:"description"`
	YearFounded time.Time `json:"year_founded"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (h House) GetID() int {
	return h.ID
}

func NewHouse(name string, country string, description string, yearFounded time.Time) *House {
	now := time.Now()
	return &House{
		Slug:        CreateSlug(name),
		Name:        name,
		Country:     country,
		Description: description,
		YearFounded: yearFounded,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

func (h House) MarshalJSON() ([]byte, error) {
	type Alias House

	year, err := strconv.Atoi(h.YearFounded.Format("2006"))

	if err != nil {
		return nil, err
	}

	return json.Marshal(&struct {
		*Alias
		YearFounded int    `json:"year_founded"`
		CreatedAt   string `json:"created_at"`
		UpdatedAt   string `json:"updated_at"`
	}{
		Alias:       (*Alias)(&h),
		YearFounded: year,
		CreatedAt:   h.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   h.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	})
}

type HouseService interface {
	List(cursor, perPage int) ([]House, error)
	Save(house *House) error
	Find(id int) (*House, error)
	FindBySlug(s string) (*House, error)
}
