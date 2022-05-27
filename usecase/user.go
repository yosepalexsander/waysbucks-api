package usecase

import (
	"context"

	"github.com/google/uuid"
	"github.com/yosepalexsander/waysbucks-api/entity"
	"github.com/yosepalexsander/waysbucks-api/repository"
	"github.com/yosepalexsander/waysbucks-api/thirdparty"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/sync/errgroup"
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
	user, err := u.repo.FindUserById(ctx, id)
	if err != nil {
		return nil, err
	}

	imageUrl, _ := thirdparty.GetImageUrl(ctx, user.Image)
	user.Image = imageUrl

	return user, nil
}

func (u *UserUseCase) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	return u.repo.FindUserByEmail(ctx, email)
}

func (u *UserUseCase) CreateNewUser(ctx context.Context, name string, email string, password string, gender string, phone string) (*entity.User, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	hashedPassword, err := hashPassword(password)
	if err != nil {
		return nil, err
	}

	user := entity.NewUser(name, email, hashedPassword, gender, phone)
	user.Id = id.String()

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
	user, err := u.repo.FindUserById(ctx, id)
	if err != nil {
		return err
	}

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		return u.repo.UpdateUser(ctx, id, newData)
	})

	g.Go(func() error {
		if newImage, ok := newData["image"]; ok && newImage != user.Image {
			return thirdparty.RemoveFile(ctx, user.Image)
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		return err
	}

	return nil
}

func (u *UserUseCase) DeleteUser(ctx context.Context, id string) error {
	user, err := u.repo.FindUserById(ctx, id)
	if err != nil {
		return err
	}

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		return u.repo.DeleteUser(ctx, id)
	})

	g.Go(func() error {
		return thirdparty.RemoveFile(ctx, user.Image)
	})

	if err := g.Wait(); err != nil {
		return err
	}

	return nil
}

func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}
