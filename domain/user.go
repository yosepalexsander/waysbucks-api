package domain

type User struct {
	Id uint64 `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
	Email string `db:"email" json:"email"`
	Password string `db:"password" json:"password"`
	Gender string `db:"gender" json:"gender"`
	Phone string `db:"phone" json:"phone"`
	Image *string `db:"image" json:"image"`
	IsAdmin uint8 `db:"is_admin" json:"is_admin"`
}