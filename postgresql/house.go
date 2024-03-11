package postgresql

import (
	"context"
	"github.com/ej-agas/perfume-db/internal"
	"github.com/jackc/pgx/v5"
)

type HouseService struct {
	DB *pgx.Conn
}

func (h HouseService) Save(house internal.House) error {
	q := `
		INSERT INTO houses (slug, name, country, description, year_founded, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := h.DB.Exec(
		context.Background(),
		q,
		house.Slug, house.Name, house.Country, house.Description, house.YearFounded, house.CreatedAt, house.UpdatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}

func (h HouseService) Find(id int) (*internal.House, error) {
	return nil, nil
}

func (h HouseService) FindBySlug(s string) (*internal.House, error) {
	return nil, nil
}
