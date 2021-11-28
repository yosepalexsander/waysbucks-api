package repository

import (
	"context"

	"github.com/yosepalexsander/waysbucks-api/entity"
)

type AddressFinder interface {
	FindUserAddress(ctx context.Context, userID int) ([]entity.Address, error)
	FindAddress(ctx context.Context, id int) (*entity.Address, error)
}

type AddressMutator interface {
	SaveAddress(ctx context.Context, userID int, address entity.Address) error
	UpdateAddress(ctx context.Context, id int, newAddress map[string]interface{}) error
	DeleteAddress(ctx context.Context, id int, userID int) error
}
