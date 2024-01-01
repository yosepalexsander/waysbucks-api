package persistance

import (
	"context"
	dbSql "database/sql"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/yosepalexsander/waysbucks-api/entity"
	"github.com/yosepalexsander/waysbucks-api/repository"
)

type cartRepo struct {
	db *sqlx.DB
}

func NewCartRepository(db *sqlx.DB) repository.CartRepository {
	return &cartRepo{db}
}

func (storage *cartRepo) FindCarts(ctx context.Context, userID string) ([]entity.Cart, error) {
	sql, _, _ := sq.Select("id", "product_id", "topping_id", "price", "qty").From("carts").Where("user_id=$1").OrderByClause("id DESC").ToSql()
	productSql, _, _ := sq.Select("id", "name", "image", "price").From("products").Where("id=$1").ToSql()
	toppingSql, _, _ := sq.Select("id", "name").From("toppings").Where("id IN $1").ToSql()

	carts := []entity.Cart{}

	rows, err := storage.db.QueryxContext(ctx, sql, userID)
	if err != nil {
		if err == dbSql.ErrNoRows {
			return carts, nil
		}

		return nil, err
	}

	for rows.Next() {
		var cart entity.Cart
		err = rows.Scan(&cart.Id, &cart.ProductId, pq.Array(&cart.ToppingIds), &cart.Price, &cart.Qty)
		_ = storage.db.QueryRowxContext(ctx, productSql, cart.ProductId).StructScan(&cart.Product)

		if len(cart.ToppingIds) < 1 {
			cart.Topping = make([]entity.CartTopping, 0)
		} else {
			rows, err := storage.db.QueryxContext(ctx, toppingSql, pq.Array(cart.ToppingIds))
			if err != nil {
				return nil, err
			}

			for rows.Next() {
				var topping entity.CartTopping

				err = rows.StructScan(&topping)
				if err != nil {
					return nil, err
				}

				cart.Topping = append(cart.Topping, topping)
			}
		}

		carts = append(carts, cart)
	}

	return carts, nil
}

func (storage *cartRepo) SaveCart(ctx context.Context, cart entity.Cart) error {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	sql, args, _ := psql.Insert("carts").
		Columns("user_id", "product_id", "price", "qty", "topping_id").
		Values(cart.UserId, cart.ProductId, cart.Price, cart.Qty, pq.Array(cart.ToppingIds)).ToSql()

	_, err := storage.db.ExecContext(ctx, sql, args...)
	if err != nil {
		return err
	}

	return nil
}

func (storage *cartRepo) UpdateCart(ctx context.Context, id int, userID string, data map[string]interface{}) error {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	sql, args, _ := psql.Update("carts").SetMap(data).Where(sq.Eq{"id": id, "user_id": userID}).ToSql()

	_, err := storage.db.ExecContext(ctx, sql, args...)
	if err != nil {
		return err
	}

	return nil
}

func (storage *cartRepo) DeleteCart(ctx context.Context, id int, userID string) error {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	sql, args, _ := psql.Delete("carts").Where(sq.Eq{"id": id, "user_id": userID}).ToSql()

	_, err := storage.db.ExecContext(ctx, sql, args...)

	if err != nil {
		return err
	}

	return nil
}
