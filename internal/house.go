package internal

import (
	"encoding/json"
	"github.com/google/uuid"
	"time"
)

type House struct {
	ID          string    `json:"id"`
	Slug        string    `json:"slug"`
	Name        string    `json:"name"`
	Country     string    `json:"country"`
	Description string    `json:"description"`
	ImageURL    string    `json:"image_url"`
	FoundedAt   time.Time `json:"founded_at"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// NewUUID generates a new UUID string
func NewUUID() string {
	return uuid.New().String()
}

// IsValidUUID checks if a string is a valid UUID
func IsValidUUID(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
}

func (h House) GetID() string {
	return h.ID
}

func NewHouse(name string, country string, description string, foundedAt time.Time) *House {
	now := time.Now()
	return &House{
		Slug:        CreateSlug(name),
		Name:        name,
		Country:     country,
		Description: description,
		FoundedAt:   foundedAt,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

func (h House) MarshalJSON() ([]byte, error) {
	type Alias House

	return json.Marshal(&struct {
		*Alias
		FoundedAt string `json:"year_founded"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
	}{
		Alias:     (*Alias)(&h),
		FoundedAt: h.FoundedAt.Format("2006-01-02T15:04:05Z07:00"),
		CreatedAt: h.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: h.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	})
}

type HouseService interface {
	List(cursor, perPage int) ([]House, error)
	Save(house *House) error
	Find(publicId string) (*House, error)
	FindBySlug(s string) (*House, error)
}
