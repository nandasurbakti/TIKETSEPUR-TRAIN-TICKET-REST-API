package models

import "time"

type Payment struct {
	ID            int        `json:"id" db:"id"`
	TicketID      int        `json:"ticket_id" db:"ticket_id"`
	PaymentMethod string     `json:"payment_method" db:"payment_method"`
	PaymentAmount float64    `json:"payment_amount" db:"payment_amount"`
	PaymentStatus string     `json:"payment_status" db:"payment_status"`
	PaymentCode   string     `json:"payment_code" db:"payment_code"`
	PaidAt        *time.Time `json:"paid_at" db:"paid_at"`
	CreatedAt     time.Time  `json:"created_at" db:"created_at"`
	ModifiedAt    time.Time  `json:"modified_at" db:"modified_at"`
}