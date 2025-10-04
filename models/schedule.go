package models

import "time"

type Schedule struct {
	ID               int       `json:"id" db:"id"`
	TrainID          int       `json:"train_id" db:"train_id"`
	DepartureStation string    `json:"departure_station" db:"departure_station"`
	ArrivalStation   string    `json:"arrival_station" db:"arrival_station"`
	DepartureTime    time.Time `json:"departure_time" db:"departure_time"`
	ArrivalTime 	 time.Time `json:"arrival_time" db:"arrival_time"`
	Price 			 float64   `json:"price" db:"price"`
	AvailableSeats 	 int 	   `json:"available_seats" db:"available_seats"`
	TrainCode        string    `db:"train_code"`
    TrainName        string    `db:"train_name"`
    TrainType        string    `db:"train_type"`
	CreatedAt  	 	 time.Time `json:"created_at" db:"created_at"`
	ModifiedAt 		 time.Time `json:"modified_at" db:"modified_at"`
}