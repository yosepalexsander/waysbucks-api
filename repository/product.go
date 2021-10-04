package repository

import (
	"context"

	"github.com/yosepalexsander/waysbucks-api/entity"
)

type ProductRepository interface {
	FindProducts(ctx context.Context) ([]entity.Product, error)
	FindProduct(ctx context.Context, id int) (*entity.Product, error)
	SaveProduct(ctx context.Context, product entity.Product) error
	UpdateProduct(ctx context.Context, id int, newProduct map[string]interface{}) error
	DeleteProduct(ctx context.Context, id int) error
}