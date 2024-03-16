package postgresql

import "github.com/jackc/pgx/v5"

type Services struct {
	House *HouseService
}

func NewServices(db *pgx.Conn) *Services {
	return &Services{
		House: &HouseService{db: db},
	}
}
