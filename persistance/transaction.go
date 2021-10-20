package persistance

import (
	"context"
	"database/sql"
	"sync"

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

func NewTransactionFinder(db *sqlx.DB) repository.TransactionFinder {
	return &transactionRepo{db: db}
}

func NewTransactionTx(db *sqlx.DB) repository.TransactionTx {
	return &transactionRepo{db: db}
}

func NewTransactionMutator(db *sqlx.DB) repository.TransactionMutator {
	return &transactionRepo{db: db}
}

var (
	orderSql, _, _   = sq.Select("id", "product_id", "topping_id", "price", "qty").From("orders").Where("transaction_id=$1").ToSql()
	productSql, _, _ = sq.Select("name", "image").From("products").Where("id=$1").ToSql()
	toppingSql, _, _ = sq.Select("id", "name").From("toppings").Where("id=$1").ToSql()
)

func (storage *transactionRepo) FindTransactions(ctx context.Context) ([]entity.Transaction, error) {
	sql, _, _ := sq.Select("id", "name", "address", "postal_code", "city", "total", "status").From("transactions").ToSql()
	orderSql, _, _ := sq.Select("id", "product_id", "topping_id", "price", "qty").From("orders").Where("transaction_id=$1").ToSql()
	productSql, _, _ := sq.Select("name", "image").From("products").Where("id=$1").ToSql()
	toppingSql, _, _ := sq.Select("id", "name").From("toppings").Where("id=$1").ToSql()

	var transactions []entity.Transaction
	rows, err := storage.db.QueryxContext(ctx, sql)
	if err != nil {
		return nil, err
	}

	var wg sync.WaitGroup
	for rows.Next() {
		var transaction entity.Transaction
		if err = rows.StructScan(&transaction); err != nil {
			return nil, err
		}
		orderRows, err := storage.db.QueryxContext(ctx, orderSql, transaction.Id)
		if err != nil {
			return nil, err
		}

		for orderRows.Next() {
			var order entity.Order
			_ = orderRows.Scan(&order.Id, &order.Product_Id, pq.Array(&order.Topping_Ids), &order.Price, &order.Qty)
			_ = storage.db.QueryRowxContext(ctx, productSql, order.Product_Id).StructScan(&order.OrderProduct)
			for _, v := range order.Topping_Ids {
				wg.Add(1)
				go func(v int64) {
					defer wg.Done()
					var topping entity.OrderTopping
					if err = storage.db.QueryRowxContext(ctx, toppingSql, v).StructScan(&topping); err != nil {
						return
					}
					order.Toppings = append(order.Toppings, topping)
				}(v)
				wg.Wait()
			}
			transaction.Orders = append(transaction.Orders, order)
		}
		transactions = append(transactions, transaction)
	}

	return transactions, nil
}

func (storage *transactionRepo) FindUserTransactions(ctx context.Context, userID int) ([]entity.Transaction, error) {
	sql, _, _ := sq.Select("t.id", "t.name", "t.address", "t.postal_code", "t.city", "t.total", "t.status",
		"o.id", "o.product_id", "o.topping_id", "o.price", "o.qty").From("transactions AS t, orders AS o").
		Where("user_id=$1 AND t.id = o.id").ToSql()

	productSql, _, _ := sq.Select("name", "image").From("products").Where("id=$1").ToSql()
	toppingSql, _, _ := sq.Select("id", "name").From("toppings").Where("id=$1").ToSql()

	var transactions []entity.Transaction

	rows, err := storage.db.QueryxContext(ctx, sql, userID)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var t entity.Transaction
		var o entity.Order
		var topping entity.OrderTopping

		if err = rows.Scan(&t.Id, &t.Name, &t.Address, &t.PostalCode, &t.City, &t.Total, &t.Status,
			&o.Id, &o.Product_Id, pq.Array(&o.Topping_Ids), &o.Price, &o.Qty); err != nil {
			return nil, err
		}
		if err = storage.db.QueryRowxContext(ctx, productSql, o.Product_Id).StructScan(&o.OrderProduct); err != nil {
			return nil, err
		}
		for _, v := range o.Topping_Ids {
			if err = storage.db.QueryRowxContext(ctx, toppingSql, v).StructScan(&topping); err != nil {
				return nil, err
			}
			o.Toppings = append(o.Toppings, topping)
		}
		t.Orders = append(t.Orders, o)
		transactions = append(transactions, t)
	}

	return transactions, nil
}

func (storage *transactionRepo) FindTransactionByID(ctx context.Context, id int) (*entity.Transaction, error) {
	sql, _, _ := sq.Select("id", "name", "address", "postal_code", "city", "total", "status").From("transactions").Where("id=$1").ToSql()
	var transaction entity.Transaction
	if err := storage.db.QueryRowxContext(ctx, sql, id).StructScan(&transaction); err != nil {
		return nil, err
	}

	rows, err := storage.db.QueryxContext(ctx, orderSql, transaction.Id)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var order entity.Order
		if err = rows.Scan(&order.Id, &order.Product_Id, pq.Array(&order.Topping_Ids), &order.Price, &order.Qty); err != nil {
			return nil, err
		}
		if err = storage.db.QueryRowxContext(ctx, productSql, order.Product_Id).StructScan(&order.OrderProduct); err != nil {
			return nil, err
		}

		for _, v := range order.Topping_Ids {
			var topping entity.OrderTopping
			if err = storage.db.QueryRowxContext(ctx, toppingSql, v).StructScan(&topping); err != nil {
				return nil, err
			}

			order.Toppings = append(order.Toppings, topping)
		}
		transaction.Orders = append(transaction.Orders, order)
	}
	return &transaction, err
}

func (storage *transactionRepo) UpdateTransaction(ctx context.Context, id int, data map[string]interface{}) error {
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

func (storage *transactionRepo) OrderTx(ctx context.Context, arg entity.TransactionTxParams) error {
	txErr := storage.ExecTx(ctx, func(tx repository.Transactioner) error {
		var err error

		id, err := tx.CreateTransaction(ctx, arg.Transaction)
		if err != nil {
			return err
		}
		for i := range arg.Order {
			arg.Order[i].Transaction_Id = id
			err := tx.CreateOrder(ctx, arg.Order[i])
			if err != nil {
				return err
			}
			err = tx.DeleteCart(ctx, arg.Order[i].Product_Id, arg.Transaction.User_Id)
			if err != nil {
				return err
			}
		}

		return nil
	})

	if txErr != nil {
		return txErr
	}

	return nil
}

func (sct *sqlConnTx) CreateTransaction(ctx context.Context, tx entity.Transaction) (int, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	var id int
	sql, args, _ := psql.Insert("transactions").Columns("user_id", "name", "address", "postal_code", "city", "phone", "total", "status").
		Values(tx.User_Id, tx.Name, tx.Address, tx.PostalCode, tx.City, tx.Phone, tx.Total, tx.Status).Suffix("RETURNING id").ToSql()

	err := sct.db.QueryRowContext(ctx, sql, args...).Scan(&id)

	if err != nil {
		return 0, err
	}
	return id, nil
}

func (sct *sqlConnTx) CreateOrder(ctx context.Context, order entity.Order) error {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	var err error
	sql, args, _ := psql.Insert("orders").Columns("transaction_id", "product_id", "topping_id", "price", "qty").
		Values(order.Transaction_Id, order.Product_Id, pq.Array(order.Topping_Ids), order.Price, order.Qty).ToSql()

	_, err = sct.db.ExecContext(ctx, sql, args...)
	return err
}

func (sct *sqlConnTx) DeleteCart(ctx context.Context, productID int, userID int) error {
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
