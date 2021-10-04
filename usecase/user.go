package usecase

import (
	"context"

	"github.com/yosepalexsander/waysbucks-api/entity"
	"github.com/yosepalexsander/waysbucks-api/repository"
	"golang.org/x/crypto/bcrypt"
)

// type UserUseCase interface{
// 	FindUserById(ctx context.Context, id int) (*entity.User, error)
// 	FindUserByEmail(ctx context.Context, email string) (*entity.User, error)
// 	CreateNewUser(ctx context.Context, user entity.User) error
// 	UpdateUser(ctx context.Context,id int, newData map[string]interface{}) (error)
// 	DeleteUser(ctx context.Context, id int) error
// 	ValidatePassword(hashedPassword string, password string) error
// 	ChangePassword(ctx context.Context, id int, newPass string) error
// }

type UserUseCase struct{
	UserRepository repository.UserRepository
}

func (u *UserUseCase) FindUserById(ctx context.Context, id int) (*entity.User, error) {
	return u.UserRepository.FindUserById(ctx, id)
}

func (u *UserUseCase) FindUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	return u.UserRepository.FindUserByEmail(ctx, email)
}

func (u *UserUseCase) CreateNewUser(ctx context.Context, user entity.User) error {
	hashedPassword, err := hashPassword(user.Password)

	if err != nil {
		return err
	}

	user.Password = hashedPassword
	
	if err := u.UserRepository.SaveUser(ctx, user); err != nil {
		return err
	}

	return nil
}

func (u *UserUseCase) ValidatePassword(hashedPassword string, password string) error {
	hashedPasswordBytes, passwordBytes := []byte(hashedPassword), []byte(password)
	
	if err := bcrypt.CompareHashAndPassword(hashedPasswordBytes, passwordBytes); err != nil {
    return err
	}

	return nil
}

func (u *UserUseCase) ChangePassword(ctx context.Context, id int, newPass string) error {
	hashedPassword, err := hashPassword(newPass)

	if err != nil {
		return err
	}

	newData := make(map[string]interface{}, 1)
	newData["password"] = hashedPassword

	if err := u.UserRepository.UpdateUser(ctx, id, newData); err != nil {
		return err
	}
	return nil
}

func (u *UserUseCase) UpdateUser(ctx context.Context,id int, newData map[string]interface{}) error {
	return u.UserRepository.UpdateUser(ctx, id, newData)
}	

func (u *UserUseCase) DeleteUser(ctx context.Context,id int) error {
	return u.UserRepository.DeleteUser(ctx, id)
}	

func hashPassword(password string) (string, error) {
	bytes, encryptErr := bcrypt.GenerateFromPassword([]byte(password), 10)
	if encryptErr != nil {
    return "", encryptErr
	}
	hashedPassword := string(bytes)

	return hashedPassword, nil
}
