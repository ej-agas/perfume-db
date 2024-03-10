package postgresql

import "github.com/ej-agas/perfume-db/internal"

type HouseService struct {
}

func (h HouseService) Save(house internal.House) error {
	return nil
}

func (h HouseService) Find(id int) (*internal.House, error) {
	return nil, nil
}

func (h HouseService) FindBySlug(s string) (*internal.House, error) {
	return nil, nil
}
