package repository

import (
	"context"

	"github.com/yosepalexsander/waysbucks-api/entity"
)

type AddressRepository interface {
	FindUserAddress(ctx context.Context, userID int) (*[]entity.UserAddress, error)
	FindAddress(ctx context.Context, id int) (*entity.UserAddress, error)
	SaveAddress(ctx context.Context, address entity.UserAddress) error
	UpdateAddress(ctx context.Context, id int, newAddress map[string]interface{}) error
	DeleteAddress(ctx context.Context, id int, userID int) error
}
