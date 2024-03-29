package postgresql

import "github.com/jackc/pgx/v5"

type Services struct {
	House     *HouseService
	Note      *NoteService
	NoteGroup *NoteGroupService
}

func NewServices(db *pgx.Conn) *Services {
	return &Services{
		House:     &HouseService{db: db},
		Note:      &NoteService{db: db},
		NoteGroup: &NoteGroupService{db: db},
	}
}
