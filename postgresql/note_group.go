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

type NoteGroupService struct {
	db *pgx.Conn
}

var ErrNoteGroupAlreadyExists = fmt.Errorf("note group already exists")

func (service NoteGroupService) List(cursor, perPage int) ([]internal.NoteGroup, error) {
	q := `SELECT * FROM note_groups WHERE id > $1 ORDER BY id LIMIT $2`
	if cursor <= 0 {
		cursor = 0
	}

	rows, err := service.db.Query(context.Background(), q, cursor, perPage)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	var noteGroups []internal.NoteGroup

	for rows.Next() {
		var note internal.NoteGroup
		err := rows.Scan(
			&note.ID,
			&note.PublicId,
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
		noteGroups = append(noteGroups, note)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return noteGroups, nil
}

func (service NoteGroupService) Save(note *internal.NoteGroup) error {
	if note.ID == 0 {
		return service.saveNewNoteGroup(note)
	}

	return service.updateNoteGroup(note)
}

func (service NoteGroupService) saveNewNoteGroup(noteGroup *internal.NoteGroup) error {
	q := `
		INSERT INTO note_groups (public_id, slug, name, description, image_url, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := service.db.Exec(
		context.Background(),
		q,
		noteGroup.PublicId,
		noteGroup.Slug,
		noteGroup.Name,
		noteGroup.Description,
		noteGroup.ImageURL,
		noteGroup.CreatedAt,
		noteGroup.UpdatedAt,
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
		return fmt.Errorf("database error: %w: %w", ErrNoteGroupAlreadyExists, pgErr)
	}

	return err
}

func (service NoteGroupService) updateNoteGroup(noteGroup *internal.NoteGroup) error {
	q := `
		UPDATE note_groups 
		SET slug = $2,
		    name = $3,
		    description = $4,
		    image_url = $5,
		    updated_at = $6
		WHERE id = $1
	`

	noteGroup.UpdatedAt = time.Now()
	_, err := service.db.Exec(context.Background(),
		q,
		noteGroup.ID,
		noteGroup.Slug,
		noteGroup.Name,
		noteGroup.Description,
		noteGroup.ImageURL,
		noteGroup.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("database error: update note group error: %w", err)
	}

	return nil
}

func (service NoteGroupService) Find(publicId string) (*internal.NoteGroup, error) {
	var noteGroup internal.NoteGroup

	q := `SELECT * FROM note_groups WHERE public_id = $1`

	if err := service.db.QueryRow(context.Background(), q, publicId).
		Scan(
			&noteGroup.ID,
			&noteGroup.PublicId,
			&noteGroup.Name,
			&noteGroup.Slug,
			&noteGroup.Description,
			&noteGroup.ImageURL,
			&noteGroup.CreatedAt,
			&noteGroup.UpdatedAt,
		); err != nil {
		return nil, err
	}

	return &noteGroup, nil
}

func (service NoteGroupService) FindBySlug(s string) (*internal.NoteGroup, error) {
	var noteGroup internal.NoteGroup

	q := `SELECT * FROM note_groups WHERE slug = $1`

	if err := service.db.QueryRow(context.Background(), q, s).
		Scan(
			&noteGroup.ID,
			&noteGroup.PublicId,
			&noteGroup.Name,
			&noteGroup.Slug,
			&noteGroup.Description,
			&noteGroup.ImageURL,
			&noteGroup.CreatedAt,
			&noteGroup.UpdatedAt,
		); err != nil {
		return nil, err
	}

	return &noteGroup, nil
}
