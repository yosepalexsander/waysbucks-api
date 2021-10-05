package usecase

import (
	"context"

	"github.com/yosepalexsander/waysbucks-api/entity"
	"github.com/yosepalexsander/waysbucks-api/repository"
)

// type AddressUseCase interface {
// 	GetUserAddress(ctx context.Context, userID int) (*[]entity.UserAddress, error)
// 	GetAddress(ctx context.Context, addressID int) (*entity.UserAddress, error)
// 	CreateNewAddress(ctx context.Context, address entity.UserAddress) error
// 	UpdateAddress(ctx context.Context, addressID int, newAddress map[string]interface{}) error
// 	DeleteAddress(ctx context.Context, addressID int, userID int) error
// }

type AddressUseCase struct {
	AddressRepository repository.AddressRepository
}

func (u *AddressUseCase) GetUserAddress(ctx context.Context, userID int) (*[]entity.UserAddress, error) {
	return u.AddressRepository.FindUserAddress(ctx, userID)
}

func (u *AddressUseCase) GetAddress(ctx context.Context, addressID int) (*entity.UserAddress, error)  {
	return u.AddressRepository.FindAddress(ctx, addressID)
}

func (u *AddressUseCase) CreateNewAddress(ctx context.Context, address entity.UserAddress) error {
	
	if err := u.AddressRepository.SaveAddress(ctx, address); err != nil {
		return err
	}
	return nil
}

func (u *AddressUseCase) UpdateAddress(ctx context.Context, addressID int, newAddress map[string]interface{}) error {
	return u.AddressRepository.UpdateAddress(ctx, addressID, newAddress)
}

func (u *AddressUseCase) DeleteAddress(ctx context.Context, addressID int, userID int) error  {
	return u.AddressRepository.DeleteAddress(ctx, addressID, userID)
}