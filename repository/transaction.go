package repository

import (
	"context"

	"github.com/yosepalexsander/waysbucks-api/entity"
)


type TransactionFinder interface {
	FindTransactions(ctx context.Context, userID int) ([]entity.Transaction, error)
}

type TransactionTx interface {
	ExecTx(ctx context.Context, fn func(Transactioner) error) error
}

type Transactioner interface {
	DeleteCart(ctx context.Context, productIds []int, userID int) error
	CreateOrder(ctx context.Context, order []entity.Order) error
	CreateTransaction(ctx context.Context, tx entity.Transaction) (int, error)
	Rollback() error 
	Commit() error
}