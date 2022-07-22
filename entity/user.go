package entity

type User struct {
	Id       string `db:"id" json:"id"`
	Name     string `db:"name" json:"name"`
	Email    string `db:"email" json:"email"`
	Password string `db:"password" json:"-"`
	Gender   string `db:"gender" json:"gender"`
	Phone    string `db:"phone" json:"phone"`
	Image    string `db:"image" json:"image"`
	IsAdmin  bool   `db:"is_admin" json:"is_admin"`
}

type Address struct {
	Id         string  `db:"id" json:"id"`
	Name       string  `db:"name" json:"name"`
	Phone      string  `db:"phone" json:"phone"`
	Address    string  `db:"address" json:"address"`
	City       string  `db:"city" json:"city"`
	PostalCode uint16  `db:"postal_code" json:"postal_code"`
	Longitude  float64 `db:"longitude" json:"longitude"`
	Latitude   float64 `db:"latitude" json:"latitude"`
	UserId     string  `db:"user_id" json:"-"`
}

type AddressRequest struct {
	Name       string  `json:"name" validate:"required"`
	Phone      string  `json:"phone" validate:"required"`
	Address    string  `json:"address" validate:"required"`
	City       string  `json:"city" validate:"required"`
	PostalCode uint16  `json:"postal_code" validate:"required"`
	Longitude  float64 `json:"longitude" validate:"required"`
	Latitude   float64 `json:"latitude" validate:"required"`
}

func NewUser(name string, email string, password string, gender string, phone string) User {
	return User{
		Name:     name,
		Email:    email,
		Password: password,
		Gender:   gender,
		Phone:    phone,
	}
}
