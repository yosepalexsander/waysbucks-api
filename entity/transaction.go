package entity

type Transaction struct {
	Id         int    `db:"id"`
	User_Id    int    `db:"user_id"`
	Name       string `db:"name"`
	Address    string `db:"address"`
	PostalCode int    `db:"postal_code"`
	City       string `db:"city"`
	Phone      string `db:"phone"`
	Total      int    `db:"total"`
	Status     string `db:"status"`
	Orders     []Order
}

type Order struct {
	Id             int     `db:"id" json:"id"`
	Transaction_Id int     `db:"transaction_id"`
	Product_Id     int     `db:"product_id"`
	Topping_Ids    []int64 `db:"topping_id"`
	OrderProduct
	Price    int `db:"price"`
	Qty      int `db:"qty"`
	Toppings []OrderTopping
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
	City       string         `json:"city" validate:"required"`
	Phone      string         `json:"phone" validate:"required"`
	Total      int            `json:"total" validate:"required"`
	Status     string         `json:"status" validate:"required"`
	Order      []OrderRequest `json:"orders"`
}

type TransactionResponse struct {
	Name       string          `json:"name"`
	Address    string          `json:"address"`
	PostalCode int             `json:"postal_code"`
	City       string          `json:"city"`
	Total      int             `json:"total"`
	Status     string          `json:"status"`
	Orders     []OrderResponse `json:"orders"`
}

type OrderResponse struct {
	Id int `json:"id"`
	OrderProduct
	Price    int            `json:"price"`
	Qty      int            `json:"qty"`
	Toppings []OrderTopping `json:"toppings"`
}
