package postgresql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/ej-agas/perfume-db/internal"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"time"
)

var (
	ErrPerfumeAlreadyExists = fmt.Errorf("perfume already exists")
)

type PerfumeService struct {
	db          *pgx.Conn
	noteService NoteService
}

func (service PerfumeService) Save(perfume *internal.Perfume) error {
	if perfume.ID == 0 {
		return service.saveNewPerfume(perfume)
	}

	return nil
}

func (service PerfumeService) saveNewPerfume(perfume *internal.Perfume) error {
	tx, err := service.db.Begin(context.Background())

	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	_, err = tx.Exec(
		context.Background(),
		`
		INSERT INTO perfumes (public_id, slug, name, description, concentration, image_url, house_id, year_released, year_discontinued, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`,
		perfume.PublicId,
		perfume.Slug,
		perfume.Name,
		perfume.Description,
		perfume.Concentration,
		perfume.ImageURL,
		perfume.House.PublicId,
		perfume.YearReleased,
		service.convertToNullIfZeroValue(perfume.YearDiscontinued),
		perfume.CreatedAt,
		perfume.UpdatedAt,
	)

	if err != nil {
		var pgErr *pgconn.PgError
		ok := errors.As(err, &pgErr)
		if !ok {
			return err
		}

		switch pgErr.Code {
		case "23505":
			return fmt.Errorf("database error: %w: %w", ErrPerfumeAlreadyExists, pgErr)
		case "23503":
			return fmt.Errorf("database error: %w: %w", ErrHouseNotFound, pgErr)
		default:
			return err
		}
	}

	for category, notes := range perfume.Notes {
		for _, note := range notes {
			_, err = tx.Exec(
				context.Background(),
				`
                INSERT INTO perfumes_notes (perfume_id, note_id, category)
                VALUES ($1, $2, $3)
                `,
				perfume.PublicId,
				note.PublicId,
				category,
			)
			if err != nil {
				return err
			}
		}
	}

	for _, perfumer := range perfume.Perfumers {
		_, err = tx.Exec(
			context.Background(),
			`
			INSERT INTO perfumes_perfumers (perfume_id, perfumer_id)
			VALUES ($1, $2)
			`,
			perfume.PublicId,
			perfumer.PublicId,
		)
		if err != nil {
			return err
		}
	}

	err = tx.Commit(context.Background())

	return err
}

func (service PerfumeService) convertToNullIfZeroValue(t time.Time) sql.NullTime {
	if t.IsZero() {
		return sql.NullTime{Valid: false}
	}

	return sql.NullTime{Time: t, Valid: true}
}

func (service PerfumeService) FindBySlug(slug string) (*internal.Perfume, error) {
	var perfume internal.Perfume
	perfume.House = &internal.House{}

	var yearDiscontinued sql.NullTime

	perfumeQuery := `
        SELECT p.id, 
               p.public_id, 
               p.slug, 
               p.name, 
               p.description, 
               p.concentration, 
               p.image_url, 
               p.year_released, 
               p.year_discontinued, 
               p.created_at, 
               p.updated_at,
			   p.house_id,
               h.slug AS house_slug,
               h.name AS house_name,
               h.country AS house_country,
               h.description AS house_description,
               h.year_founded AS house_year_founded,
               h.created_at AS house_created_at,
               h.updated_at AS house_updated_at
        FROM perfumes p
        LEFT JOIN houses h ON p.house_id = h.public_id
        WHERE p.slug = $1
	`

	row := service.db.QueryRow(context.Background(), perfumeQuery, slug)
	err := row.Scan(
		&perfume.ID,
		&perfume.PublicId,
		&perfume.Slug,
		&perfume.Name,
		&perfume.Description,
		&perfume.Concentration,
		&perfume.ImageURL,
		&perfume.YearReleased,
		&yearDiscontinued,
		&perfume.CreatedAt,
		&perfume.UpdatedAt,
		&perfume.House.PublicId,
		&perfume.House.Slug,
		&perfume.House.Name,
		&perfume.House.Country,
		&perfume.House.Description,
		&perfume.House.YearFounded,
		&perfume.House.CreatedAt,
		&perfume.House.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("perfume with slug '%s' not found", slug)
		}
		return nil, err
	}

	if yearDiscontinued.Valid {
		perfume.YearDiscontinued = yearDiscontinued.Time
	}

	perfumersQuery := `
		SELECT
			p.id,
			p.public_id,
			p.slug,
			p.name,
			p.nationality,
			p.image_url,
			p.birth_date,
			p.created_at,
			p.updated_at
		FROM perfumes_perfumers
				 LEFT JOIN perfumers p ON perfumes_perfumers.perfumer_id = p.public_id
		WHERE perfume_id = $1;
`
	perfumerRows, err := service.db.Query(context.Background(), perfumersQuery, perfume.PublicId)
	if err != nil {
		return nil, err
	}
	defer perfumerRows.Close()

	for perfumerRows.Next() {
		var perfumer internal.Perfumer
		err := perfumerRows.Scan(
			&perfumer.ID,
			&perfumer.PublicId,
			&perfumer.Slug,
			&perfumer.Name,
			&perfumer.Nationality,
			&perfumer.ImageURL,
			&perfumer.BirthDate,
			&perfumer.CreatedAt,
			&perfumer.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		perfume.Perfumers = append(perfume.Perfumers, &perfumer)
	}

	notesQuery := `
		SELECT
		   category,
		   n.id,
		   n.public_id,
		   n.slug,
		   n.name,
		   n.description,
		   n.image_url,
		   n.note_group_id
		FROM perfumes_notes
				 LEFT JOIN notes n ON perfumes_notes.note_id = n.public_id
		WHERE perfume_id = $1;
`
	noteRows, err := service.db.Query(context.Background(), notesQuery, perfume.PublicId)
	if err != nil {
		return nil, err
	}
	defer noteRows.Close()

	// Initialize the map to store notes
	perfume.Notes = make(map[internal.NoteCategory][]*internal.Note)

	// Process the notes
	for noteRows.Next() {
		var note internal.Note
		var category string
		err := noteRows.Scan(
			&category,
			&note.ID,
			&note.PublicId,
			&note.Slug,
			&note.Name,
			&note.Description,
			&note.ImageURL,
			&note.NoteGroupId,
		)
		if err != nil {
			return nil, err
		}

		noteCategory, err := internal.NoteCategoryFromString(category)

		if err != nil {
			return nil, fmt.Errorf("error: invalid note category '%s': %w", category, err)
		}

		// Map the note to the appropriate category
		perfume.Notes[noteCategory] = append(perfume.Notes[noteCategory], &note)
	}

	return &perfume, nil
}

//func (service PerfumeService) FindBySlug(s string) (*Perfume, error)           {}
//func (service PerfumeService) FindMany(publicIds []string) ([]*Perfume, error) {}
