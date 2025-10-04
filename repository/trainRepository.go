package repository

import (
	"tiketsepur/models"

	"github.com/jmoiron/sqlx"
)

type TrainRepository interface {
	Create(train *models.Train) error
	FindByID(id int) (*models.Train, error)
	FindAll() ([]models.Train, error)
	Update(id int, train *models.Train) error
	Delete(id int) error
}

type trainRepository struct {
	db *sqlx.DB
}

func NewTrainRepository(db *sqlx.DB) TrainRepository {
	return &trainRepository{db: db}
}

func (r *trainRepository) Create(train *models.Train) error {
	query := `INSERT INTO trains (train_code, train_name, train_type, total_seats, created_at) 
			  VALUES ($1, $2, $3, $4, NOW()) RETURNING id`
	return r.db.QueryRow(query, train.TrainCode, train.TrainName, train.TrainType, train.TotalSeats).Scan(&train.ID)
}

func (r *trainRepository) FindByID(id int) (*models.Train, error) {
	var train models.Train
	query := `SELECT * FROM trains WHERE id = $1`
	err := r.db.Get(&train, query, id)
	if err != nil {
		return nil, err
	}
	return &train, nil
}

func (r *trainRepository) FindAll() ([]models.Train, error) {
	var trains []models.Train
	query := `SELECT * FROM trains ORDER BY created_at DESC`
	err := r.db.Select(&trains, query)
	return trains, err
}

func (r *trainRepository) Update(id int, train *models.Train) error {
	query := `UPDATE trains SET train_code = $1, train_name = $2, train_type = $3, 
			  total_seats = $4, modified_at = NOW() WHERE id = $5`
	_, err := r.db.Exec(query, train.TrainCode, train.TrainName, train.TrainType, train.TotalSeats, id)
	return err
}

func (r *trainRepository) Delete(id int) error {
	query := `DELETE FROM trains WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}