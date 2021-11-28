package usecase

import (
	"context"

	"github.com/yosepalexsander/waysbucks-api/entity"
	"github.com/yosepalexsander/waysbucks-api/repository"
)

type AddressUseCase struct {
	Finder  repository.AddressFinder
	Mutator repository.AddressMutator
}

func NewAddressUseCase(rf repository.AddressFinder, rm repository.AddressMutator) AddressUseCase {
	return AddressUseCase{rf, rm}
}

func (u *AddressUseCase) GetUserAddress(ctx context.Context, userID int) ([]entity.Address, error) {
	return u.Finder.FindUserAddress(ctx, userID)
}

func (u *AddressUseCase) GetAddress(ctx context.Context, addressID int) (*entity.Address, error) {
	return u.Finder.FindAddress(ctx, addressID)
}

func (u *AddressUseCase) CreateNewAddress(ctx context.Context, userID int, newAddress entity.AddressRequest) error {
	address := addressFromRequest(newAddress)
	if err := u.Mutator.SaveAddress(ctx, userID, address); err != nil {
		return err
	}
	return nil
}

func (u *AddressUseCase) UpdateAddress(ctx context.Context, addressID int, newAddress map[string]interface{}) error {
	return u.Mutator.UpdateAddress(ctx, addressID, newAddress)
}

func (u *AddressUseCase) DeleteAddress(ctx context.Context, addressID int, userID int) error {
	return u.Mutator.DeleteAddress(ctx, addressID, userID)
}

func addressFromRequest(req entity.AddressRequest) entity.Address {
	return entity.Address{
		Name:       req.Name,
		Phone:      req.Phone,
		Address:    req.Address,
		City:       req.City,
		PostalCode: req.PostalCode,
	}
}
