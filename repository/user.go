package repository

import (
	"context"

	"github.com/yosepalexsander/waysbucks-api/entity"
)

type UserRepository interface {
	UserFinder
	UserMutator
}

type UserFinder interface {
	FindUsers(ctx context.Context) ([]entity.User, error)
	FindUserById(ctx context.Context, id int) (*entity.User, error)
	FindUserByEmail(ctx context.Context, email string) (*entity.User, error)
}

type UserMutator interface {
	SaveUser(ctx context.Context, user entity.User) error
	UpdateUser(ctx context.Context, id int, newData map[string]interface{}) error
	DeleteUser(ctx context.Context, id int) error
}
