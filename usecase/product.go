package usecase

import (
	"context"
	"sync"

	"github.com/yosepalexsander/waysbucks-api/entity"
	"github.com/yosepalexsander/waysbucks-api/helper"
	"github.com/yosepalexsander/waysbucks-api/repository"
)

type ProductUseCase struct {
	ProductRepository repository.ProductRepository
}

func (u *ProductUseCase) GetAllProduct(ctx context.Context) ([]entity.Product, error)  {
	products, err := u.ProductRepository.FindProducts(ctx)

	if err != nil {
		return nil, err
	}

	var wg sync.WaitGroup
	// done := make(chan bool, 1)

	for i := range products {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			imageUrl, _ := helper.GetImageUrl(ctx, products[i].Image)
			products[i].Image = imageUrl
		}(i)
	}
	wg.Wait()
	return products, nil
}

func (u *ProductUseCase) GetProduct(ctx context.Context, productID int) (*entity.Product, error) {
	return u.ProductRepository.FindProduct(ctx, productID)
}

func (u *ProductUseCase) CreateProduct(ctx context.Context, product entity.Product) error {
	return u.ProductRepository.SaveProduct(ctx, product)
}

func (u *ProductUseCase) UpdateProduct(ctx context.Context, id int, newData map[string]interface{}) error {
	return u.ProductRepository.UpdateProduct(ctx, id, newData)
}

func (u *ProductUseCase) DeleteProduct(ctx context.Context, id int) error {
	return u.ProductRepository.DeleteProduct(ctx, id)
}