package entity

type Cart struct {
	Id int `db:"id" json:"id"`
	User_Id int `db:"user_id" json:"user_id,omitempty"`
	Product_Id int `db:"product_id" json:"product_id,omitempty"`
	ToppingIds []int64 `db:"topping_id" json:"topping_id,omitempty"` 
	Price int `db:"price" json:"price"`
	Qty int `db:"qty" json:"qty"`
	Product CartProduct `json:"product"`
	Topping []CartTopping `json:"toppings"`
}

type CartProduct struct {
	Id int `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
	Image string `db:"image" json:"image"`
	Price int `db:"price" json:"price"`
}

type CartTopping struct {
	Id int `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
	Price int `db:"price" json:"price"`
}
