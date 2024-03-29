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

type NoteService struct {
	db *pgx.Conn
}

var ErrNoteAlreadyExists = fmt.Errorf("note already exists")

func (service NoteService) List(cursor, perPage int) ([]internal.Note, error) {
	q := `SELECT * FROM notes WHERE id > $1 ORDER BY id LIMIT $2`
	if cursor <= 0 {
		cursor = 0
	}

	rows, err := service.db.Query(context.Background(), q, cursor, perPage)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	var notes []internal.Note

	for rows.Next() {
		var note internal.Note
		err := rows.Scan(
			&note.ID,
			&note.Slug,
			&note.Name,
			&note.Description,
			&note.ImageURL,
			&note.CreatedAt,
			&note.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		notes = append(notes, note)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return notes, nil
}

func (service NoteService) Save(note *internal.Note) error {
	if note.ID == 0 {
		return service.saveNewNote(note)
	}

	return service.updateNote(note)
}

func (service NoteService) saveNewNote(note *internal.Note) error {
	q := `
		INSERT INTO notes (public_id, slug, name, description, image_url, note_group_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err := service.db.Exec(
		context.Background(),
		q,
		note.PublicId,
		note.Slug,
		note.Name,
		note.Description,
		note.ImageURL,
		note.NoteGroupId,
		note.CreatedAt,
		note.UpdatedAt,
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
		return fmt.Errorf("database error: %w: %w", ErrNoteAlreadyExists, pgErr)
	}

	return err
}

func (service NoteService) updateNote(note *internal.Note) error {
	q := `
		UPDATE notes 
		SET slug = $2,
		    name = $3,
		    description = $4,
		    updated_at = $5
		WHERE id = $1
	`

	note.UpdatedAt = time.Now()
	_, err := service.db.Exec(context.Background(),
		q,
		note.ID,
		note.Slug,
		note.Name,
		note.Description,
		note.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("database error: update note error: %w", err)
	}

	return nil
}

func (service NoteService) Find(id int) (*internal.Note, error) {
	var note internal.Note

	q := `SELECT * FROM notes WHERE id = $1`

	if err := service.db.QueryRow(context.Background(), q, id).
		Scan(
			&note.ID,
			&note.Name,
			&note.Slug,
			&note.Description,
			&note.ImageURL,
			&note.CreatedAt,
			&note.UpdatedAt,
		); err != nil {
		return nil, err
	}

	return &note, nil
}

func (service NoteService) FindBySlug(s string) (*internal.Note, error) {
	var note internal.Note

	q := `SELECT * FROM notes WHERE slug = $1`

	if err := service.db.QueryRow(context.Background(), q, s).
		Scan(
			&note.ID,
			&note.PublicId,
			&note.Name,
			&note.Slug,
			&note.Description,
			&note.ImageURL,
			&note.NoteGroupId,
			&note.CreatedAt,
			&note.UpdatedAt,
		); err != nil {
		return nil, err
	}

	return &note, nil
}
