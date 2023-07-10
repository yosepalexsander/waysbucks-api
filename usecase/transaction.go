package usecase

import (
	"context"

	"github.com/yosepalexsander/waysbucks-api/entity"
	"github.com/yosepalexsander/waysbucks-api/repository"
)

type TransactionUseCase struct {
	repo repository.TransactionRepository
}

func NewTransactionUseCase(repo repository.TransactionRepository) TransactionUseCase {
	return TransactionUseCase{repo}
}

func (u *TransactionUseCase) FindTransactions(ctx context.Context) ([]entity.Transaction, error) {
	transactions, err := u.repo.FindTransactions(ctx)
	if err != nil {
		return nil, err
	}

	return transactions, nil
}

func (u *TransactionUseCase) GetUserTransactions(ctx context.Context, userID string) ([]entity.Transaction, error) {
	transactions, err := u.repo.FindUserTransactions(ctx, userID)
	if err != nil {
		return nil, err
	}

	return transactions, nil
}

func (u *TransactionUseCase) GetDetailTransaction(ctx context.Context, id string) (*entity.Transaction, error) {
	transaction, err := u.repo.FindTransactionByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return transaction, err
}

func (u *TransactionUseCase) MakeTransaction(ctx context.Context, request entity.TransactionRequest) (*entity.Transaction, error) {
	transaction := entity.NewTransaction(request)
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
			arg.Order[i].TransactionId = id
			err := tx.CreateOrder(ctx, arg.Order[i])
			if err != nil {
				return err
			}

			err = tx.DeleteCart(ctx, arg.Order[i].ProductId, arg.Transaction.UserId)
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
