package storage

import (
	"context"
	"log"
	"sync"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/yosepalexsander/waysbucks-api/domain"
)
type UserStorage struct {
	DB *sqlx.DB
}

type (
	UserFinder interface {
		FindUserById(ctx context.Context, id uint64) (*domain.User, error)
		FindUserByEmail(ctx context.Context, email string) (*domain.User, error)
	}
	UserSaver interface {
		SaveUser(ctx context.Context, user domain.User) error
		UpdateUser(ctx context.Context, user domain.User) (*domain.User, error)
	}
	UserDelete interface {
		DeleteUser(ctx context.Context, id uint64) error
	}
)

func (storage UserStorage) FindUserById(ctx context.Context, id uint64) (*domain.User, error) {
	var wg sync.WaitGroup
	var err error
	var user domain.User

	wg.Add(1)
	go func ()  {
		defer wg.Done()
		sql, _, _ := sq.
		Select("id", "name", "email", "password", "gender", "phone", "image", "is_admin").
		From("users").Where("id=$1").ToSql()
		err = storage.DB.QueryRowxContext(ctx, sql, id).StructScan(&user)
	}()
	wg.Wait()

	if err != nil {
		log.Println(err)
		return &user, err
	}
	
	return &user, nil
}

func (storage UserStorage) FindUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	var wg sync.WaitGroup
	var err error
	var user domain.User
	
	wg.Add(1)
	go func ()  {
		defer wg.Done()
		sql, _, _ := sq.
		Select("id", "name", "email", "password", "gender", "phone", "image", "is_admin").
		From("users").Where("email=$1").ToSql()
		err = storage.DB.QueryRowxContext(ctx, sql, email).StructScan(&user)
	}()
	wg.Wait()

	if err != nil {
		log.Println(err)
		return &user, err
	}
	
	return &user, nil
}

func (storage UserStorage) SaveUser(ctx context.Context, user domain.User) error {
	var wg sync.WaitGroup
	var err error
	wg.Add(1)
	go func() {
		defer wg.Done()
		psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
		sql, args, _ := psql.
		Insert("users").
		Columns("name", "email", "password", "gender", "phone", "image", "is_admin").
		Values(user.Name, user.Email, user.Password, user.Gender, user.Phone, user.Image, user.IsAdmin).ToSql()
		_, err = storage.DB.ExecContext(ctx, sql, args...)
	}()
	wg.Wait()

	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (storage UserStorage) UpdateUser(ctx context.Context, user domain.User) (*domain.User, error) {
	
	return nil, nil
}

func (storage UserStorage) DeleteUser(ctx context.Context, id uint64) error {
	var wg sync.WaitGroup
	var err error
	wg.Add(1)
	go func() {
		defer wg.Done()
		sql, _, _ := sq.
		Delete("users").Where("id=$1").ToSql()
		_, err = storage.DB.ExecContext(ctx, sql, id)
	}()
	wg.Wait()

	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

