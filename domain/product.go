package domain

type Product struct {
	Name string `json:"name"`
	Price int16 `json:"price"`
	Image string `json:"image"`
}