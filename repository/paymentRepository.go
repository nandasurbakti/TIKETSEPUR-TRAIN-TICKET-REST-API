package repository

import (
	"tiketsepur/models"

	"github.com/jmoiron/sqlx"
)

type PaymentRepository interface {
	Create(payment *models.Payment, tx *sqlx.Tx) error
	FindByID(id int) (*models.Payment, error)
	FindByTicketID(ticketID int) (*models.Payment, error)
	FindByPaymentCode(code string) (*models.Payment, error)
	UpdateStatus(id int, status string) error
	UpdateStatusTx(id int, status string, tx *sqlx.Tx) error
	UpdateStatusTxByTicketID(ticketID int, status string, tx *sqlx.Tx) error
}

type paymentRepository struct {
	db *sqlx.DB
}

func NewPaymentRepository(db *sqlx.DB) PaymentRepository {
	return &paymentRepository{db: db}
}

func (r *paymentRepository) Create(payment *models.Payment, tx *sqlx.Tx) error {
	query := `INSERT INTO payments (ticket_id, payment_method, payment_amount, 
			  payment_status, payment_code) 
			  VALUES ($1, $2, $3, $4, $5) RETURNING id`
	return tx.QueryRow(query, payment.TicketID, payment.PaymentMethod,
		payment.PaymentAmount, payment.PaymentStatus, payment.PaymentCode).Scan(&payment.ID)
}

func (r *paymentRepository) FindByID(id int) (*models.Payment, error) {
	var payment models.Payment
	query := `SELECT * FROM payments WHERE id = $1`
	err := r.db.Get(&payment, query, id)
	if err != nil {
		return nil, err
	}
	return &payment, nil
}

func (r *paymentRepository) FindByTicketID(ticketID int) (*models.Payment, error) {
	var payment models.Payment
	query := `SELECT * FROM payments WHERE ticket_id = $1`
	err := r.db.Get(&payment, query, ticketID)
	if err != nil {
		return nil, err
	}
	return &payment, nil
}

func (r *paymentRepository) FindByPaymentCode(code string) (*models.Payment, error) {
	var payment models.Payment
	query := `SELECT * FROM payments WHERE payment_code = $1`
	err := r.db.Get(&payment, query, code)
	if err != nil {
		return nil, err
	}
	return &payment, nil
}

func (r *paymentRepository) UpdateStatus(id int, status string) error {
	query := `UPDATE payments SET payment_status = $1, paid_at = NOW(), 
			  modified_at = NOW() WHERE id = $2`
	_, err := r.db.Exec(query, status, id)
	return err
}

func (r *paymentRepository) UpdateStatusTx(id int, status string, tx *sqlx.Tx) error {
    query := `UPDATE payments 
              SET payment_status = $1, paid_at = NOW(), modified_at = NOW() 
              WHERE id = $2`
    _, err := tx.Exec(query, status, id)
    return err
}

func (r *paymentRepository) UpdateStatusTxByTicketID(ticketID int, status string, tx *sqlx.Tx) error {
	query := `UPDATE payments SET payment_status = $1, modified_at = NOW() WHERE ticket_id = $2`
	_, err := tx.Exec(query, status, ticketID)
	return err
}