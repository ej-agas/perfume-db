package postgresql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/ej-agas/perfume-db/internal"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrPerfumeAlreadyExists = fmt.Errorf("perfume already exists")
)

type PerfumeService struct {
	db          *pgxpool.Pool
	noteService NoteService
}

func (service PerfumeService) Save(perfume *internal.Perfume) error {
	if perfume.ID == "" {
		return service.saveNewPerfume(perfume)
	}

	return service.updatePerfume(perfume)
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
		INSERT INTO perfumes (slug, name, description, concentration, image_url, house_id, released_at, discontinued_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`,
		perfume.Slug,
		perfume.Name,
		perfume.Description,
		perfume.Concentration,
		perfume.ImageURL,
		perfume.House.ID,
		perfume.ReleasedAt,
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
				perfume.ID,
				note.ID,
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
			perfume.ID,
			perfumer.ID,
		)
		if err != nil {
			return err
		}
	}

	err = tx.Commit(context.Background())

	return err
}

func (service PerfumeService) updatePerfume(perfume *internal.Perfume) error {
	conn, err := service.db.Acquire(context.Background())

	if err != nil {
		return fmt.Errorf("%w: %w", ErrAcquiringConn, err)
	}

	defer conn.Release()

	tx, err := conn.Begin(context.Background())

	if err != nil {
		return fmt.Errorf("%w: %w", ErrStartingDBTx, err)
	}

	defer tx.Rollback(context.Background())

	_, err = tx.Exec(
		context.Background(),
		`
		UPDATE perfumes SET
			 slug = $1,
			 name = $2,
			 description = $3,
			 concentration = $4,
			 image_url = $5,
			 house_id = $6,
			 released_at = $7,
			 discontinued_at = $8,
			 created_at = $9,
			 updated_at = $10
		WHERE id = $11
	`,
		perfume.Slug,
		perfume.Name,
		perfume.Description,
		perfume.Concentration,
		perfume.ImageURL,
		perfume.House.ID,
		perfume.ReleasedAt,
		service.convertToNullIfZeroValue(perfume.YearDiscontinued),
		perfume.CreatedAt,
		perfume.UpdatedAt,
		perfume.ID,
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
		placeholders := make([]string, len(notes))
		args := make([]interface{}, len(notes))
		for i, note := range notes {
			placeholders[i] = fmt.Sprintf("$%d", i+1)
			args[i] = note.ID
		}

		var q string
		if len(args) > 0 {
			q = fmt.Sprintf(
				`DELETE FROM perfumes_notes WHERE note_id NOT IN (%s) AND perfume_id = $%d AND category = $%d`,
				strings.Join(placeholders, ", "),
				len(args)+1,
				len(args)+2,
			)

			args = append(args, perfume.ID, category)
		} else {
			q = `DELETE FROM perfumes_notes WHERE perfume_id = $1 AND category = $2`
			args = append(args, perfume.ID, category)
		}

		if _, err := tx.Exec(context.Background(), q, args...); err != nil {
			return fmt.Errorf("update perfume error: delete from perfume_notes query error: %w", err)
		}

		for _, note := range notes {
			_, err = tx.Exec(
				context.Background(),
				`
	           INSERT INTO perfumes_notes (perfume_id, note_id, category)
	           VALUES ($1, $2, $3) ON CONFLICT (perfume_id, note_id) DO NOTHING 
	           `,
				perfume.ID,
				note.ID,
				category,
			)
			if err != nil {
				return fmt.Errorf("update perfume error: insert into perfume_notes query error: %w", err)
			}
		}
	}

	placeholders := make([]string, len(perfume.Perfumers))
	args := make([]interface{}, len(perfume.Perfumers))
	for i, perfumer := range perfume.Perfumers {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = perfumer.ID
	}

	var q string
	if len(args) > 0 {
		q = fmt.Sprintf(
			"DELETE FROM perfumes_perfumers WHERE perfumer_id NOT IN (%s) AND perfume_id = $%d",
			strings.Join(placeholders, ", "),
			len(args)+1,
		)
		args = append(args, perfume.ID)
	} else {
		q = `DELETE FROM perfumes_perfumers WHERE perfume_id = $1`
		args = append(args, perfume.ID)
	}
	fmt.Println(q, args)
	if _, err := service.db.Exec(context.Background(), q, args...); err != nil {
		return fmt.Errorf("update perfume error: delete from perfume_perfumers query error: %w", err)
	}

	for _, perfumer := range perfume.Perfumers {
		_, err = tx.Exec(
			context.Background(),
			`
			INSERT INTO perfumes_perfumers (perfume_id, perfumer_id)
			VALUES ($1, $2) ON CONFLICT (perfume_id, perfumer_id) DO NOTHING
			`,
			perfume.ID,
			perfumer.ID,
		)
		if err != nil {
			return fmt.Errorf("update perfume error: insert into perfume_perfumers query error: %w", err)
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

func (service PerfumeService) Find(publicId string) (*internal.Perfume, error) {
	var perfume internal.Perfume
	perfume.House = &internal.House{}

	var yearDiscontinued sql.NullTime

	perfumeQuery := `
        SELECT p.id, 
               p.slug, 
               p.name, 
               p.description, 
               p.concentration, 
               p.image_url, 
               p.released_at, 
               p.discontinued_at, 
               p.created_at, 
               p.updated_at,
			   p.house_id,
               h.slug AS house_slug,
               h.name AS house_name,
               h.country AS house_country,
               h.description AS house_description,
               h.founded_at AS house_year_founded,
               h.created_at AS house_created_at,
               h.updated_at AS house_updated_at
        FROM perfumes p
        LEFT JOIN houses h ON p.house_id = h.id
        WHERE p.id = $1
	`

	row := service.db.QueryRow(context.Background(), perfumeQuery, publicId)
	err := row.Scan(
		&perfume.ID,
		&perfume.Slug,
		&perfume.Name,
		&perfume.Description,
		&perfume.Concentration,
		&perfume.ImageURL,
		&perfume.ReleasedAt,
		&yearDiscontinued,
		&perfume.CreatedAt,
		&perfume.UpdatedAt,
		&perfume.House.ID,
		&perfume.House.Slug,
		&perfume.House.Name,
		&perfume.House.Country,
		&perfume.House.Description,
		&perfume.House.FoundedAt,
		&perfume.House.CreatedAt,
		&perfume.House.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("perfume with public_id '%s' not found", publicId)
		}
		return nil, err
	}

	if yearDiscontinued.Valid {
		perfume.YearDiscontinued = yearDiscontinued.Time
	}

	perfumersQuery := `
		SELECT
			p.id,
			p.slug,
			p.name,
			p.nationality,
			p.image_url,
			p.birth_date,
			p.created_at,
			p.updated_at
		FROM perfumes_perfumers
				 LEFT JOIN perfumers p ON perfumes_perfumers.perfumer_id = p.id
		WHERE perfume_id = $1;
`
	perfumerRows, err := service.db.Query(context.Background(), perfumersQuery, perfume.ID)
	if err != nil {
		return nil, err
	}
	defer perfumerRows.Close()

	for perfumerRows.Next() {
		var perfumer internal.Perfumer
		err := perfumerRows.Scan(
			&perfumer.ID,
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
		   n.slug,
		   n.name,
		   n.description,
		   n.image_url,
		   n.note_group_id
		FROM perfumes_notes
				 LEFT JOIN notes n ON perfumes_notes.note_id = n.id
		WHERE perfume_id = $1;
`
	noteRows, err := service.db.Query(context.Background(), notesQuery, perfume.ID)
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

func (service PerfumeService) FindBySlug(slug string) (*internal.Perfume, error) {
	var perfume internal.Perfume
	perfume.House = &internal.House{}

	var yearDiscontinued sql.NullTime

	perfumeQuery := `
        SELECT p.id, 
               p.slug, 
               p.name, 
               p.description, 
               p.concentration, 
               p.image_url, 
               p.released_at, 
               p.discontinued_at, 
               p.created_at, 
               p.updated_at,
			   p.house_id,
               h.slug AS house_slug,
               h.name AS house_name,
               h.country AS house_country,
               h.description AS house_description,
               h.founded_at AS house_year_founded,
               h.created_at AS house_created_at,
               h.updated_at AS house_updated_at
        FROM perfumes p
        LEFT JOIN houses h ON p.house_id = h.id
        WHERE p.slug = $1
	`

	conn, err := service.db.Acquire(context.Background())
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrAcquiringConn, err)
	}

	defer conn.Release()

	row := conn.QueryRow(context.Background(), perfumeQuery, slug)
	err = row.Scan(
		&perfume.ID,
		&perfume.Slug,
		&perfume.Name,
		&perfume.Description,
		&perfume.Concentration,
		&perfume.ImageURL,
		&perfume.ReleasedAt,
		&yearDiscontinued,
		&perfume.CreatedAt,
		&perfume.UpdatedAt,
		&perfume.House.ID,
		&perfume.House.Slug,
		&perfume.House.Name,
		&perfume.House.Country,
		&perfume.House.Description,
		&perfume.House.FoundedAt,
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
        p.slug,
        p.name,
        p.nationality,
        COALESCE(p.image_url, '') as image_url,
        p.birth_date,
        p.created_at,
        p.updated_at
    FROM perfumes_perfumers
    LEFT JOIN perfumers p ON perfumes_perfumers.perfumer_id = p.id
    WHERE perfume_id = $1;
`
	perfumerRows, err := service.db.Query(context.Background(), perfumersQuery, perfume.ID)
	if err != nil {
		return nil, fmt.Errorf("error querying perfumers: %w", err)
	}
	defer perfumerRows.Close()

	var perfumers []*internal.Perfumer
	for perfumerRows.Next() {
		var perfumer internal.Perfumer
		var birthDate sql.NullTime

		err := perfumerRows.Scan(
			&perfumer.ID,
			&perfumer.Slug,
			&perfumer.Name,
			&perfumer.Nationality,
			&perfumer.ImageURL,
			&birthDate,
			&perfumer.CreatedAt,
			&perfumer.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning perfumer: %w", err)
		}

		if birthDate.Valid {
			perfumer.BirthDate = birthDate.Time
		}

		perfumers = append(perfumers, &perfumer)
	}

	if err = perfumerRows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating perfumers: %w", err)
	}

	perfume.Perfumers = perfumers

	notesQuery := `
		SELECT
		   category,
		   n.id,
		   n.slug,
		   n.name,
		   n.description,
		   n.image_url,
		   n.note_group_id
		FROM perfumes_notes
				 LEFT JOIN notes n ON perfumes_notes.note_id = n.id
		WHERE perfume_id = $1;
`
	noteRows, err := service.db.Query(context.Background(), notesQuery, perfume.ID)
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

func (service PerfumeService) FindByPerfumer(perfumerID string) ([]*internal.Perfume, error) {
	query := `
        SELECT 
            p.id, 
            p.slug, 
            p.name, 
            p.description, 
            p.concentration, 
            COALESCE(p.image_url, '') as image_url,
            p.released_at, 
            p.discontinued_at, 
            p.created_at, 
            p.updated_at,
            p.house_id,
            h.slug AS house_slug,
            h.name AS house_name,
            h.country AS house_country,
            COALESCE(h.image_url, '') AS house_image_url
        FROM perfumes p
        JOIN houses h ON p.house_id = h.id
        JOIN perfumes_perfumers pp ON p.id = pp.perfume_id
        WHERE pp.perfumer_id = $1
        ORDER BY p.name
    `

	rows, err := service.db.Query(context.Background(), query, perfumerID)
	if err != nil {
		return nil, fmt.Errorf("error querying perfumes by perfumer: %w", err)
	}
	defer rows.Close()

	var perfumes []*internal.Perfume
	for rows.Next() {
		var p internal.Perfume
		var yearDiscontinued sql.NullTime
		var house internal.House

		err := rows.Scan(
			&p.ID,
			&p.Slug,
			&p.Name,
			&p.Description,
			&p.Concentration,
			&p.ImageURL,
			&p.ReleasedAt,
			&yearDiscontinued,
			&p.CreatedAt,
			&p.UpdatedAt,
			&house.ID,
			&house.Slug,
			&house.Name,
			&house.Country,
			&house.ImageURL,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning perfume row: %w", err)
		}

		if yearDiscontinued.Valid {
			p.YearDiscontinued = yearDiscontinued.Time
		}

		p.House = &house
		p.Notes = make(map[internal.NoteCategory][]*internal.Note)
		perfumes = append(perfumes, &p)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating perfume rows: %w", err)
	}

	return perfumes, nil
}

//func (service PerfumeService) FindBySlug(s string) (*Perfume, error)           {}
//func (service PerfumeService) FindMany(publicIds []string) ([]*Perfume, error) {}
