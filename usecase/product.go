package usecase

import (
	"context"
	"database/sql"

	"github.com/yosepalexsander/waysbucks-api/entity"
	"github.com/yosepalexsander/waysbucks-api/helper"
	"github.com/yosepalexsander/waysbucks-api/repository"
	"github.com/yosepalexsander/waysbucks-api/thirdparty"
	"golang.org/x/sync/errgroup"
)

type ProductUseCase struct {
	repo repository.ProductRepository
}

func NewProductUseCase(repo repository.ProductRepository) ProductUseCase {
	return ProductUseCase{repo}
}

func (u *ProductUseCase) GetProducts(ctx context.Context, params map[string][]string) ([]entity.Product, error) {
	whereClauses, orderClauses := helper.QueryParamsToSqlClauses(params)
	products, err := u.repo.FindProducts(ctx, whereClauses, orderClauses)

	switch {
	case err != nil:
		return nil, err
	case len(products) == 0:
		return nil, sql.ErrNoRows
	}

	g, ctx := errgroup.WithContext(ctx)
	for i := range products {
		i := i
		g.Go(func() error {
			imageUrl, err := thirdparty.GetImageUrl(ctx, products[i].Image)

			if err != nil {
				return err
			}

			products[i].Image = imageUrl
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return nil, thirdparty.ErrServiceUnavailable
	}

	return products, nil
}

func (u *ProductUseCase) GetProduct(ctx context.Context, productID int) (*entity.Product, error) {
	product, err := u.repo.FindProduct(ctx, productID)

	if err != nil {
		return nil, err
	}
	if product == nil {
		return nil, sql.ErrNoRows
	}

	imageUrl, _ := thirdparty.GetImageUrl(ctx, product.Image)
	product.Image = imageUrl

	return product, nil
}

func (u *ProductUseCase) CreateProduct(ctx context.Context, productReq entity.ProductRequest) error {
	product := entity.NewProduct(productReq)

	return u.repo.SaveProduct(ctx, product)
}

func (u *ProductUseCase) UpdateProduct(ctx context.Context, id int, newData map[string]interface{}) error {
	return u.repo.UpdateProduct(ctx, id, newData)
}

func (u *ProductUseCase) DeleteProduct(ctx context.Context, id int) error {
	return u.repo.DeleteProduct(ctx, id)
}

func (u *ProductUseCase) GetToppings(ctx context.Context) ([]entity.ProductTopping, error) {
	toppings, err := u.repo.FindToppings(ctx)

	switch {
	case err != nil:
		return nil, err
	case len(toppings) == 0:
		return nil, sql.ErrNoRows
	}

	g, ctx := errgroup.WithContext(ctx)

	for i := range toppings {
		i := i
		g.Go(func() error {
			imageUrl, err := thirdparty.GetImageUrl(ctx, toppings[i].Image)
			if err != nil {
				return err
			}

			toppings[i].Image = imageUrl

			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return nil, thirdparty.ErrServiceUnavailable
	}

	return toppings, nil
}

func (u *ProductUseCase) GetTopping(ctx context.Context, id int) (*entity.ProductTopping, error) {
	return u.repo.FindTopping(ctx, id)
}

func (u *ProductUseCase) CreateTopping(ctx context.Context, toppingReq entity.ProductToppingRequest) error {
	topping := entity.NewProductTopping(toppingReq)

	if err := u.repo.SaveTopping(ctx, topping); err != nil {
		_ = thirdparty.RemoveFile(ctx, topping.Name)
		return err
	}

	return nil
}

func (u *ProductUseCase) UpdateTopping(ctx context.Context, id int, newData map[string]interface{}) error {
	return u.repo.UpdateTopping(ctx, id, newData)
}

func (u *ProductUseCase) DeleteTopping(ctx context.Context, id int) error {
	return u.repo.DeleteTopping(ctx, id)
}
