package models

import "time"

type Train struct {
	ID         int       `json:"id" db:"id"`
	TrainCode  string    `json:"train_code" db:"train_code"`
	TrainName  string    `json:"train_name" db:"train_name"`
	TrainType  string    `json:"train_type" db:"train_type"`
	TotalSeats int       `json:"total_seats" db:"total_seats"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	ModifiedAt time.Time `json:"modified_at" db:"modified_at"`
}
