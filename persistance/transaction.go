package persistance

import (
	"context"
	"database/sql"
	"log"

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
	return &transactionRepo{
		db: db,
	}
}

func (storage transactionRepo) FindTransactions(ctx context.Context, userID int) ([]entity.Transaction, error)  {
	sql, _, _ := sq.Select("id", "name", "address", "postcode", "total", "status").From("transactions").Where("user_id=$1").ToSql()
	orderSql, _, _ := sq.Select("id", "product_id", "topping_id", "price", "qty").From("orders").Where("transaction_id=$1").ToSql()
	productSql, _, _ := sq.Select("name", "image").From("products").Where("id=$1").ToSql()
	toppingSql, _, _ := sq.Select("id", "name").From("toppings").Where("id=$1").ToSql()
	
	var transactions []entity.Transaction

	rows, err := storage.db.QueryxContext(ctx, sql, userID)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var transaction entity.Transaction
		var order entity.Order
		if err = rows.StructScan(&transaction); err != nil {
			return nil, err
		}
		err = storage.db.QueryRowxContext(ctx, orderSql, transaction.Id).Scan(&order.Id, &order.Product_Id, pq.Array(&order.Topping_Ids), &order.Price, &order.Qty)
		_ = storage.db.QueryRowxContext(ctx, productSql, order.Product_Id).StructScan(&order.OrderProduct)

		for _, v := range order.Topping_Ids {
			var topping entity.OrderTopping
			if err = storage.db.QueryRowxContext(ctx, toppingSql, v).StructScan(&topping); err != nil {
				return nil, err
			}
			
			order.Toppings = append(order.Toppings, topping)
		}

		transaction.Orders = append(transaction.Orders, order)
		transactions = append(transactions, transaction)
	}
	return transactions, err
}

func (storage transactionRepo) ExecTx(ctx context.Context, fn func(repository.Transactioner) error) error {
	tx, err := storage.TxBegin(ctx)
	if err != nil {
		return err
	}

	err = fn(tx)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			log.Println(rbErr)
			return rbErr
		}
		return err
	}

	return tx.Commit()
}

func (storage transactionRepo) OrderTx(ctx context.Context, arg entity.TransactionTxParams) error {
	txErr := storage.ExecTx(ctx, func(tx repository.Transactioner) error {
		var err error

		id, err := tx.CreateTransaction(ctx, arg.Transaction)
		if err != nil {
			return err
		} 
		for i := range arg.Order {
			arg.Order[i].Transaction_Id = id
		}
		
		err = tx.CreateOrder(ctx, arg.Order)
		if err != nil {
			return err
		} 

		err = tx.DeleteCart(ctx, arg.ProductIds, arg.Transaction.User_Id) 
		if err != nil {
			return err
		} 

		return nil
	})

	if txErr != nil {
		return txErr
	}

	return nil
}

func (sct *sqlConnTx) CreateTransaction(ctx context.Context, tx entity.Transaction) (int, error)  {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	
	var id int
	sql, args, _ := psql.Insert("transactions").Columns("user_id", "name", "address", "postcode", "phone", "total", "status").
	Values(tx.User_Id, tx.Name, tx.Address, tx.PostCode, tx.Phone, tx.Total, tx.Status).Suffix("RETURNING id").ToSql()
	
	err := sct.db.QueryRowContext(ctx, sql, args...).Scan(&id)
	
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (sct *sqlConnTx) CreateOrder(ctx context.Context, order []entity.Order) error  {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	var err error
	for _, v := range order {
		sql, args, _ := psql.Insert("orders").Columns("transaction_id", "product_id", "price", "qty").
		Values(v.Transaction_Id, v.Product_Id, v.Price, v.Qty).ToSql()		
		_, err = sct.db.ExecContext(ctx, sql, args...)
	}

	return err
}

func (sct *sqlConnTx) DeleteCart(ctx context.Context, productIds []int, userID int) error {
	var err error
	for _, v := range productIds {
		sql, _, _ := sq.Delete("carts").Where("product_id=$1 AND user_id=$2").ToSql()

		_, err = sct.db.ExecContext(ctx, sql, v, userID)
		
		if err != nil {
			return err
		}
	}
	return nil
}
func (storage transactionRepo) TxBegin(ctx context.Context) (repository.Transactioner, error) {
	tx, err := storage.db.BeginTx(ctx, nil)
	sct := sqlConnTx{tx}
	return &sct, err
}


func (sct *sqlConnTx) Rollback() error {
	return sct.db.Rollback()
}

func (sct *sqlConnTx) Commit() error {
	return sct.db.Commit()
}