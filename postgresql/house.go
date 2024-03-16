package postgresql

import (
	"context"
	"errors"
	"fmt"
	"github.com/ej-agas/perfume-db/internal"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type HouseService struct {
	db *pgx.Conn
}

var ErrHouseAlreadyExists error = fmt.Errorf("error house already exists")

func (service HouseService) Save(house *internal.House) error {
	q := `
		INSERT INTO houses (slug, name, country, description, year_founded, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := service.db.Exec(
		context.Background(),
		q,
		house.Slug, house.Name, house.Country, house.Description, house.YearFounded, house.CreatedAt, house.UpdatedAt,
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

	// Iterate through the rows and scan the results into House objects
	var houses []internal.House
	for rows.Next() {
		var house internal.House
		err := rows.Scan(
			&house.ID,
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

func (service HouseService) Find(id int) (*internal.House, error) {
	var house internal.House

	q := `SELECT * FROM houses WHERE id=$1`

	if err := service.db.QueryRow(context.Background(), q, id).
		Scan(
			&house.ID,
			&house.Slug,
			&house.Name,
			&house.Country,
			&house.Description,
			&house.YearFounded,
			&house.CreatedAt,
			&house.UpdatedAt,
		); err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &house, nil
}

func (service HouseService) FindBySlug(s string) (*internal.House, error) {
	var house internal.House

	q := `SELECT * FROM houses WHERE slug=$1`

	if err := service.db.QueryRow(context.Background(), q, s).
		Scan(
			&house.ID,
			&house.Slug,
			&house.Name,
			&house.Country,
			&house.Description,
			&house.YearFounded,
			&house.CreatedAt,
			&house.UpdatedAt,
		); err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &house, nil
}
