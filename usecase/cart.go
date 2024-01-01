package usecase

import (
	"context"

	"github.com/yosepalexsander/waysbucks-api/entity"
	"github.com/yosepalexsander/waysbucks-api/repository"
)

type CartUseCase struct {
	repo repository.CartRepository
}

func NewCartUseCase(r repository.CartRepository) CartUseCase {
	return CartUseCase{r}
}

func (u *CartUseCase) FindCarts(ctx context.Context, userID string) ([]entity.Cart, error) {
	carts, err := u.repo.FindCarts(ctx, userID)
	if err != nil {
		return nil, err
	}

	return carts, nil
}

func (u *CartUseCase) SaveCart(ctx context.Context, req entity.CartRequest, userId string) error {
	cart := entity.NewCart(req.ProductId, req.Price, req.Qty, req.ToppingIds, userId)

	err := u.repo.SaveCart(ctx, cart)
	if err != nil {
		return err
	}

	return nil
}

func (u *CartUseCase) UpdateCart(ctx context.Context, id int, userID string, data map[string]interface{}) error {
	return u.repo.UpdateCart(ctx, id, userID, data)
}

func (u *CartUseCase) DeleteCart(ctx context.Context, id int, userID string) error {
	return u.repo.DeleteCart(ctx, id, userID)
}
