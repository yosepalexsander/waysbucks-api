package entity

import "time"

type Product struct {
	Id          string    `db:"id" json:"id"`
	Name        string    `db:"name" json:"name" validate:"required"`
	Description string    `db:"description" json:"description" validate:"required"`
	Image       string    `db:"image" json:"image" validate:"required"`
	Price       int       `db:"price" json:"price" validate:"required"`
	IsAvailable bool      `db:"is_available" json:"is_available"`
	Created_At  time.Time `db:"created_at" json:"created_at"`
	Updated_At  time.Time `db:"updated_at" json:"updated_at"`
}

type ProductRequest struct {
	Id          int    `json:"id"`
	Name        string `json:"name" schema:"name,required"`
	Description string `json:"description" schema:"description,required"`
	Image       string `json:"image" schema:"image"`
	Price       int    `json:"price" schema:"price,required"`
	IsAvailable bool   `json:"is_available" schema:"is_available"`
}

type ProductTopping struct {
	Id          int    `db:"id" json:"id"`
	Name        string `db:"name" json:"name"`
	Image       string `db:"image" json:"image"`
	Price       int    `db:"price" json:"price"`
	IsAvailable bool   `db:"is_available" json:"is_available"`
}

type ProductToppingRequest struct {
	Name        string `json:"name" validate:"required"`
	Image       string `json:"image" validate:"required"`
	Price       int    `json:"price" validate:"required"`
	IsAvailable bool   `json:"is_available"`
}
