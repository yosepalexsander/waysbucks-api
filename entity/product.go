package entity

type Product struct {
	Id int `db:"id" json:"id"`
	Name string `db:"name" json:"name" validate:"required"`
	Description string `db:"description" json:"description" validate:"required"`
	Image string `db:"image" json:"image" validate:"required"`
	Price int `db:"price" json:"price" validate:"required"`
	Is_Available bool `db:"is_available" json:"is_available"`
}

