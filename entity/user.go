package entity

type User struct {
	Id       int    `db:"id" json:"id"`
	Name     string `db:"name" json:"name"`
	Email    string `db:"email" json:"email"`
	Password string `db:"password" json:"-"`
	Gender   string `db:"gender" json:"gender"`
	Phone    string `db:"phone" json:"phone"`
	Image    string `db:"image" json:"image"`
	IsAdmin  bool   `db:"is_admin" json:"is_admin"`
}

type Address struct {
	Id         int    `db:"id" json:"id"`
	UserId     int    `db:"user_id" json:"-"`
	Name       string `db:"name" json:"name"`
	Phone      string `db:"phone" json:"phone"`
	Address    string `db:"address" json:"address"`
	City       string `db:"city" json:"city"`
	PostalCode uint16 `db:"postal_code" json:"postal_code"`
}

type AddressRequest struct {
	Name       string `json:"name" validate:"required"`
	Phone      string `json:"phone" validate:"required"`
	Address    string `json:"address" validate:"required"`
	City       string `json:"city" validate:"required"`
	PostalCode uint16 `json:"postal_code" validate:"required"`
}
