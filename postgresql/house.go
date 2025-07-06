package postgresql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/ej-agas/perfume-db/internal"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type HouseService struct {
	db *pgxpool.Pool
}

var (
	ErrHouseAlreadyExists = fmt.Errorf("error house already exists")
	ErrHouseNotFound      = fmt.Errorf("house not found")
)

func (service HouseService) List(cursor, perPage int) ([]internal.House, error) {
	q := `SELECT id, slug, name, country, description, COALESCE(image_url, '') as image_url, founded_at, created_at, updated_at FROM houses ORDER BY id LIMIT $1`

	rows, err := service.db.Query(context.Background(), q, perPage)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	var houses []internal.House
	for rows.Next() {
		var (
			house  internal.House
			pgUUID pgtype.UUID
		)
		err := rows.Scan(
			&pgUUID,
			&house.Slug,
			&house.Name,
			&house.Country,
			&house.Description,
			&house.ImageURL,
			&house.FoundedAt,
			&house.CreatedAt,
			&house.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		// Convert pgtype.UUID to string
		if !pgUUID.Valid {
			return nil, fmt.Errorf("invalid UUID value")
		}
		// Convert [16]byte to uuid.UUID and then to string
		uid, err := uuid.FromBytes(pgUUID.Bytes[:])
		if err != nil {
			return nil, fmt.Errorf("failed to convert UUID: %w", err)
		}
		house.ID = uid.String()

		houses = append(houses, house)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return houses, nil
}

func (service HouseService) Save(house *internal.House) error {
	if house.ID == "" {
		return service.saveNewHouse(house)
	}

	return service.updateHouse(house)
}

func (service HouseService) saveNewHouse(house *internal.House) error {
	q := `
		INSERT INTO houses (id, slug, name, country, description, founded_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING created_at, updated_at
	`

	// Generate new UUID if not set
	if house.ID == "" {
		house.ID = internal.NewUUID()
	}

	// Parse UUID string to bytes
	uid, err := uuid.Parse(house.ID)
	if err != nil {
		return fmt.Errorf("invalid UUID format: %w", err)
	}

	row := service.db.QueryRow(
		context.Background(),
		q,
		uid,
		house.Slug,
		house.Name,
		house.Country,
		house.Description,
		house.FoundedAt,
		time.Now(),
		time.Now(),
	)

	err = row.Scan(
		&house.CreatedAt,
		&house.UpdatedAt,
	)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "23505": // unique_violation
				return fmt.Errorf("%w: %s", ErrHouseAlreadyExists, pgErr.Detail)
			}
		}
		return fmt.Errorf("failed to save house: %w", err)
	}

	return nil
}

func (service HouseService) updateHouse(house *internal.House) error {
	q := `
		UPDATE houses 
		SET slug = $2,
		    name = $3,
		    country = $4,
		    description = $5,
		    founded_at = $6,
		    updated_at = $7
		WHERE id = $1
	`

	// Parse UUID string to bytes
	uid, err := uuid.Parse(house.ID)
	if err != nil {
		return fmt.Errorf("invalid UUID format: %w", err)
	}

	house.UpdatedAt = time.Now()
	_, err = service.db.Exec(context.Background(),
		q,
		uid,
		house.Slug,
		house.Name,
		house.Country,
		house.Description,
		house.FoundedAt,
		house.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("update house error: %w", err)
	}

	return nil
}

func (service HouseService) Find(publicId string) (*internal.House, error) {
	q := "SELECT * FROM houses WHERE id = $1"

	var (
		house  internal.House
		pgUUID pgtype.UUID
	)

	err := service.db.QueryRow(context.Background(), q, publicId).Scan(
		&pgUUID,
		&house.Slug,
		&house.Name,
		&house.Country,
		&house.Description,
		&house.FoundedAt,
		&house.CreatedAt,
		&house.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to find house: %w", err)
	}

	// Convert [16]byte to uuid.UUID and then to string
	if !pgUUID.Valid {
		return nil, fmt.Errorf("invalid UUID value")
	}
	uid, err := uuid.FromBytes(pgUUID.Bytes[:])
	if err != nil {
		return nil, fmt.Errorf("failed to convert UUID: %w", err)
	}
	house.ID = uid.String()

	return &house, nil
}

func (service HouseService) FindBySlug(s string) (*internal.House, error) {
	q := `
        SELECT 
            id, 
            slug, 
            name, 
            country, 
            description, 
            COALESCE(image_url, '') as image_url, 
            founded_at, 
            created_at, 
            updated_at 
        FROM houses 
        WHERE slug = $1
    `

	var (
		house  internal.House
		pgUUID pgtype.UUID
	)

	err := service.db.QueryRow(context.Background(), q, s).Scan(
		&pgUUID,
		&house.Slug,
		&house.Name,
		&house.Country,
		&house.Description,
		&house.ImageURL,
		&house.FoundedAt,
		&house.CreatedAt,
		&house.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to find house by slug: %w", err)
	}

	// Convert [16]byte to uuid.UUID and then to string
	if !pgUUID.Valid {
		return nil, fmt.Errorf("invalid UUID value")
	}
	uid, err := uuid.FromBytes(pgUUID.Bytes[:])
	if err != nil {
		return nil, fmt.Errorf("failed to convert UUID: %w", err)
	}
	house.ID = uid.String()

	return &house, nil
}

func (service HouseService) FindPerfumesByHouse(house internal.House) (*[]*internal.Perfume, error) {
	q := `
		SELECT 
		    id,
		    slug,
		    name,
		    description,
		    concentration,
		    COALESCE(image_url, '') as image_url,
		    released_at,
		    discontinued_at,
		    created_at,
		    updated_at
		FROM perfumes 
		WHERE house_id = $1
		ORDER BY id DESC 
		LIMIT 25
	`
	perfumes := make([]*internal.Perfume, 0, 25) // Slice of pointers

	rows, err := service.db.Query(context.Background(), q, house.ID) // Use house.ID
	if err != nil {
		return &perfumes, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var (
			perfume          internal.Perfume
			pgUUID           pgtype.UUID
			yearDiscontinued sql.NullTime
		)
		err := rows.Scan(
			&pgUUID,
			&perfume.Slug,
			&perfume.Name,
			&perfume.Description,
			&perfume.Concentration,
			&perfume.ImageURL,
			&perfume.ReleasedAt,
			&yearDiscontinued,
			&perfume.CreatedAt,
			&perfume.UpdatedAt,
		)
		if err != nil {
			return &perfumes, fmt.Errorf("failed to scan row: %w", err)
		}

		if !pgUUID.Valid {
			return &perfumes, fmt.Errorf("invalid pgUUID value")
		}

		uuidStr, err := uuid.FromBytes(pgUUID.Bytes[:])
		if err != nil {
			return &perfumes, fmt.Errorf("failed to convert UUID: %w", err)
		}

		perfume.ID = uuidStr.String()
		perfume.House = &house

		if yearDiscontinued.Valid {
			perfume.YearDiscontinued = yearDiscontinued.Time
		}

		// Take address of the local variable and add to the slice
		perfumes = append(perfumes, &perfume)
	}

	return &perfumes, nil
}
