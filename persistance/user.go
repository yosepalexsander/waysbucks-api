package persistance

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/yosepalexsander/waysbucks-api/entity"
	"github.com/yosepalexsander/waysbucks-api/repository"
)

type userRepo struct {
	DB *sqlx.DB
}

func NewUserFinder(DB *sqlx.DB) repository.UserFinder {
	return &userRepo{DB}
}

func NewUserMutator(DB *sqlx.DB) repository.UserMutator {
	return &userRepo{DB}
}

func (storage *userRepo) FindUsers(ctx context.Context) ([]entity.User, error) {
	sql, _, _ := sq.
		Select("id", "name", "email", "gender", "phone", "image").
		From("users").Where("is_admin = $1").ToSql()

	users := []entity.User{}

	rows, err := storage.DB.QueryxContext(ctx, sql, 0)
	for rows.Next() {
		user := entity.User{}
		err = rows.StructScan(&user)
		users = append(users, user)
	}

	if err != nil {
		return nil, err
	}

	return users, nil
}
func (storage *userRepo) FindUserById(ctx context.Context, id int) (*entity.User, error) {
	var user entity.User

	sql, _, _ := sq.
		Select("id", "name", "email", "gender", "phone", "image", "is_admin").
		From("users").Where("id=$1").ToSql()
	err := storage.DB.QueryRowxContext(ctx, sql, id).StructScan(&user)

	if err != nil {
		return &user, err
	}

	return &user, nil
}

func (storage *userRepo) FindUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	var user entity.User

	sql, _, _ := sq.
		Select("id", "name", "email", "password", "gender", "phone", "image", "is_admin").
		From("users").Where("email=$1").ToSql()
	err := storage.DB.QueryRowxContext(ctx, sql, email).StructScan(&user)

	if err != nil {
		return &user, err
	}

	return &user, nil
}

func (storage *userRepo) SaveUser(ctx context.Context, user entity.User) error {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	sql, args, _ := psql.
		Insert("users").
		Columns("name", "email", "password", "gender", "phone", "image", "is_admin").
		Values(user.Name, user.Email, user.Password, user.Gender, user.Phone, user.Image, user.IsAdmin).ToSql()
	_, err := storage.DB.ExecContext(ctx, sql, args...)

	if err != nil {
		return err
	}

	return nil
}

func (storage *userRepo) UpdateUser(ctx context.Context, id int, newData map[string]interface{}) error {

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	sql, args, _ := psql.
		Update("users").SetMap(newData).
		Where(sq.Eq{"id": id}).ToSql()

	_, err := storage.DB.ExecContext(ctx, sql, args...)

	if err != nil {
		return err
	}

	return nil
}

func (storage *userRepo) DeleteUser(ctx context.Context, id int) error {
	sql, _, _ := sq.
		Delete("users").Where("id=$1").ToSql()
	_, err := storage.DB.ExecContext(ctx, sql, id)

	if err != nil {
		return err
	}

	return nil
}
