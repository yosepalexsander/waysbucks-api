package persistance

import (
	"context"
	"log"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/yosepalexsander/waysbucks-api/domain"
)

type UserRepo struct {
	DB *sqlx.DB
}

type (
	UserFinder interface {
		FindUserById(ctx context.Context, id uint64) (*domain.User, error)
		FindUserByEmail(ctx context.Context, email string) (*domain.User, error)
	}
	UserSaver interface {
		SaveUser(ctx context.Context, user domain.User) error
		UpdateUser(ctx context.Context,id uint64, newData map[string]interface{}) (*domain.User, error)
	}
	UserRemover interface {
		DeleteUser(ctx context.Context, id uint64) error
	}
)

func (storage UserRepo) FindUserById(ctx context.Context, id uint64) (*domain.User, error) {
	var user domain.User
	
	sql, _, _ := sq.
	Select("id", "name", "email", "gender", "phone", "image", "is_admin").
	From("users").Where("id=$1").ToSql()
	err := storage.DB.QueryRowxContext(ctx, sql, id).StructScan(&user)

	if err != nil {
		log.Println(err)
		return &user, err
	}
	
	return &user, nil
}

func (storage UserRepo) FindUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User
	
	sql, _, _ := sq.
	Select("id", "name", "email", "password", "gender", "phone", "image", "is_admin").
	From("users").Where("email=$1").ToSql()
	err := storage.DB.QueryRowxContext(ctx, sql, email).StructScan(&user)

	if err != nil {
		log.Println(err)
		return &user, err
	}
	
	return &user, nil
}

func (storage UserRepo) SaveUser(ctx context.Context, user domain.User) error {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	sql, args, _ := psql.
	Insert("users").
	Columns("name", "email", "password", "gender", "phone", "image", "is_admin").
	Values(user.Name, user.Email, user.Password, user.Gender, user.Phone, user.Image, user.IsAdmin).ToSql()
	_, err := storage.DB.ExecContext(ctx, sql, args...)

	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (storage UserRepo) UpdateUser(ctx context.Context, id uint64, newData map[string]interface{}) (*domain.User, error) {

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	sql, args, _ := psql.
	Update("users").SetMap(newData).
	Where(sq.Eq{"id": id}).ToSql()

	_, err := storage.DB.ExecContext(ctx, sql, args...)

	if err != nil {
		log.Println(err)
		return nil, err
	}
	return nil, nil
}

func (storage UserRepo) DeleteUser(ctx context.Context, id uint64) error {
	sql, _, _ := sq.
	Delete("users").Where("id=$1").ToSql()
	_, err := storage.DB.ExecContext(ctx, sql, id)
	
	if err != nil {
		log.Println(err)
		return err
	}
	
	return nil
}

