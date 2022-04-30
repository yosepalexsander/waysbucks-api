package entity

import "time"

type Product struct {
	Id          int       `db:"id" json:"id"`
	Name        string    `db:"name" json:"name"`
	Description string    `db:"description" json:"description"`
	Image       string    `db:"image" json:"image"`
	Price       int       `db:"price" json:"price"`
	IsAvailable bool      `db:"is_available" json:"is_available"`
	Created_At  time.Time `db:"created_at" json:"created_at"`
	Updated_At  time.Time `db:"updated_at" json:"updated_at"`
}

type ProductRequest struct {
	Name        string `schema:"name,required"`
	Description string `schema:"description,required"`
	Image       string `schema:"image"`
	Price       int    `schema:"price,required"`
	IsAvailable bool   `schema:"is_available"`
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

func NewProduct(req ProductRequest) Product {
	return Product{
		Name:        req.Name,
		Description: req.Description,
		Image:       req.Image,
		Price:       req.Price,
		IsAvailable: req.IsAvailable,
	}
}

func NewProductTopping(req ProductToppingRequest) ProductTopping {
	return ProductTopping{
		Name:        req.Name,
		Image:       req.Image,
		Price:       req.Price,
		IsAvailable: req.IsAvailable,
	}
}
