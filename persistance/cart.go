package persistance

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/yosepalexsander/waysbucks-api/entity"
)

type CartRepo struct {
	DB *sqlx.DB
}

func (storage CartRepo) FindCarts(ctx context.Context, userID int) ([]entity.Cart, error) {
	sql, _, _ := sq.Select("id", "product_id", "topping_id", "price", "qty").From("carts").Where("user_id=$1").ToSql()
	productSql, _, _ := sq.Select("id", "name", "image", "price").From("products").Where("id=$1").ToSql()
	toppingSql, _, _ := sq.Select("id", "name").From("toppings").Where("id= $1").ToSql()

	carts := []entity.Cart{}

	rows, err := storage.DB.QueryxContext(ctx, sql, userID)

	for rows.Next() {
		var cart entity.Cart
		err = rows.Scan(&cart.Id, &cart.Product_Id, pq.Array(&cart.ToppingIds), &cart.Price, &cart.Qty)
		_ = storage.DB.QueryRowxContext(ctx, productSql, cart.Product_Id).StructScan(&cart.Product)
		for _, v := range cart.ToppingIds {
			var topping entity.CartTopping
			toppingErr := storage.DB.QueryRowxContext(ctx, toppingSql, v).StructScan(&topping)
			if toppingErr != nil {
				err = toppingErr
			}
			cart.Topping = append(cart.Topping, topping)
		}

		carts = append(carts, cart)
	}
	if err != nil {
		return nil, err
	}
	return carts, nil
}

func (storage CartRepo) SaveToCart(ctx context.Context, cart entity.Cart) error {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	sql, args, _ := psql.Insert("carts").
		Columns("user_id", "product_id", "price", "qty", "topping_id").
		Values(cart.User_Id, cart.Product_Id, cart.Price, cart.Qty, pq.Array(cart.ToppingIds)).ToSql()

	_, err := storage.DB.ExecContext(ctx, sql, args...)
	if err != nil {
		return err
	}

	return nil
}

func (storage CartRepo) UpdateCart(ctx context.Context, id int, userID int, data map[string]interface{}) error {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	sql, args, _ := psql.Update("carts").SetMap(data).Where(sq.Eq{"id": id, "user_id": userID}).ToSql()

	_, err := storage.DB.ExecContext(ctx, sql, args...)
	if err != nil {
		return err
	}

	return nil
}

func (storage CartRepo) DeleteCart(ctx context.Context, id int, userID int) error {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	sql, args, _ := psql.Delete("carts").Where(sq.Eq{"id": id, "user_id": userID}).ToSql()

	_, err := storage.DB.ExecContext(ctx, sql, args...)

	if err != nil {
		return err
	}

	return nil
}
