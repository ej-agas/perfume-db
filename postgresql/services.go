package postgresql

import "github.com/jackc/pgx/v5"

type Services struct {
	House     *HouseService
	Note      *NoteService
	NoteGroup *NoteGroupService
	Perfumer  *PerfumerService
	Perfume   *PerfumeService
}

func NewServices(db *pgx.Conn) *Services {
	return &Services{
		House:     &HouseService{db: db},
		Note:      &NoteService{db: db},
		NoteGroup: &NoteGroupService{db: db},
		Perfumer:  &PerfumerService{db: db},
		Perfume:   &PerfumeService{db: db},
	}
}
