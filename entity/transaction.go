package entity

type Transaction struct {
	Id         int    `db:"id"`
	User_Id    int    `db:"user_id" json:"-"`
	Name       string `db:"name" json:"name"`
	Address    string `db:"address" json:"address"`
	PostalCode int    `db:"postal_code" json:"postal_code"`
	Phone      string `db:"phone" json:"phone"`
	Total      int    `db:"total" json:"total"`
	Status     string `db:"status" json:"status"`
	Orders     []Order `json:"orders"`
}

type Order struct {
	Id             int     `db:"id" json:"id"`
	Transaction_Id int     `db:"transaction_id" json:"-"`
	Product_Id     int     `db:"product_id" json:"product_id,omitempty"`
	Topping_Ids    []int64 `db:"topping_id" json:"topping_id,omitempty"`
	OrderProduct
	Price    int `db:"price" json:"price"`
	Qty      int `db:"qty" json:"qty"`
	Toppings []OrderTopping `json:"toppings"`
}

type TransactionTxParams struct {
	Transaction Transaction
	Order       []Order
}

type OrderProduct struct {
	Name  string `db:"name" json:"name"`
	Image string `db:"image" json:"image"`
}

type OrderTopping struct {
	Id   int    `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
}

// HTTP models for transaction
type OrderRequest struct {
	Product_Id  int     `json:"product_id" validate:"required"`
	Qty         int     `json:"qty" validate:"required"`
	Price       int     `json:"price" validate:"required"`
	Topping_Ids []int64 `json:"topping_id"`
}

type TransactionRequest struct {
	User_Id    int
	Name       string         `json:"name" validate:"required"`
	Address    string         `json:"address" validate:"required"`
	PostalCode int            `json:"postal_code" validate:"required"`
	Phone      string         `json:"phone" validate:"required"`
	Total      int            `json:"total" validate:"required"`
	Status     string         `json:"status" validate:"required"`
	Order      []OrderRequest `json:"orders" validate:"required"`
}