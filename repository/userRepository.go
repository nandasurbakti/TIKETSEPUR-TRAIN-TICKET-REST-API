package repository

import (
	"tiketsepur/models"

	"github.com/jmoiron/sqlx"
)

type UserRepository interface {
	Create(user *models.User) error
	FindByID(id int) (*models.User, error)
	FindByEmail(email string) (*models.User, error)
	FindAll() ([]models.User, error)
	Update(id int, user *models.User) error
	Delete(id int) error
}

type userRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *models.User) error {
	query := `INSERT INTO users (email, password, full_name, phone, role, created_at) 
			  VALUES ($1, $2, $3, $4, $5, NOW()) RETURNING id`
	return r.db.QueryRow(query, user.Email, user.Password, user.FullName, user.Phone, user.Role).Scan(&user.ID)
}

func (r *userRepository) FindByID(id int) (*models.User, error) {
	var user models.User
	query := `SELECT * FROM users WHERE id = $1`
	err := r.db.Get(&user, query, id)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	query := `SELECT * FROM users WHERE email = $1`
	err := r.db.Get(&user, query, email)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindAll() ([]models.User, error) {
	var users []models.User
	query := `SELECT * FROM users ORDER BY created_at DESC`
	err := r.db.Select(&users, query)
	return users, err
}

func (r *userRepository) Update(id int, user *models.User) error {
	query := `UPDATE users SET email = $1, full_name = $2, phone = $3, role = $4, modified_at = NOW() WHERE id = $5`
	_, err := r.db.Exec(query, user.Email, user.FullName, user.Phone, user.Role, id)
	return err
}

func (r *userRepository) Delete(id int) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}