package dto

import "time"

type CreateScheduleRequest struct {
	TrainID          int       `json:"train_id" binding:"required"`
	DepartureStation string    `json:"departure_station" binding:"required"`
	ArrivalStation   string    `json:"arrival_station" binding:"required"`
	DepartureTime    time.Time `json:"departure_time" binding:"required"`
	ArrivalTime      time.Time `json:"arrival_time" binding:"required"`
	Price            float64   `json:"price" binding:"required,min=0"`
	AvailableSeats   int       `json:"available_seats" binding:"required,min=0"`
}

type UpdateScheduleRequest struct {
	TrainID          *int       `json:"train_id"`
	DepartureStation *string    `json:"departure_station"`
	ArrivalStation   *string    `json:"arrival_station"`
	DepartureTime    *time.Time `json:"departure_time"`
	ArrivalTime      *time.Time `json:"arrival_time"`
	Price            *float64   `json:"price" binding:"omitempty,min=0"`
	AvailableSeats   *int       `json:"available_seats" binding:"omitempty,min=0"`
}

type SearchScheduleRequest struct {
	DepartureStation string `form:"departure_station"`
	ArrivalStation   string `form:"arrival_station"`
	Date             string `form:"date"`
}