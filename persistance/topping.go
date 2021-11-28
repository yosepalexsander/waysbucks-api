package persistance

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/yosepalexsander/waysbucks-api/entity"
	"github.com/yosepalexsander/waysbucks-api/repository"
)

type toppingRepo struct {
	DB *sqlx.DB
}

func NewToppingRepo(DB *sqlx.DB) repository.ToppingRepository {
	return &toppingRepo{DB}
}

func (s *toppingRepo) FindToppings(ctx context.Context) ([]entity.ProductTopping, error) {
	sql, _, _ := sq.
		Select("id", "name", "image", "price", "is_available").
		From("toppings").OrderByClause("created_at DESC").ToSql()

	var toppings []entity.ProductTopping

	rows, err := s.DB.QueryxContext(ctx, sql)

	for rows.Next() {
		var topping entity.ProductTopping
		err = rows.StructScan(&topping)

		toppings = append(toppings, topping)
	}

	if err != nil {
		return nil, err
	}

	return toppings, nil
}

func (s *toppingRepo) FindTopping(ctx context.Context, id int) (*entity.ProductTopping, error) {
	sql, _, _ := sq.
		Select("id", "name", "image", "price", "is_available").
		From("toppings").Where("id=$1").ToSql()

	var topping entity.ProductTopping

	err := s.DB.QueryRowxContext(ctx, sql, id).StructScan(&topping)
	if err != nil {
		return nil, err
	}

	return &topping, nil
}

func (s *toppingRepo) SaveTopping(ctx context.Context, topping entity.ProductTopping) error {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	sql, args, _ := psql.Insert("toppings").
		Columns("name", "image", "price", "is_available").
		Values(topping.Name, topping.Image, topping.Price, topping.Is_Available).ToSql()

	_, err := s.DB.ExecContext(ctx, sql, args...)
	if err != nil {
		return err
	}

	return nil
}

func (s *toppingRepo) UpdateTopping(ctx context.Context, id int, newData map[string]interface{}) error {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	sql, args, _ := psql.Update("toppings").SetMap(newData).Where(sq.Eq{"id": id}).ToSql()

	_, err := s.DB.ExecContext(ctx, sql, args...)
	if err != nil {
		return err
	}

	return nil
}

func (s *toppingRepo) DeleteTopping(ctx context.Context, id int) error {
	sql, _, _ := sq.Delete("toppings").Where("id=$1").ToSql()

	_, err := s.DB.ExecContext(ctx, sql, id)
	if err != nil {
		return err
	}

	return nil
}
