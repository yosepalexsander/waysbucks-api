package storage

import (
	"context"

	"github.com/yosepalexsander/waysbucks-api/domain"
)
type UserFinder interface {
	FindUserById(ctx context.Context, id uint32)
}
type UserSaver interface {
	SaveUser(ctx context.Context, user domain.User)
	UpdateUser(ctx context.Context, user domain.User)
}

func FindUserById(ctx context.Context, id uint32) (*domain.User, error) {
	return nil, nil
}

func SaveUser(ctx context.Context, user domain.User) error {
	return nil
}

func UpdateUser(ctx context.Context, user domain.User) (*domain.User, error) {

	return nil, nil
}