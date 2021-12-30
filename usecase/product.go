package usecase

import (
	"context"
	"database/sql"
	"mime/multipart"

	"github.com/yosepalexsander/waysbucks-api/entity"
	"github.com/yosepalexsander/waysbucks-api/helper"
	"github.com/yosepalexsander/waysbucks-api/repository"
	"github.com/yosepalexsander/waysbucks-api/thirdparty"
	"golang.org/x/sync/errgroup"
)

type ProductUseCase struct {
	repository.ProductFinder
	repository.ProductMutator
	repository.ToppingRepository
}

func NewProductUseCase(rpf repository.ProductFinder, rpm repository.ProductMutator, rt repository.ToppingRepository) ProductUseCase {
	return ProductUseCase{rpf, rpm, rt}
}

func (u *ProductUseCase) GetProducts(ctx context.Context, params map[string][]string) ([]entity.Product, error) {
	whereClauses, orderClauses := helper.QueryParamsToSqlClauses(params)
	products, err := u.ProductFinder.FindProducts(ctx, whereClauses, orderClauses)

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
	product, err := u.ProductFinder.FindProduct(ctx, productID)
	if err != nil {
		return nil, err
	}

	if imageUrl, err := thirdparty.GetImageUrl(ctx, product.Image); err != nil {
		return nil, thirdparty.ErrServiceUnavailable
	} else {
		product.Image = imageUrl
		return product, nil
	}
}

func (u *ProductUseCase) CreateProduct(ctx context.Context, productReq entity.ProductRequest) error {
	product := productFromRequest(productReq)

	if err := u.ProductMutator.SaveProduct(ctx, product); err != nil {
		_ = thirdparty.RemoveFile(ctx, product.Name)
		return err
	}

	return nil
}

func (u *ProductUseCase) UpdateProduct(ctx context.Context, id int, newData map[string]interface{}) error {
	return u.ProductMutator.UpdateProduct(ctx, id, newData)
}

func (u *ProductUseCase) DeleteProduct(ctx context.Context, id int) error {
	return u.ProductMutator.DeleteProduct(ctx, id)
}

func (u *ProductUseCase) GetToppings(ctx context.Context) ([]entity.ProductTopping, error) {
	toppings, err := u.ToppingRepository.FindToppings(ctx)

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
	return u.FindTopping(ctx, id)
}

func (u *ProductUseCase) CreateTopping(ctx context.Context, toppingReq entity.ProductToppingRequest) error {
	topping := toppingFromRequest(toppingReq)

	if err := u.ToppingRepository.SaveTopping(ctx, topping); err != nil {
		_ = thirdparty.RemoveFile(ctx, topping.Name)
		return err
	}

	return nil
}

func (u *ProductUseCase) UpdateTopping(ctx context.Context, id int, newData map[string]interface{}) error {
	return u.ToppingRepository.UpdateTopping(ctx, id, newData)
}

func (u *ProductUseCase) UpdateImage(ctx context.Context, file multipart.File, oldName string, newName string) error {
	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		if err := thirdparty.UploadFile(ctx, file, newName); err != nil {
			return err
		}
		return nil
	})

	g.Go(func() error {
		if err := thirdparty.RemoveFile(ctx, oldName); err != nil {
			return err
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		return err
	}

	return nil
}

func (u *ProductUseCase) DeleteTopping(ctx context.Context, id int) error {
	return u.ToppingRepository.DeleteTopping(ctx, id)
}

func productFromRequest(req entity.ProductRequest) entity.Product {
	return entity.Product{
		Name:         req.Name,
		Description:  req.Description,
		Image:        req.Image,
		Price:        req.Price,
		Is_Available: req.Is_Available,
	}
}

func toppingFromRequest(req entity.ProductToppingRequest) entity.ProductTopping {
	return entity.ProductTopping{
		Name:         req.Name,
		Image:        req.Image,
		Price:        req.Price,
		Is_Available: req.Is_Available,
	}
}
