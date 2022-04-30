package persistance

import (
	"context"
	"database/sql"
	"encoding/json"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/yosepalexsander/waysbucks-api/entity"
	"github.com/yosepalexsander/waysbucks-api/repository"
)

type transactionRepo struct {
	db *sqlx.DB
}

type sqlConnTx struct {
	db *sql.Tx
}

func NewTransactionRepository(db *sqlx.DB) repository.TransactionRepository {
	return &transactionRepo{db}
}

func (storage *transactionRepo) FindTransactions(ctx context.Context) ([]entity.Transaction, error) {
	sql, _, _ := sq.Select("t.id", "t.name", "t.address", "t.phone", "t.city", "t.postal_code", "t.total", "t.status",
		"json_agg(json_build_object('id', o.id, 'name', p.name,'image', p.image, 'topping_id', o.topping_id, 'price', o.price, 'qty', o.qty) ORDER BY o.id) AS order").
		From("transactions AS t, orders AS o, products AS p").Where("t.id = o.transaction_id AND o.product_id = p.id").GroupBy("t.id").
		OrderByClause("t.created_at DESC").ToSql()

	toppingSql, _, _ := sq.Select("id", "name").From("toppings").Where("id=$1").ToSql()

	var transactions []entity.Transaction
	rows, err := storage.db.QueryxContext(ctx, sql)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var t entity.Transaction
		var orderJSON []byte
		if err = rows.Scan(&t.Id, &t.Name, &t.Address, &t.Phone, &t.City, &t.PostalCode, &t.Total, &t.Status, &orderJSON); err != nil {
			return nil, err
		}
		_ = json.Unmarshal(orderJSON, &t.Orders)
		for i := range t.Orders {
			for _, v := range t.Orders[i].Topping_Ids {
				var topping entity.OrderTopping
				if err = storage.db.QueryRowxContext(ctx, toppingSql, v).Scan(&topping.Id, &topping.Name); err != nil {
					return nil, err
				}

				t.Orders[i].Toppings = append(t.Orders[i].Toppings, topping)
			}
			t.Orders[i].Topping_Ids = nil
		}
		transactions = append(transactions, t)
	}
	return transactions, nil
}

func (storage *transactionRepo) FindUserTransactions(ctx context.Context, userID string) ([]entity.Transaction, error) {
	sql, _, _ := sq.Select("t.id", "t.name", "t.address", "t.phone", "t.city", "t.postal_code", "t.total", "t.status",
		"json_agg(json_build_object('id', o.id, 'name', p.name,'image', p.image, 'topping_id', o.topping_id, 'price', o.price, 'qty', o.qty) ORDER BY o.id) AS order").
		From("transactions AS t, orders AS o, products AS p").Where("t.id = o.transaction_id AND t.user_id = $1 AND o.product_id = p.id").GroupBy("t.id").
		OrderByClause("t.created_at DESC").ToSql()

	toppingSql, _, _ := sq.Select("id", "name").From("toppings").Where("id=$1").ToSql()

	var transactions []entity.Transaction
	rows, err := storage.db.QueryxContext(ctx, sql, userID)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var t entity.Transaction
		var orderJSON []byte
		if err = rows.Scan(&t.Id, &t.Name, &t.Address, &t.Phone, &t.City, &t.PostalCode, &t.Total, &t.Status, &orderJSON); err != nil {
			return nil, err
		}
		_ = json.Unmarshal(orderJSON, &t.Orders)
		for i := range t.Orders {
			for _, v := range t.Orders[i].Topping_Ids {
				var topping entity.OrderTopping
				if err = storage.db.QueryRowxContext(ctx, toppingSql, v).Scan(&topping.Id, &topping.Name); err != nil {
					return nil, err
				}

				t.Orders[i].Toppings = append(t.Orders[i].Toppings, topping)
			}
			t.Orders[i].Topping_Ids = nil
		}
		transactions = append(transactions, t)
	}
	return transactions, nil
}

func (storage *transactionRepo) FindTransactionByID(ctx context.Context, id string) (*entity.Transaction, error) {
	sql, _, _ := sq.Select("t.id", "t.name", "t.address", "t.phone", "t.city", "t.postal_code", "t.total", "t.status",
		"json_agg(json_build_object('id', o.id, 'name', p.name,'image', p.image, 'topping_id', o.topping_id, 'price', o.price, 'qty', o.qty) ORDER BY o.id) AS order").
		From("transactions AS t, orders AS o, products AS p").Where("t.id = $1 AND t.id = o.transaction_id AND o.product_id = p.id").GroupBy("t.id").
		OrderByClause("t.created_at DESC").ToSql()
	toppingSql, _, _ := sq.Select("id", "name").From("toppings").Where("id=$1").ToSql()

	var t entity.Transaction
	var orderJSON []byte
	row := storage.db.QueryRowxContext(ctx, sql, id)
	if err := row.Scan(&t.Id, &t.Name, &t.Address, &t.Phone, &t.City, &t.PostalCode, &t.Total, &t.Status, &orderJSON); err != nil {
		return nil, err
	}
	_ = json.Unmarshal(orderJSON, &t.Orders)
	for i := range t.Orders {
		for _, v := range t.Orders[i].Topping_Ids {
			var topping entity.OrderTopping
			if err := storage.db.QueryRowxContext(ctx, toppingSql, v).Scan(&topping.Id, &topping.Name); err != nil {
				return nil, err
			}

			t.Orders[i].Toppings = append(t.Orders[i].Toppings, topping)
		}
		t.Orders[i].Topping_Ids = nil
	}
	return &t, nil
}

func (storage *transactionRepo) UpdateTransaction(ctx context.Context, id string, data map[string]interface{}) error {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	sql, args, _ := psql.Update("transactions").SetMap(data).ToSql()

	_, err := storage.db.ExecContext(ctx, sql, args...)
	if err != nil {
		return err
	}

	return nil
}

func (storage *transactionRepo) TxBegin(ctx context.Context) (repository.Transactioner, error) {
	tx, err := storage.db.BeginTx(ctx, nil)
	sct := sqlConnTx{tx}
	return &sct, err
}

func (storage *transactionRepo) ExecTx(ctx context.Context, fn func(repository.Transactioner) error) error {
	tx, err := storage.TxBegin(ctx)
	if err != nil {
		return err
	}

	err = fn(tx)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return rbErr
		}
		return err
	}

	return tx.Commit()
}

func (sct *sqlConnTx) CreateTransaction(ctx context.Context, tx entity.Transaction) (string, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	var id string
	sql, args, _ := psql.Insert("transactions").Columns("id", "user_id", "name", "address", "city", "postal_code", "phone", "total", "status").
		Values(tx.Id, tx.UserId, tx.Name, tx.Address, tx.City, tx.PostalCode, tx.Phone, tx.Total, tx.Status).Suffix("RETURNING id").ToSql()

	err := sct.db.QueryRowContext(ctx, sql, args...).Scan(&id)

	if err != nil {
		return "", err
	}
	return id, nil
}

func (sct *sqlConnTx) CreateOrder(ctx context.Context, order entity.Order) error {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	var err error
	sql, args, _ := psql.Insert("orders").Columns("transaction_id", "product_id", "topping_id", "price", "qty").
		Values(order.Transaction_Id, order.ProductId, pq.Array(order.Topping_Ids), order.Price, order.Qty).ToSql()

	_, err = sct.db.ExecContext(ctx, sql, args...)
	return err
}

func (sct *sqlConnTx) DeleteCart(ctx context.Context, productID int, userID string) error {
	var err error
	sql, _, _ := sq.Delete("carts").Where("product_id=$1 AND user_id=$2").ToSql()

	_, err = sct.db.ExecContext(ctx, sql, productID, userID)
	return err
}

func (sct *sqlConnTx) Rollback() error {
	return sct.db.Rollback()
}

func (sct *sqlConnTx) Commit() error {
	return sct.db.Commit()
}
