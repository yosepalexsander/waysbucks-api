package entity

type User struct {
	Id int `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
	Email string `db:"email" json:"email"`
	Password string `db:"password" json:"password,omitempty"`
	Gender string `db:"gender" json:"gender"`
	Phone string `db:"phone" json:"phone"`
	Image string `db:"image" json:"image,omitempty"`
	IsAdmin bool `db:"is_admin" json:"is_admin"`
}

type UserAddress struct {
	Id int `db:"id" json:"id"`
	UserId int `db:"user_id" json:"user_id,omitempty"`
	Name string `db:"name" json:"name" validate:"required"`
	Phone string `db:"phone" json:"phone" validate:"required"`
	Address string `db:"address" json:"address" validate:"required"`
	City string `db:"city" json:"city" validate:"required"`
	PostalCode uint16 `db:"postal_code" json:"postal_code" validate:"required"`
}