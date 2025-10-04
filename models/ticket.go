package models

import "time"

type Ticket struct {
	ID                int       `json:"id" db:"id"`
	UserID            int       `json:"user_id" db:"user_id"`
	ScheduleID        int       `json:"schedule_id" db:"schedule_id"`
	SeatNumber        string    `json:"seat_number" db:"seat_number"`
	PassengerName     string    `json:"passenger_name" db:"passenger_name"`
	PassengerIDNumber string    `json:"passenger_id_number" db:"passenger_id_number"`
	Status            string    `json:"status" db:"status"`
	BookingCode       string    `json:"booking_code" db:"booking_code"`
	TotalPrice        float64   `json:"total_price" db:"total_price"`
	CreatedAt         time.Time `json:"created_at" db:"created_at"`
	ModifiedAt        *time.Time `json:"modified_at" db:"modified_at"`

	User *User `json:"user,omitempty" db:"-"`
	Schedule *Schedule `json:"schedule,omitempty" db:"-"`
	Payment *Payment `json:"payment,omitempty" db:"-"`
}

type TicketWithDetails struct {
	ID                int        `json:"id" db:"id"`
	UserID            int        `json:"user_id" db:"user_id"`
	ScheduleID        int        `json:"schedule_id" db:"schedule_id"`
	SeatNumber        string     `json:"seat_number" db:"seat_number"`
	PassengerName     string     `json:"passenger_name" db:"passenger_name"`
	PassengerIDNumber string     `json:"passenger_id_number" db:"passenger_id_number"`
	Status            string     `json:"status" db:"status"`
	BookingCode       string     `json:"booking_code" db:"booking_code"`
	TotalPrice        float64    `json:"total_price" db:"total_price"`
	CreatedAt         time.Time  `json:"created_at" db:"created_at"`
	ModifiedAt        *time.Time `json:"modified_at" db:"modified_at"`
	DepartureStation  string     `db:"departure_station"`
	ArrivalStation    string     `db:"arrival_station"`
	DepartureTime     time.Time  `db:"departure_time"`
	TrainName         string     `db:"train_name"`
	TrainCode         string     `db:"train_code"`
	TrainType         string     `db:"train_type"`
	PaymentCode       *string    `db:"payment_code"`
	PaymentStatus     *string    `db:"payment_status"`

	User     *User     `json:"user,omitempty" db:"-"`
	Schedule *Schedule `json:"schedule,omitempty" db:"-"`
	Payment  *Payment  `json:"payment,omitempty" db:"-"`
}