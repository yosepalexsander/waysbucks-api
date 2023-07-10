package persistance

import (
	"context"
	dbSql "database/sql"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/yosepalexsander/waysbucks-api/entity"
	"github.com/yosepalexsander/waysbucks-api/repository"
)

type productRepo struct {
	db *sqlx.DB
}

func NewProductRepository(db *sqlx.DB) repository.ProductRepository {
	return &productRepo{db}
}

func (storage *productRepo) FindProducts(ctx context.Context, whereClauses []string, orderClause string) ([]entity.Product, error) {
	sq := sq.Select("id", "name", "description", "image", "price", "is_available", "created_at", "updated_at").
		From("products")

	for _, v := range whereClauses {
		sq = sq.Where(v)
	}

	if orderClause != "" {
		sq = sq.OrderByClause(orderClause)
	} else {
		sq = sq.OrderByClause("created_at DESC")
	}

	sql, _, _ := sq.ToSql()

	products := []entity.Product{}

	rows, err := storage.db.QueryxContext(ctx, sql)
	if err != nil {
		if err == dbSql.ErrNoRows {
			return products, nil
		}

		return nil, err
	}

	for rows.Next() {
		product := entity.Product{}
		err = rows.StructScan(&product)
		products = append(products, product)
	}

	return products, nil
}

func (storage *productRepo) FindProduct(ctx context.Context, id int) (*entity.Product, error) {
	sql, _, _ := sq.
		Select("id", "name", "description", "image", "price", "is_available").
		From("products").
		Where("id=$1").ToSql()

	var product entity.Product
	err := storage.db.QueryRowxContext(ctx, sql, id).StructScan(&product)
	if err != nil {
		return nil, err
	}

	return &product, nil
}

func (storage *productRepo) SaveProduct(ctx context.Context, product entity.Product) error {
	sql, args, _ := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Insert("products").
		Columns("name", "description", "image", "price", "is_available").
		Values(product.Name, product.Description, product.Image, product.Price, product.IsAvailable).ToSql()

	_, err := storage.db.ExecContext(ctx, sql, args...)

	if err != nil {
		return err
	}

	return nil
}

func (storage *productRepo) UpdateProduct(ctx context.Context, id int, newProduct map[string]interface{}) error {
	sql, args, _ := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Update("products").SetMap(newProduct).
		Where(sq.Eq{"id": id}).ToSql()

	_, err := storage.db.ExecContext(ctx, sql, args...)

	if err != nil {
		return err
	}

	return nil
}

func (storage *productRepo) DeleteProduct(ctx context.Context, id int) error {
	sql, _, _ := sq.Delete("products").Where("id=$1").ToSql()

	_, err := storage.db.ExecContext(ctx, sql, id)

	if err != nil {
		return err
	}

	return nil
}

func (s *productRepo) FindToppings(ctx context.Context) ([]entity.ProductTopping, error) {
	sql, _, _ := sq.
		Select("id", "name", "image", "price", "is_available").
		From("toppings").OrderByClause("created_at DESC").ToSql()

	toppings := []entity.ProductTopping{}

	rows, err := s.db.QueryxContext(ctx, sql)

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

func (s *productRepo) FindTopping(ctx context.Context, id int) (*entity.ProductTopping, error) {
	sql, _, _ := sq.
		Select("id", "name", "image", "price", "is_available").
		From("toppings").Where("id=$1").ToSql()

	var topping entity.ProductTopping

	err := s.db.QueryRowxContext(ctx, sql, id).StructScan(&topping)
	if err != nil {
		return nil, err
	}

	return &topping, nil
}

func (s *productRepo) SaveTopping(ctx context.Context, topping entity.ProductTopping) error {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	sql, args, _ := psql.Insert("toppings").
		Columns("name", "image", "price", "is_available").
		Values(topping.Name, topping.Image, topping.Price, topping.IsAvailable).ToSql()

	_, err := s.db.ExecContext(ctx, sql, args...)
	if err != nil {
		return err
	}

	return nil
}

func (s *productRepo) UpdateTopping(ctx context.Context, id int, newData map[string]interface{}) error {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	sql, args, _ := psql.Update("toppings").SetMap(newData).Where(sq.Eq{"id": id}).ToSql()

	_, err := s.db.ExecContext(ctx, sql, args...)
	if err != nil {
		return err
	}

	return nil
}

func (s *productRepo) DeleteTopping(ctx context.Context, id int) error {
	sql, _, _ := sq.Delete("toppings").Where("id=$1").ToSql()

	_, err := s.db.ExecContext(ctx, sql, id)
	if err != nil {
		return err
	}

	return nil
}
