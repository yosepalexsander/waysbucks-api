package entity

type Cart struct {
	Id         int           `db:"id" json:"id"`
	Price      int           `db:"price" json:"price"`
	Qty        int           `db:"qty" json:"qty"`
	ProductId  int           `db:"product_id" json:"-"`
	ToppingIds []int64       `db:"topping_id" json:"-"`
	UserId     string        `db:"user_id" json:"-"`
	Product    CartProduct   `json:"product"`
	Topping    []CartTopping `json:"toppings"`
}

type CartRequest struct {
	Price      int     `json:"price" validate:"required"`
	Qty        int     `json:"qty" validate:"required"`
	ProductId  int     `json:"product_id" validate:"required"`
	ToppingIds []int64 `json:"topping_id"`
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

func NewCart(productID int, price int, qty int, toppingID []int64, userID string) Cart {
	return Cart{
		Price:      price,
		Qty:        qty,
		ProductId:  productID,
		ToppingIds: toppingID,
		UserId:     userID,
	}
}
