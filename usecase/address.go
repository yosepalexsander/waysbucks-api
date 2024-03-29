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

func (u *AddressUseCase) FindUserAddresses(ctx context.Context, userID string) ([]entity.Address, error) {
	return u.repo.FindAllUserAddresses(ctx, userID)
}

func (u *AddressUseCase) GetAddress(ctx context.Context, id string) (*entity.Address, error) {
	return u.repo.FindAddress(ctx, id)
}

func (u *AddressUseCase) CreateNewAddress(ctx context.Context, userID string, newAddress entity.AddressRequest) error {
	address := addressFromRequest(newAddress)

	id, err := uuid.NewRandom()
	if err != nil {
		return err
	}

	address.Id = id.String()

	if err := u.repo.SaveAddress(ctx, userID, address); err != nil {
		return err
	}

	return nil
}

func (u *AddressUseCase) UpdateAddress(ctx context.Context, id string, newAddress map[string]interface{}) error {
	return u.repo.UpdateAddress(ctx, id, newAddress)
}

func (u *AddressUseCase) DeleteAddress(ctx context.Context, id string, userID string) error {
	return u.repo.DeleteAddress(ctx, id, userID)
}

func addressFromRequest(req entity.AddressRequest) entity.Address {
	return entity.Address{
		Name:       req.Name,
		Phone:      req.Phone,
		Address:    req.Address,
		City:       req.City,
		PostalCode: req.PostalCode,
		Longitude:  req.Longitude,
		Latitude:   req.Latitude,
	}
}
