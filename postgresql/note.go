package postgresql

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/ej-agas/perfume-db/internal"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type NoteService struct {
	db *pgxpool.Pool
}

var (
	ErrNoteAlreadyExists = fmt.Errorf("note already exists")
	ErrNoteGroupNotFound = fmt.Errorf("note group not found")
	ErrNoteNotFound      = fmt.Errorf("note not found")
)

func (service NoteService) List(cursor, perPage int) ([]internal.Note, error) {
	if cursor <= 0 {
		cursor = 0
	}

	q := `SELECT * FROM notes WHERE id > $1 ORDER BY id LIMIT $2`
	rows, err := service.db.Query(context.Background(), q, cursor, perPage)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	var notes []internal.Note

	for rows.Next() {
		var note internal.Note
		if err := rows.Scan(
			&note.ID,
			&note.PublicId,
			&note.Slug,
			&note.Name,
			&note.Description,
			&note.ImageURL,
			&note.NoteGroupId,
			&note.CreatedAt,
			&note.UpdatedAt,
		); err != nil {
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

	switch pgErr.Code {
	case "23505":
		return fmt.Errorf("database error: %w: %w", ErrNoteAlreadyExists, pgErr)
	case "23503":
		return fmt.Errorf("database error: %w: %w", ErrNoteGroupNotFound, pgErr)
	default:
		return err
	}
}

func (service NoteService) updateNote(note *internal.Note) error {
	q := `
		UPDATE notes 
		SET slug = $2,
		    name = $3,
		    description = $4,
		    note_group_id = $5,
		    updated_at = $6
		WHERE id = $1
	`

	note.UpdatedAt = time.Now()
	_, err := service.db.Exec(context.Background(),
		q,
		note.ID,
		note.Slug,
		note.Name,
		note.Description,
		note.NoteGroupId,
		note.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("database error: update note error: %w", err)
	}

	return nil
}

func (service NoteService) Find(publicId string) (*internal.Note, error) {
	var note internal.Note

	q := `SELECT * FROM notes WHERE public_id = $1`

	if err := service.db.QueryRow(context.Background(), q, publicId).
		Scan(
			&note.ID,
			&note.PublicId,
			&note.Slug,
			&note.Name,
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

func (service NoteService) FindBySlug(s string) (*internal.Note, error) {
	var note internal.Note

	q := `SELECT * FROM notes WHERE slug = $1`

	if err := service.db.QueryRow(context.Background(), q, s).
		Scan(
			&note.ID,
			&note.PublicId,
			&note.Slug,
			&note.Name,
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

func (service NoteService) FindMany(publicIds []string) ([]*internal.Note, error) {
	notes := make([]*internal.Note, 0)

	if len(publicIds) == 0 {
		return notes, nil
	}

	placeholders := make([]string, len(publicIds))
	args := make([]interface{}, len(publicIds))
	for i, id := range publicIds {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = id
	}
	q := fmt.Sprintf("SELECT * FROM notes WHERE public_id IN (%s)", strings.Join(placeholders, ", "))

	rows, err := service.db.Query(context.Background(), q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	found := make(map[string]bool)
	for _, id := range publicIds {
		found[id] = false
	}

	for rows.Next() {
		var note internal.Note
		if err := rows.Scan(
			&note.ID,
			&note.PublicId,
			&note.Slug,
			&note.Name,
			&note.Description,
			&note.ImageURL,
			&note.NoteGroupId,
			&note.CreatedAt,
			&note.UpdatedAt,
		); err != nil {
			return nil, err
		}

		found[note.PublicId] = true
		notes = append(notes, &note)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	for id, ok := range found {
		if !ok {
			return nil, fmt.Errorf("%w: note with public_id '%s' not found", ErrNoteNotFound, id)
		}
	}

	return notes, nil
}
