package repository

import (
	"context"

	"github.com/yosepalexsander/waysbucks-api/entity"
)

type CartRepository interface {
	FindCarts(ctx context.Context, userID int) ([]entity.Cart, error)
	SaveCart(ctx context.Context, cart entity.Cart) error
	UpdateCart(ctx context.Context, id int, userID int, data map[string]interface{}) error
	DeleteCart(ctx context.Context, id int, userID int) error
}
