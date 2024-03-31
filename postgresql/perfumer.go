package postgresql

import (
	"context"
	"errors"
	"fmt"
	"github.com/ej-agas/perfume-db/internal"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

var (
	ErrPerfumerAlreadyExists = fmt.Errorf("perfumer already exists")
)

type PerfumerService struct {
	db *pgx.Conn
}

func (service PerfumerService) List(cursor, perPage int) ([]internal.Perfumer, error) {
	if cursor <= 0 {
		cursor = 0
	}

	q := `SELECT * FROM perfumers WHERE id > $1 ORDER BY id LIMIT $2`
	rows, err := service.db.Query(context.Background(), q, cursor, perPage)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	var perfumers []internal.Perfumer

	for rows.Next() {
		var perfumer internal.Perfumer
		if err := rows.Scan(
			&perfumer.ID,
			&perfumer.PublicId,
			&perfumer.Slug,
			&perfumer.Name,
			&perfumer.Nationality,
			&perfumer.ImageURL,
			&perfumer.BirthDate,
			&perfumer.CreatedAt,
			&perfumer.UpdatedAt,
		); err != nil {
			return nil, err
		}

		perfumers = append(perfumers, perfumer)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return perfumers, nil
}

func (service PerfumerService) Save(perfumer *internal.Perfumer) error {
	if perfumer.ID == 0 {
		return service.saveNewPerfumer(perfumer)
	}

	return service.updatePerfumer(perfumer)
}

func (service PerfumerService) saveNewPerfumer(perfumer *internal.Perfumer) error {
	q := `
		INSERT INTO perfumers (public_id, slug, name, nationality, image_url, birth_date, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err := service.db.Exec(
		context.Background(),
		q,
		perfumer.PublicId,
		perfumer.Slug,
		perfumer.Name,
		perfumer.Nationality,
		perfumer.ImageURL,
		perfumer.BirthDate,
		perfumer.CreatedAt,
		perfumer.UpdatedAt,
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
		return fmt.Errorf("%w: %w", ErrPerfumerAlreadyExists, pgErr)
	}

	return err
}

func (service PerfumerService) updatePerfumer(perfumer *internal.Perfumer) error {
	q := `
		UPDATE perfumers 
		SET slug = $2,
		    name = $3,
		    nationality = $4,
		    image_url = $5,
		    birth_date = $6,
		    updated_at = $7
		WHERE id = $1
	`

	_, err := service.db.Exec(
		context.Background(),
		q,
		perfumer.ID,
		perfumer.Slug,
		perfumer.Name,
		perfumer.Nationality,
		perfumer.ImageURL,
		perfumer.BirthDate,
		perfumer.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("update perfumer error: %w", err)
	}

	return nil
}

func (service PerfumerService) Find(publicId string) (*internal.Perfumer, error) {
	var perfumer internal.Perfumer

	q := `SELECT * FROM perfumers WHERE public_id = $1`

	if err := service.db.QueryRow(context.Background(), q, publicId).
		Scan(
			&perfumer.ID,
			&perfumer.PublicId,
			&perfumer.Slug,
			&perfumer.Name,
			&perfumer.Nationality,
			&perfumer.ImageURL,
			&perfumer.BirthDate,
			&perfumer.CreatedAt,
			&perfumer.UpdatedAt,
		); err != nil {
		return nil, err
	}

	return &perfumer, nil
}

func (service PerfumerService) FindBySlug(s string) (*internal.Perfumer, error) {
	var perfumer internal.Perfumer

	q := `SELECT * FROM perfumers WHERE slug = $1`

	if err := service.db.QueryRow(context.Background(), q, s).
		Scan(
			&perfumer.ID,
			&perfumer.PublicId,
			&perfumer.Slug,
			&perfumer.Name,
			&perfumer.Nationality,
			&perfumer.ImageURL,
			&perfumer.BirthDate,
			&perfumer.CreatedAt,
			&perfumer.UpdatedAt,
		); err != nil {
		return nil, err
	}

	return &perfumer, nil
}
