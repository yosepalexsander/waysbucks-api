package usecase

import (
	"context"
	"errors"
	"log"
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
func (u *TransactionUseCase) GetTransactions(ctx context.Context, userID int) ([]entity.TransactionResponse, error) {
	transactions, err := u.Finder.FindTransactions(ctx)
	if err != nil {
		return nil, err
	}

	var transactionResponse []entity.TransactionResponse

	for _, v := range transactions {
		transaction, err := transactionResponseFromDB(&v)
		if err != nil {
			return nil, err
		}
		transactionResponse = append(transactionResponse, transaction)
	}

	return transactionResponse, nil
}

func (u *TransactionUseCase) GetUserTransactions(ctx context.Context, userID int) ([]entity.TransactionResponse, error) {
	transactions, err := u.Finder.FindUserTransactions(ctx, userID)
	if err != nil {
		return nil, err
	}

	var transactionResponse []entity.TransactionResponse

	for _, v := range transactions {
		transaction, err := transactionResponseFromDB(&v)
		if err != nil {
			return nil, err
		}
		transactionResponse = append(transactionResponse, transaction)
	}

	return transactionResponse, nil
}

func (u *TransactionUseCase) GetDetailTransaction(ctx context.Context, id int) (entity.TransactionResponse, error) {
	transaction, err := u.Finder.FindTransactionByID(ctx, id)
	if err != nil {
		return entity.TransactionResponse{}, err
	}

	transactionResponse, err := transactionResponseFromDB(transaction)

	if err == nil {
		return transactionResponse, nil
	}
	return transactionResponse, err
}

func (u *TransactionUseCase) MakeTransaction(ctx context.Context, request entity.TransactionRequest) error {
	transaction := transactionFromRequest(request)
	log.Println(transaction)
	return u.Transactioner.OrderTx(ctx, transaction)
}

func (u *TransactionUseCase) UpdateTransaction(ctx context.Context, id int, data map[string]interface{}) error {
	return u.Mutator.UpdateTransaction(ctx, id, data)
}

func transactionResponseFromDB(t *entity.Transaction) (entity.TransactionResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var OrderResponse []entity.OrderResponse

	g, ctx := errgroup.WithContext(ctx)
	for i := range t.Orders {
		i := i
		g.Go(func() error {
			imageUrl, err := helper.GetImageUrl(ctx, t.Orders[i].Image)
			if err == nil {
				t.Orders[i].Image = imageUrl
				OrderResponse = append(OrderResponse, orderResponseFromDB(t.Orders[i]))
			}
			return err
		})
	}

	if err := g.Wait(); err != nil {
		return entity.TransactionResponse{}, errors.New("object storage service unavailable")
	}
	return entity.TransactionResponse{
		Name:       t.Name,
		Address:    t.Address,
		PostalCode: t.PostalCode,
		City:       t.City,
		Total:      t.Total,
		Status:     t.Status,
		Orders:     OrderResponse,
	}, nil
}

func orderResponseFromDB(o entity.Order) entity.OrderResponse {
	return entity.OrderResponse{
		Id:           o.Id,
		OrderProduct: o.OrderProduct,
		Toppings:     o.Toppings,
		Price:        o.Price,
		Qty:          o.Qty,
	}
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
