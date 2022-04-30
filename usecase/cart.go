package usecase

import (
	"context"
	"database/sql"

	"github.com/yosepalexsander/waysbucks-api/entity"
	"github.com/yosepalexsander/waysbucks-api/repository"
	"github.com/yosepalexsander/waysbucks-api/thirdparty"
	"golang.org/x/sync/errgroup"
)

type CartUseCase struct {
	repo repository.CartRepository
}

func NewCartUseCase(r repository.CartRepository) CartUseCase {
	return CartUseCase{r}
}

func (u *CartUseCase) GetCarts(ctx context.Context, userID string) ([]entity.Cart, error) {
	carts, err := u.repo.FindCarts(ctx, userID)
	switch {
	case err != nil:
		return nil, err
	case len(carts) == 0:
		return nil, sql.ErrNoRows
	}

	g, ctx := errgroup.WithContext(ctx)

	for i := range carts {
		i := i
		g.Go(func() error {
			imageUrl, err := thirdparty.GetImageUrl(ctx, carts[i].Product.Image)
			if err == nil && imageUrl != "" {
				carts[i].Product.Image = imageUrl
			}
			return err
		})
	}

	if err := g.Wait(); err != nil {
		return nil, thirdparty.ErrServiceUnavailable
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
