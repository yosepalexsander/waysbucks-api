package usecase

import (
	"context"

	"github.com/yosepalexsander/waysbucks-api/entity"
	"github.com/yosepalexsander/waysbucks-api/repository"
)


type TransactionUseCase struct {
	Finder repository.TransactionFinder
	Transactioner repository.TransactionTx
}

func (u *TransactionUseCase) GetUserTransactions(ctx context.Context, userID int)   {
	
}

func (u *TransactionUseCase) MakeTransaction(ctx context.Context, arg entity.TransactionTxParams) error {

	
	return nil
}