package usecase

import (
	"context"
	"mime/multipart"
	"sync"

	"github.com/yosepalexsander/waysbucks-api/entity"
	"github.com/yosepalexsander/waysbucks-api/repository"
	"github.com/yosepalexsander/waysbucks-api/thirdparty"
	"golang.org/x/crypto/bcrypt"
)

type UserUseCase struct {
	Finder  repository.UserFinder
	Mutator repository.UserMutator
}

func NewUserUseCase(rf repository.UserFinder, rm repository.UserMutator) UserUseCase {
	return UserUseCase{rf, rm}
}

func (u *UserUseCase) FindUsers(ctx context.Context) ([]entity.User, error) {
	return u.Finder.FindUsers(ctx)
}
func (u *UserUseCase) FindUserById(ctx context.Context, id int) (*entity.User, error) {
	user, err := u.Finder.FindUserById(ctx, id)

	if err != nil {
		return nil, err
	}
	imageUrl, err := thirdparty.GetImageUrl(ctx, user.Image)
	if imageUrl != "" || err == nil {
		user.Image = imageUrl
	}
	return user, nil
}

func (u *UserUseCase) FindUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	return u.Finder.FindUserByEmail(ctx, email)
}

func (u *UserUseCase) CreateNewUser(ctx context.Context, user entity.User) error {
	hashedPassword, err := hashPassword(user.Password)

	if err != nil {
		return err
	}

	user.Password = hashedPassword

	if err := u.Mutator.SaveUser(ctx, user); err != nil {
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

	if err := u.Mutator.UpdateUser(ctx, id, newData); err != nil {
		return err
	}
	return nil
}

func (u *UserUseCase) UpdateUser(ctx context.Context, id int, newData map[string]interface{}) error {
	return u.Mutator.UpdateUser(ctx, id, newData)
}

func (u *UserUseCase) UpdateImage(ctx context.Context, file multipart.File, oldName string, newName string) error {
	wg := &sync.WaitGroup{}
	var uploadErr error
	wg.Add(2)
	go func() {
		defer wg.Done()
		if err := thirdparty.UploadFile(ctx, file, newName); err != nil {
			uploadErr = err
			return
		}
	}()
	go func() {
		defer wg.Done()
		if err := thirdparty.RemoveFile(ctx, oldName); err != nil {
			uploadErr = err
			return
		}
	}()
	wg.Wait()

	if uploadErr != nil {
		return uploadErr
	}

	return nil
}
func (u *UserUseCase) DeleteUser(ctx context.Context, id int) error {
	return u.Mutator.DeleteUser(ctx, id)
}

func hashPassword(password string) (string, error) {
	bytes, encryptErr := bcrypt.GenerateFromPassword([]byte(password), 10)
	if encryptErr != nil {
		return "", encryptErr
	}
	hashedPassword := string(bytes)

	return hashedPassword, nil
}
