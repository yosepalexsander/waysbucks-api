package persistance

import (
	"context"
	"log"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/yosepalexsander/waysbucks-api/entity"
)

type ProductRepo struct {
	DB *sqlx.DB
}

func (storage ProductRepo) FindProducts(ctx context.Context) ([]entity.Product, error) {
	sql, _, _ := sq.Select("id", "name", "description", "image", "price", "is_available").From("products").OrderByClause("created_at DESC").ToSql()

	var products []entity.Product

	rows, err := storage.DB.QueryxContext(ctx, sql)

	if err != nil {
		return nil, err
	}
	for rows.Next() {
		product := entity.Product{}

		if err := rows.StructScan(&product); err != nil {
			log.Println(err)
		}
		
		products = append(products, product)
	}

	return products, nil
}

func (storage ProductRepo) FindProduct(ctx context.Context, id int) (*entity.Product, error) {
	sql, _, _ := sq.
	Select("name", "description", "image", "price", "is_available").
	From("products").
	Where("id=$1").ToSql()

	var product entity.Product
	err := storage.DB.QueryRowxContext(ctx, sql, id).StructScan(&product)
	if err != nil {
		return nil, err
	}

	return &product, nil
}

func (storage ProductRepo) SaveProduct(ctx context.Context, product entity.Product) error {
	sql, args, _ := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
	Insert("products").
	Columns("name", "description", "image", "price", "is_available").
	Values(product.Name, product.Description, product.Image, product.Price, product.IsAvailable).ToSql()

	_, err := storage.DB.ExecContext(ctx, sql, args...)

	if err != nil {
		return err
	}
	
	return nil
}

func (storage ProductRepo) UpdateProduct(ctx context.Context, id int, newProduct map[string]interface{}) error {
	sql, args, _ := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
	Update("products").SetMap(newProduct).
	Where(sq.Eq{"id": id}).ToSql()

	_, err := storage.DB.ExecContext(ctx, sql, args...)

	if err != nil {
		log.Println("update error: ", err)
		return err
	}
	
	return nil
}

func (storage ProductRepo) DeleteProduct(ctx context.Context, id int) error {
	sql, _, _ := sq.Delete("products").Where("id=$1").ToSql()

	_, err := storage.DB.ExecContext(ctx, sql, id)

	if err != nil {
		return err
	}
	
	return nil
}