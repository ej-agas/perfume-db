package postgresql

import (
	"context"
	"errors"
	"fmt"
	"github.com/ej-agas/perfume-db/internal"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"time"
)

type HouseService struct {
	db *pgx.Conn
}

var ErrHouseAlreadyExists = fmt.Errorf("error house already exists")

func (service HouseService) List(cursor, perPage int) ([]internal.House, error) {
	q := `SELECT * FROM houses WHERE id > $1 ORDER BY id LIMIT $2`
	if cursor <= 0 {
		cursor = 0
	}

	rows, err := service.db.Query(context.Background(), q, cursor, perPage)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	var houses []internal.House
	for rows.Next() {
		var house internal.House
		err := rows.Scan(
			&house.ID,
			&house.PublicId,
			&house.Slug,
			&house.Name,
			&house.Country,
			&house.Description,
			&house.YearFounded,
			&house.CreatedAt,
			&house.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		houses = append(houses, house)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return houses, nil
}

func (service HouseService) Save(house *internal.House) error {
	if house.ID == 0 {
		return service.saveNewHouse(house)
	}

	return service.updateHouse(house)
}

func (service HouseService) saveNewHouse(house *internal.House) error {
	q := `
		INSERT INTO houses (public_id, slug, name, country, description, year_founded, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err := service.db.Exec(
		context.Background(),
		q,
		house.PublicId,
		house.Slug,
		house.Name,
		house.Country,
		house.Description,
		house.YearFounded,
		house.CreatedAt,
		house.UpdatedAt,
	)

	if err == nil {
		return nil
	}

	var pgErr *pgconn.PgError
	ok := errors.As(err, &pgErr)
	if !ok {
		return err
	}

	if pgErr.Code == "23505" {
		return fmt.Errorf("%w: %w", ErrHouseAlreadyExists, pgErr)
	}

	return err
}

func (service HouseService) updateHouse(house *internal.House) error {
	q := `
		UPDATE houses 
		SET slug = $2,
		    name = $3,
		    country = $4,
		    description = $5,
		    year_founded = $6,
		    updated_at = $7
		WHERE id = $1
	`

	house.UpdatedAt = time.Now()
	_, err := service.db.Exec(context.Background(),
		q,
		house.ID,
		house.Slug,
		house.Name,
		house.Country,
		house.Description,
		house.YearFounded,
		house.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("update house error: %w", err)
	}

	return nil
}

func (service HouseService) Find(publicId string) (*internal.House, error) {
	var house internal.House

	q := `SELECT * FROM houses WHERE public_id = $1`

	if err := service.db.QueryRow(context.Background(), q, publicId).
		Scan(
			&house.ID,
			&house.PublicId,
			&house.Slug,
			&house.Name,
			&house.Country,
			&house.Description,
			&house.YearFounded,
			&house.CreatedAt,
			&house.UpdatedAt,
		); err != nil {
		return nil, err
	}

	return &house, nil
}

func (service HouseService) FindBySlug(s string) (*internal.House, error) {
	var house internal.House

	q := `SELECT * FROM houses WHERE slug = $1`

	if err := service.db.QueryRow(context.Background(), q, s).
		Scan(
			&house.ID,
			&house.PublicId,
			&house.Slug,
			&house.Name,
			&house.Country,
			&house.Description,
			&house.YearFounded,
			&house.CreatedAt,
			&house.UpdatedAt,
		); err != nil {
		return nil, err
	}

	return &house, nil
}
