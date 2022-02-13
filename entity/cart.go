package entity

type Cart struct {
	Id         int           `db:"id" json:"id"`
	Price      int           `db:"price" json:"price"`
	Qty        int           `db:"qty" json:"qty"`
	Product_Id int           `db:"product_id" json:"product_id,omitempty"`
	ToppingIds []int64       `db:"topping_id" json:"topping_id,omitempty"`
	User_Id    string        `db:"user_id" json:"-"`
	Product    CartProduct   `json:"product"`
	Topping    []CartTopping `json:"toppings"`
}

type CartProduct struct {
	Id    int    `db:"id" json:"id"`
	Name  string `db:"name" json:"name"`
	Image string `db:"image" json:"image"`
	Price int    `db:"price" json:"price"`
}

type CartTopping struct {
	Id   int    `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
}
