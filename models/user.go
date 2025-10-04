package models

type User struct {
	ID         int    `json:"id" db:"id"`
	Email      string `json:"email" db:"email"`
	Password   string `json:"-" db:"password"`
	FullName   string `json:"full_name" db:"full_name"`
	Phone      string `db:"phone" json:"phone"`
	Role       string `db:"role" json:"role"`
	CreatedAt  string `json:"created_at" db:"created_at"`
	ModifiedAt string `json:"modified_at" db:"modified_at"`
}