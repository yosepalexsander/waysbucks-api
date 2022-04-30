package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/yosepalexsander/waysbucks-api/entity"
	"github.com/yosepalexsander/waysbucks-api/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserUseCase struct {
	repo repository.UserRepository
}

func NewUserUseCase(repo repository.UserRepository) UserUseCase {
	return UserUseCase{repo}
}

func (u *UserUseCase) FindUsers(ctx context.Context) ([]entity.User, error) {
	return u.repo.FindUsers(ctx)
}

func (u *UserUseCase) GetProfile(ctx context.Context, id string) (*entity.User, error) {
	return u.repo.FindUserById(ctx, id)
}

func (u *UserUseCase) FindUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	return u.repo.FindUserByEmail(ctx, email)
}

func (u *UserUseCase) CreateNewUser(ctx context.Context, name string, email string, password string, gender string, phone string) (*entity.User, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	user := entity.NewUser(name, email, password, gender, phone)

	hashedPassword, err := hashPassword(user.Password)
	if err != nil {
		return nil, err
	}

	user.Id = id.String()
	user.Password = hashedPassword

	if err := u.repo.SaveUser(ctx, user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (u *UserUseCase) ValidatePassword(hashedPassword string, password string) error {
	hashedPasswordBytes, passwordBytes := []byte(hashedPassword), []byte(password)

	if err := bcrypt.CompareHashAndPassword(hashedPasswordBytes, passwordBytes); err != nil {
		return err
	}

	return nil
}

func (u *UserUseCase) ChangePassword(ctx context.Context, id string, newPass string) error {
	hashedPassword, err := hashPassword(newPass)
	if err != nil {
		return err
	}

	newData := make(map[string]interface{}, 1)
	newData["password"] = hashedPassword

	if err := u.repo.UpdateUser(ctx, id, newData); err != nil {
		return err
	}

	return nil
}

func (u *UserUseCase) UpdateUser(ctx context.Context, id string, newData map[string]interface{}) error {
	return u.repo.UpdateUser(ctx, id, newData)
}

func (u *UserUseCase) DeleteUser(ctx context.Context, id string) error {
	return u.repo.DeleteUser(ctx, id)
}

func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}
