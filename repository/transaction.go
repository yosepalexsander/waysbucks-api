package repository

import (
	"context"

	"github.com/yosepalexsander/waysbucks-api/entity"
)

type TransactionFinder interface {
	FindTransactions(ctx context.Context) ([]entity.Transaction, error)
	FindUserTransactions(ctx context.Context, userID int) ([]entity.Transaction, error)
	FindTransactionByID(ctx context.Context, id int) (*entity.Transaction, error)
}

type TransactionMutator interface {
	UpdateTransaction(ctx context.Context, id int, data map[string]interface{}) error
}

type TransactionTx interface {
	ExecTx(ctx context.Context, fn func(Transactioner) error) error
	OrderTx(ctx context.Context, arg entity.TransactionTxParams) error
}

type Transactioner interface {
	DeleteCart(ctx context.Context, productID int, userID int) error
	CreateOrder(ctx context.Context, order entity.Order) error
	CreateTransaction(ctx context.Context, tx entity.Transaction) (int, error)
	Rollback() error
	Commit() error
}
