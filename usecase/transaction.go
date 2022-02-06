package usecase

import (
	"context"
	"time"

	"github.com/yosepalexsander/waysbucks-api/entity"
	"github.com/yosepalexsander/waysbucks-api/helper"
	"github.com/yosepalexsander/waysbucks-api/repository"
	"github.com/yosepalexsander/waysbucks-api/thirdparty"
	"golang.org/x/sync/errgroup"
)

type TransactionUseCase struct {
	repo repository.TransactionRepository
}

func NewTransactionUseCase(repo repository.TransactionRepository) TransactionUseCase {
	return TransactionUseCase{repo}
}

func (u *TransactionUseCase) GetTransactions(ctx context.Context) ([]entity.Transaction, error) {
	transactions, err := u.repo.FindTransactions(ctx)
	if err != nil {
		return nil, err
	}

	g := new(errgroup.Group)
	for _, t := range transactions {
		t := t
		g.Go(func() error {
			for i := range t.Orders {
				imageUrl, err := thirdparty.GetImageUrl(ctx, t.Orders[i].Image)
				if err != nil {
					return err
				}
				t.Orders[i].Image = imageUrl
			}
			return nil
		})

	}

	if err := g.Wait(); err != nil {
		return nil, thirdparty.ErrServiceUnavailable
	}

	return transactions, nil
}

func (u *TransactionUseCase) GetUserTransactions(ctx context.Context, userID int) ([]entity.Transaction, error) {
	transactions, err := u.repo.FindUserTransactions(ctx, userID)
	if err != nil {
		return nil, err
	}

	g := new(errgroup.Group)
	for _, t := range transactions {
		t := t
		g.Go(func() error {
			for i := range t.Orders {
				imageUrl, err := thirdparty.GetImageUrl(ctx, t.Orders[i].Image)
				if err != nil {
					return err
				}
				t.Orders[i].Image = imageUrl
			}
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return nil, thirdparty.ErrServiceUnavailable
	}

	return transactions, nil
}

func (u *TransactionUseCase) GetDetailTransaction(ctx context.Context, id string) (*entity.Transaction, error) {
	transaction, err := u.repo.FindTransactionByID(ctx, id)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	for i := range transaction.Orders {
		imageUrl, err := thirdparty.GetImageUrl(ctx, transaction.Orders[i].Image)
		if err != nil {
			cancel()
			return nil, thirdparty.ErrServiceUnavailable
		}
		transaction.Orders[i].Image = imageUrl
	}
	return transaction, err
}

func (u *TransactionUseCase) MakeTransaction(ctx context.Context, request entity.TransactionRequest) (*entity.Transaction, error) {
	transaction := transactionFromRequest(request)
	if err := u.orderTx(ctx, transaction); err != nil {
		return nil, err
	}
	newTransaction, err := u.GetDetailTransaction(ctx, transaction.Transaction.Id)
	newTransaction.Email = transaction.Transaction.Email
	newTransaction.ServiceFee = transaction.Transaction.ServiceFee
	if err != nil {
		return nil, err
	}

	return newTransaction, nil
}

func (u *TransactionUseCase) orderTx(ctx context.Context, arg entity.TransactionTxParams) error {
	txErr := u.repo.ExecTx(ctx, func(tx repository.Transactioner) error {
		var err error

		id, err := tx.CreateTransaction(ctx, arg.Transaction)
		if err != nil {
			return err
		}
		for i := range arg.Order {
			arg.Order[i].Transaction_Id = id
			err := tx.CreateOrder(ctx, arg.Order[i])
			if err != nil {
				return err
			}

			err = tx.DeleteCart(ctx, arg.Order[i].Product_Id, arg.Transaction.User_Id)
			if err != nil {
				return err
			}
		}

		return nil
	})

	if txErr != nil {
		return txErr
	}

	return nil
}

func (u *TransactionUseCase) UpdateTransaction(ctx context.Context, id string, data map[string]interface{}) error {
	return u.repo.UpdateTransaction(ctx, id, data)
}

func transactionFromRequest(r entity.TransactionRequest) entity.TransactionTxParams {
	var orders []entity.Order

	for _, v := range r.Order {
		orders = append(orders, orderFromRequest(v))
	}

	return entity.TransactionTxParams{
		Transaction: entity.Transaction{
			Id:         "ORDER-" + helper.RandString(20),
			User_Id:    r.User_Id,
			Name:       r.Name,
			Email:      r.Email,
			Address:    r.Address,
			City:       r.City,
			PostalCode: r.PostalCode,
			Phone:      r.Phone,
			Total:      r.Total,
			ServiceFee: r.ServiceFee,
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
