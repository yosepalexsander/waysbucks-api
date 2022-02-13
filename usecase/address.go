package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/yosepalexsander/waysbucks-api/entity"
	"github.com/yosepalexsander/waysbucks-api/repository"
)

type AddressUseCase struct {
	repo repository.AddressRepository
}

func NewAddressUseCase(repo repository.AddressRepository) AddressUseCase {
	return AddressUseCase{repo}
}

func (u *AddressUseCase) GetUserAddresses(ctx context.Context, userID string) ([]entity.Address, error) {
	return u.repo.FindAllUserAddresses(ctx, userID)
}

func (u *AddressUseCase) GetAddress(ctx context.Context, addressID int) (*entity.Address, error) {
	return u.repo.FindAddress(ctx, addressID)
}

func (u *AddressUseCase) CreateNewAddress(ctx context.Context, userID string, newAddress entity.AddressRequest) error {
	address := addressFromRequest(newAddress)
	address.Id = uuid.NewString()
	if err := u.repo.SaveAddress(ctx, userID, address); err != nil {
		return err
	}
	return nil
}

func (u *AddressUseCase) UpdateAddress(ctx context.Context, addressID int, newAddress map[string]interface{}) error {
	return u.repo.UpdateAddress(ctx, addressID, newAddress)
}

func (u *AddressUseCase) DeleteAddress(ctx context.Context, addressID int, userID string) error {
	return u.repo.DeleteAddress(ctx, addressID, userID)
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
