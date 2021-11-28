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
	CartRepository repository.CartRepository
}

func NewCartUseCase(r repository.CartRepository) CartUseCase {
	return CartUseCase{r}
}

func (u *CartUseCase) GetUserCarts(ctx context.Context, userID int) ([]entity.Cart, error) {
	carts, err := u.CartRepository.FindCarts(ctx, userID)
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
				carts[i].Product_Id = 0
				carts[i].ToppingIds = nil
			}
			return err
		})
	}

	if err := g.Wait(); err != nil {
		return nil, thirdparty.ErrServiceUnavailable
	}
	return carts, nil
}

func (u *CartUseCase) SaveToCart(ctx context.Context, cart entity.Cart) error {
	err := u.CartRepository.SaveToCart(ctx, cart)
	if err != nil {
		return err
	}

	return nil
}

func (u *CartUseCase) UpdateCart(ctx context.Context, id int, userID int, data map[string]interface{}) error {
	return u.CartRepository.UpdateCart(ctx, id, userID, data)
}

func (u *CartUseCase) DeleteCart(ctx context.Context, id, userID int) error {
	return u.CartRepository.DeleteCart(ctx, id, userID)
}
