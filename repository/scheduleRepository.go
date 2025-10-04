package repository

import (
	"errors"
	"tiketsepur/models"

	"github.com/jmoiron/sqlx"
)

type ScheduleRepository interface {
	Create(schedule *models.Schedule) error
	FindByID(id int) (*models.Schedule, error)
	FindAll() ([]models.Schedule, error)
	Search(departure, arrival, date string) ([]models.Schedule, error)
	Update(id int, schedule *models.Schedule) error
	Delete(id int) error
	DecrementSeat(id int, tx *sqlx.Tx) error
	IncrementSeat(id int, tx *sqlx.Tx) error
}

type scheduleRepository struct {
	db *sqlx.DB
}

func NewScheduleRepository(db *sqlx.DB) ScheduleRepository {
	return &scheduleRepository{db: db}
}

func (r *scheduleRepository) Create(schedule *models.Schedule) error {
	query := `INSERT INTO schedules (train_id, departure_station, arrival_station, 
			  departure_time, arrival_time, price, available_seats, created_at) 
			  VALUES ($1, $2, $3, $4, $5, $6, $7, NOW()) RETURNING id`
	return r.db.QueryRow(query, schedule.TrainID, schedule.DepartureStation,
		schedule.ArrivalStation, schedule.DepartureTime, schedule.ArrivalTime,
		schedule.Price, schedule.AvailableSeats).Scan(&schedule.ID)
}

func (r *scheduleRepository) FindByID(id int) (*models.Schedule, error) {
	var schedule models.Schedule
	query := `SELECT s.id, s.train_id, s.departure_station, s.arrival_station, 
				s.departure_time, s.arrival_time, s.price, s.available_seats, 
				s.created_at, s.modified_at,
				t.id as train_id, t.train_code, t.train_name, t.train_type
				FROM schedules s
				JOIN trains t ON s.train_id = t.id
				WHERE s.id = $1`
	err := r.db.Get(&schedule, query, id)
	if err != nil {
		return nil, err
	}
	return &schedule, nil
}

func (r *scheduleRepository) FindAll() ([]models.Schedule, error) {
	var schedules []models.Schedule
	query := `SELECT s.id, s.train_id, s.departure_station, s.arrival_station, 
			s.departure_time, s.arrival_time, s.price, s.available_seats, 
			s.created_at, s.modified_at,
			t.id as train_id, t.train_code, t.train_name, t.train_type
			FROM schedules s
			JOIN trains t ON s.train_id = t.id
			ORDER BY s.departure_time ASC`
	err := r.db.Select(&schedules, query)
	return schedules, err
}

func (r *scheduleRepository) Search(departure, arrival, date string) ([]models.Schedule, error) {
	var schedules []models.Schedule
	var err error

	baseQuery := `SELECT s.id, s.train_id, s.departure_station, s.arrival_station, 
			s.departure_time, s.arrival_time, s.price, s.available_seats, 
			s.created_at, s.modified_at,
			t.id as train_id, t.train_code, t.train_name, t.train_type
			FROM schedules s
			JOIN trains t ON s.train_id = t.id
			WHERE s.departure_station ILIKE $1 
			AND s.arrival_station ILIKE $2 
			AND s.available_seats > 0`

	if date != "" {
		baseQuery += " AND DATE(s.departure_time) = $3 ORDER BY s.departure_time ASC"
		err = r.db.Select(&schedules, baseQuery, "%"+departure+"%", "%"+arrival+"%", date)
	} else {
		baseQuery += " ORDER BY s.departure_time ASC"
		err = r.db.Select(&schedules, baseQuery, "%"+departure+"%", "%"+arrival+"%")
	}

	return schedules, err
}


func (r *scheduleRepository) Update(id int, schedule *models.Schedule) error {
	query := `UPDATE schedules SET train_id = $1, departure_station = $2, 
			  arrival_station = $3, departure_time = $4, arrival_time = $5, 
			  price = $6, available_seats = $7, modified_at = NOW() WHERE id = $8`
	_, err := r.db.Exec(query, schedule.TrainID, schedule.DepartureStation,
		schedule.ArrivalStation, schedule.DepartureTime, schedule.ArrivalTime,
		schedule.Price, schedule.AvailableSeats, id)
	return err
}

func (r *scheduleRepository) Delete(id int) error {
	query := `DELETE FROM schedules WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *scheduleRepository) DecrementSeat(id int, tx *sqlx.Tx) error {
	query := `UPDATE schedules SET available_seats = available_seats - 1, 
			  modified_at = NOW() WHERE id = $1 AND available_seats > 0`
	result, err := tx.Exec(query, id)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return ErrNoSeatsAvailable
	}
	return nil
}

func (r *scheduleRepository) IncrementSeat(id int, tx *sqlx.Tx) error {
	query := `UPDATE schedules SET available_seats = available_seats + 1, 
			  modified_at = NOW() WHERE id = $1`
	_, err := tx.Exec(query, id)
	return err
}

var ErrNoSeatsAvailable = errors.New("kursi tidak tersedia")