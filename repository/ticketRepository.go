package repository

import (
	"tiketsepur/models"

	"github.com/jmoiron/sqlx"
)

type TicketRepository interface {
	Create(ticket *models.Ticket, tx *sqlx.Tx) error
	FindByID(id int) (*models.TicketWithDetails, error)
	FindByBookingCode(code string) (*models.TicketWithDetails, error)
	FindByUserID(userID int) ([]models.TicketWithDetails, error)
	FindAll() ([]models.TicketWithDetails, error)
	UpdateStatus(id int, status string, tx *sqlx.Tx) error
	CheckSeatAvailability(scheduleID int, seatNumber string) (bool, error)
}

type ticketRepository struct {
	db *sqlx.DB
}

func NewTicketRepository(db *sqlx.DB) TicketRepository {
	return &ticketRepository{db: db}
}

func (r *ticketRepository) Create(ticket *models.Ticket, tx *sqlx.Tx) error {
	query := `INSERT INTO tickets (user_id, schedule_id, seat_number, passenger_name, 
			  passenger_id_number, status, booking_code, total_price, created_at, modified_at) 
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW(), NOW()) RETURNING id, created_at, modified_at`
	return tx.QueryRow(query, ticket.UserID, ticket.ScheduleID, ticket.SeatNumber,
		ticket.PassengerName, ticket.PassengerIDNumber, ticket.Status,
		ticket.BookingCode, ticket.TotalPrice).Scan(&ticket.ID, &ticket.CreatedAt, &ticket.ModifiedAt)
}

func (r *ticketRepository) FindByID(id int) (*models.TicketWithDetails, error) {
	var ticket models.TicketWithDetails
	query := `SELECT t.id, t.user_id, t.schedule_id, t.seat_number, 
    		  t.passenger_name, t.passenger_id_number, t.status, 
    		  t.booking_code, t.total_price, t.created_at, t.modified_at, s.departure_station, s.arrival_station, s.departure_time,
			  tr.train_name, tr.train_code, tr.train_type, p.payment_code, p.payment_status FROM tickets t
			  JOIN schedules s ON t.schedule_id = s.id
			  JOIN trains tr ON s.train_id = tr.id
			  LEFT JOIN payments p ON p.ticket_id = t.id
			  WHERE t.id = $1`
	err := r.db.Get(&ticket, query, id)
	if err != nil {
		return nil, err
	}
	return &ticket, nil
}

func (r *ticketRepository) FindByBookingCode(code string) (*models.TicketWithDetails, error) {
	var ticket models.TicketWithDetails
	query := `SELECT t.id, t.user_id, t.schedule_id, t.seat_number, 
    		  t.passenger_name, t.passenger_id_number, t.status, 
    		  t.booking_code, t.total_price, t.created_at, t.modified_at, s.departure_station, s.arrival_station, s.departure_time,
			  tr.train_name, tr.train_code, tr.train_type, p.payment_code, p.payment_status FROM tickets t
			  JOIN schedules s ON t.schedule_id = s.id
			  JOIN trains tr ON s.train_id = tr.id
			  LEFT JOIN payments p ON p.ticket_id = t.id
			  WHERE t.booking_code = $1`
	err := r.db.Get(&ticket, query, code)
	if err != nil {
		return nil, err
	}
	return &ticket, nil
}

func (r *ticketRepository) FindByUserID(userID int) ([]models.TicketWithDetails, error) {
	var tickets []models.TicketWithDetails
	query := `SELECT t.id, t.user_id, t.schedule_id, t.seat_number, 
			t.passenger_name, t.passenger_id_number, t.status, 
			t.booking_code, t.total_price, t.created_at, t.modified_at,
			s.departure_station AS departure_station,
			s.arrival_station AS arrival_station,
			s.departure_time AS departure_time,
			tr.train_name AS train_name,
			tr.train_code AS train_code,
			tr.train_type AS train_type,
			p.payment_code AS payment_code,
			p.payment_status AS payment_status
			FROM tickets t
			JOIN schedules s ON t.schedule_id = s.id
			JOIN trains tr ON s.train_id = tr.id
			LEFT JOIN payments p ON p.ticket_id = t.id
			WHERE t.user_id = $1
			ORDER BY t.created_at DESC`
	err := r.db.Select(&tickets, query, userID)
	return tickets, err
}

func (r *ticketRepository) FindAll() ([]models.TicketWithDetails, error) {
	var tickets []models.TicketWithDetails
	query := `SELECT t.id, t.user_id, t.schedule_id, t.seat_number, 
    		  t.passenger_name, t.passenger_id_number, t.status, 
    		  t.booking_code, t.total_price, t.created_at, t.modified_at, s.departure_station, s.arrival_station, s.departure_time,
			  tr.train_name, tr.train_code, tr.train_type,p.payment_code, p.payment_status FROM tickets t
			  JOIN schedules s ON t.schedule_id = s.id
			  JOIN trains tr ON s.train_id = tr.id
			  LEFT JOIN payments p ON p.ticket_id = t.id
			  ORDER BY t.created_at DESC`
	err := r.db.Select(&tickets, query)
	return tickets, err
}

func (r *ticketRepository) UpdateStatus(id int, status string, tx *sqlx.Tx) error {
	query := `UPDATE tickets SET status = $1, modified_at = NOW() WHERE id = $2`
	_, err := tx.Exec(query, status, id)
	return err
}

func (r *ticketRepository) CheckSeatAvailability(scheduleID int, seatNumber string) (bool, error) {
	var count int
	query := `SELECT COUNT(*) FROM tickets WHERE schedule_id = $1 AND seat_number = $2 
			  AND status != 'cancelled'`
	err := r.db.Get(&count, query, scheduleID, seatNumber)
	return count == 0, err
}