package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/yosepalexsander/waysbucks-api/entity"
	"github.com/yosepalexsander/waysbucks-api/helper"
	"github.com/yosepalexsander/waysbucks-api/repository"
	"golang.org/x/sync/errgroup"
)

type TransactionUseCase struct {
	Finder        repository.TransactionFinder
	Transactioner repository.TransactionTx
	Mutator       repository.TransactionMutator
}

func NewTransactionUseCase(f repository.TransactionFinder, t repository.TransactionTx, m repository.TransactionMutator) TransactionUseCase {
	return TransactionUseCase{f, t, m}
}
func (u *TransactionUseCase) GetTransactions(ctx context.Context) ([]entity.Transaction, error) {
	transactions, err := u.Finder.FindTransactions(ctx)
	if err != nil {
		return nil, err
	}

	g := new(errgroup.Group)
	for _, t := range transactions {
		t := t
		g.Go(func() error {
			for i := range t.Orders {
				imageUrl, err := helper.GetImageUrl(ctx, t.Orders[i].Image)
				if err != nil {
					return err
				}
				t.Orders[i].Image = imageUrl
			}
			return nil
		})

	}

	if err := g.Wait(); err != nil {
		return nil, errors.New("object storage service unavailable")
	}

	return transactions, nil
}

func (u *TransactionUseCase) GetUserTransactions(ctx context.Context, userID int) ([]entity.Transaction, error) {
	transactions, err := u.Finder.FindUserTransactions(ctx, userID)
	if err != nil {
		return nil, err
	}

	g := new(errgroup.Group)
	for _, t := range transactions {
		t := t
		g.Go(func() error {
			for i := range t.Orders {
				imageUrl, err := helper.GetImageUrl(ctx, t.Orders[i].Image)
				if err != nil {
					return err
				}
				t.Orders[i].Image = imageUrl
			}
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return nil, errors.New("object storage service unavailable")
	}

	return transactions, nil
}

func (u *TransactionUseCase) GetDetailTransaction(ctx context.Context, id int) (*entity.Transaction, error) {
	transaction, err := u.Finder.FindTransactionByID(ctx, id)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	for _, order := range transaction.Orders {
		imageUrl, err := helper.GetImageUrl(ctx, order.Image)
		if err != nil {
			cancel()
			return nil, errors.New("object storage service unavailable")
		}
		order.Image = imageUrl
	}
	return transaction, err
}

func (u *TransactionUseCase) MakeTransaction(ctx context.Context, request entity.TransactionRequest) error {
	transaction := transactionFromRequest(request)
	return u.Transactioner.OrderTx(ctx, transaction)
}

func (u *TransactionUseCase) UpdateTransaction(ctx context.Context, id int, data map[string]interface{}) error {
	return u.Mutator.UpdateTransaction(ctx, id, data)
}

func transactionFromRequest(r entity.TransactionRequest) entity.TransactionTxParams {
	var orders []entity.Order

	for _, v := range r.Order {
		orders = append(orders, orderFromRequest(v))
	}

	return entity.TransactionTxParams{
		Transaction: entity.Transaction{
			User_Id:    r.User_Id,
			Name:       r.Name,
			Address:    r.Address,
			PostalCode: r.PostalCode,
			City:       r.City,
			Phone:      r.Phone,
			Total:      r.Total,
			Status:     r.Status,
		},
		Order: orders,
	}
}

func orderFromRequest(r entity.OrderRequest) entity.Order {
	return entity.Order{
		Product_Id:  r.Product_Id,
		Qty:         r.Qty,
		Price:       r.Price,
		Topping_Ids: r.Topping_Ids,
	}
}
