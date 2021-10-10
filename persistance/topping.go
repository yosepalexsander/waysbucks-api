package persistance

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/yosepalexsander/waysbucks-api/entity"
)

type ToppingRepo struct {
	DB *sqlx.DB
}

func (s ToppingRepo) FindToppings(ctx context.Context) ([]entity.ProductTopping, error)  {
	sql, _, _ := sq.
	Select("id", "name", "image", "price", "is_available").
	From("toppings").ToSql()
	
	var toppings []entity.ProductTopping

	rows, err := s.DB.QueryxContext(ctx, sql)

	for rows.Next(){
		var topping entity.ProductTopping
		err = rows.StructScan(&topping)

		toppings = append(toppings, topping)
	}

	if err != nil {
		return nil, err
	}

	return toppings, nil
}

func (s ToppingRepo) FindTopping(ctx context.Context, id int) (*entity.ProductTopping, error)  {
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

func (s ToppingRepo) SaveTopping(ctx context.Context, topping entity.ProductTopping) error  {
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

func (s ToppingRepo) UpdateTopping(ctx context.Context, id int, newData map[string]interface{}) error {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	sql, args, _ := psql.Update("toppings").SetMap(newData).Where(sq.Eq{"id": id}).ToSql()

	_, err := s.DB.ExecContext(ctx, sql, args...)
	if err != nil {
		return err
	}

	return nil
}

func (s ToppingRepo) DeleteTopping(ctx context.Context, id int) error {
	sql, _, _ := sq.Delete("toppings").Where("id=$1").ToSql()
	
	_, err := s.DB.ExecContext(ctx, sql, id)
	if err != nil {
		return  err
	}

	return nil
} 