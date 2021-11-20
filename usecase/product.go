package usecase

import (
	"context"
	"database/sql"
	"mime/multipart"
	"sync"

	"github.com/yosepalexsander/waysbucks-api/entity"
	"github.com/yosepalexsander/waysbucks-api/repository"
	"github.com/yosepalexsander/waysbucks-api/thirdparty"
)

type ProductUseCase struct {
	repository.ProductRepository
	repository.ToppingRepository
}

func (u *ProductUseCase) GetProducts(ctx context.Context) ([]entity.Product, error)  {
	products, err := u.ProductRepository.FindProducts(ctx)

	if err != nil {
		return nil, err
	}
	var wg sync.WaitGroup
	
	for i := range products {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			imageUrl, _ := thirdparty.GetImageUrl(ctx, products[i].Image)
			products[i].Image = imageUrl
		}(i)
	}

	wg.Wait()
	
	return products, nil
}

func (u *ProductUseCase) GetProduct(ctx context.Context, productID int) (*entity.Product, error) {
	product, err := u.ProductRepository.FindProduct(ctx, productID)

	if err != nil {
		return nil, err
	}

	imageUrl, _ := thirdparty.GetImageUrl(ctx, product.Image)
	product.Image = imageUrl
	
	return product, nil
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

func (u *ProductUseCase) GetToppings(ctx context.Context) ([]entity.ProductTopping, error) {
	toppings, err := u.ToppingRepository.FindToppings(ctx)
	
	switch {
		case err != nil:
			return nil, err
		case len(toppings) == 0:
			return nil, sql.ErrNoRows
	}
	
	var wg sync.WaitGroup

	for i := range toppings {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			imageUrl, _ := thirdparty.GetImageUrl(ctx, toppings[i].Image)
			toppings[i].Image = imageUrl
		}(i)
	}
	wg.Wait()
	
	return toppings, nil
}

func (u *ProductUseCase) GetTopping(ctx context.Context, id int) (*entity.ProductTopping, error) {
	return u.FindTopping(ctx, id)
}

func (u *ProductUseCase) CreateTopping(ctx context.Context, topping entity.ProductTopping) error {
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
	wg := &sync.WaitGroup{}
	var uploadErr error
	wg.Add(2)
	go func() {
		defer wg.Done()
		if err := thirdparty.UploadFile(ctx, file, newName); err != nil  {
			uploadErr = err
			return
		}
	}()
	go func() {
		defer wg.Done()
		if err := thirdparty.RemoveFile(ctx, oldName); err != nil {
			uploadErr = err
			return
		}
	}()
	wg.Wait()

	if uploadErr != nil {
		return uploadErr
	}

	return nil
}

func (u *ProductUseCase) DeleteTopping(ctx context.Context, id int) error {
	return u.ToppingRepository.DeleteTopping(ctx, id)
}