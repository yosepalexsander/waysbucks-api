package repository

import (
	"context"

	"github.com/yosepalexsander/waysbucks-api/entity"
)

type ProductRepository interface {
	ProductFinder
	ProductMutator
	ProductRemover
}

type ProductFinder interface {
	FindProducts(ctx context.Context, whereClauses []string, orderClause string) ([]entity.Product, error)
	FindProduct(ctx context.Context, id int) (*entity.Product, error)
	FindToppings(ctx context.Context) ([]entity.ProductTopping, error)
	FindTopping(ctx context.Context, id int) (*entity.ProductTopping, error)
}

type ProductMutator interface {
	SaveProduct(ctx context.Context, product entity.Product) error
	UpdateProduct(ctx context.Context, id int, newProduct map[string]interface{}) error
	SaveTopping(ctx context.Context, topping entity.ProductTopping) error
	UpdateTopping(ctx context.Context, id int, newData map[string]interface{}) error
}

type ProductRemover interface {
	DeleteProduct(ctx context.Context, id int) error
	DeleteTopping(ctx context.Context, id int) error
}
