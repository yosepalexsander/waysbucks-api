package entity

type Transaction struct {
	Id int `db:"id"`
	User_Id int `db:"user_id"`
	Name string `db:"name" json:"name"`
	Address string `db:"address" json:"address,omitempty"`
	PostCode int `db:"postcode" json:"postcode,omitempty"`
	Phone string `db:"phone" json:"phone,omitempty"`
	Total int `db:"total" json:"total"` 
	Status string `db:"status" json:"status"`
	Orders []Order
}

type Order struct {
	Id int `db:"id" json:"id"`
	Transaction_Id int `db:"transaction_id"`
	Product_Id int `db:"product_id"`
	Topping_Ids []int `db:"topping_id"`
	OrderProduct 
	Toppings []OrderTopping `json:"toppings"`
	Price int `db:"price" json:"price"`
	Qty int `db:"qty" json:"qty"`
}

type TransactionTxParams struct {
	Transaction Transaction
	Order []Order 
	ProductIds []int
}

type OrderProduct struct {
	Name string `db:"name" json:"name"`
	Image string `db:"image" json:"image"`
}

type OrderTopping struct {
	Id int `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
}