package entity

type Transaction struct {
	Id         string  `db:"id"`
	Name       string  `db:"name" json:"name"`
	Email      string  `json:"email"`
	Phone      string  `db:"phone" json:"phone"`
	Address    string  `db:"address" json:"address"`
	City       string  `db:"city" json:"city"`
	PostalCode int     `db:"postal_code" json:"postal_code"`
	Total      int     `db:"total" json:"total"`
	ServiceFee int     `json:"service_fee"`
	Status     string  `db:"status" json:"status"`
	User_Id    string  `db:"user_id" json:"-"`
	Orders     []Order `json:"orders"`
}

type Order struct {
	Id             int     `db:"id" json:"id"`
	Price          int     `db:"price" json:"price"`
	Qty            int     `db:"qty" json:"qty"`
	Product_Id     int     `db:"product_id" json:"product_id,omitempty"`
	Topping_Ids    []int64 `db:"topping_id" json:"topping_id,omitempty"`
	Transaction_Id string  `db:"transaction_id" json:"-"`
	OrderProduct
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

type OrderRequest struct {
	Qty         int     `json:"qty" validate:"required"`
	Price       int     `json:"price" validate:"required"`
	Product_Id  int     `json:"product_id" validate:"required"`
	Topping_Ids []int64 `json:"topping_id"`
}

type TransactionRequest struct {
	Email      string         `json:"email" validate:"required"`
	Name       string         `json:"name" validate:"required"`
	Address    string         `json:"address" validate:"required"`
	City       string         `json:"city" validate:"required"`
	Phone      string         `json:"phone" validate:"required"`
	ServiceFee int            `json:"service_fee" validate:"required"`
	PostalCode int            `json:"postal_code" validate:"required"`
	Total      int            `json:"total" validate:"required"`
	Status     string         `json:"status" validate:"required"`
	Order      []OrderRequest `json:"orders" validate:"required"`
	User_Id    string
}
