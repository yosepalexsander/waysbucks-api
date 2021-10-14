package usecase

import (
	"context"
	"database/sql"
	"sync"

	"github.com/yosepalexsander/waysbucks-api/entity"
	"github.com/yosepalexsander/waysbucks-api/helper"
	"github.com/yosepalexsander/waysbucks-api/repository"
)


type CartUseCase struct {
	CartRepository repository.CartRepository
}

func (u *CartUseCase) GetUserCarts(ctx context.Context, userID int) ([]entity.Cart, error) {
	carts, err := u.CartRepository.FindCarts(ctx, userID)
	switch {
		case err != nil:
			return nil, err
		case len(carts) == 0:
			return nil, sql.ErrNoRows
	}
	var wg sync.WaitGroup

	wg.Add(len(carts))
	for i := range carts {
		go func(i int) {
			defer wg.Done()
			imageUrl, err := helper.GetImageUrl(ctx, carts[i].Product.Image)
			if err == nil && imageUrl != "" {
				carts[i].Product.Image = imageUrl
			}
			carts[i].Product_Id = 0
			carts[i].ToppingIds = nil
		}(i)
	}
	wg.Wait()
	
	return carts, nil
}

func (u CartUseCase) SaveToCart(ctx context.Context, cart entity.Cart) error  {
	err := u.CartRepository.SaveToCart(ctx, cart)
	if err != nil {
		return err
	}

	return nil
}

func (u CartUseCase) UpdateCart(ctx context.Context, id int, userID int, data map[string]interface{}) error  {
	return u.CartRepository.UpdateCart(ctx, id, userID, data)
}

func (u CartUseCase) DeleteCart(ctx context.Context, id, userID int) error {
	return u.CartRepository.DeleteCart(ctx, id, userID)
}