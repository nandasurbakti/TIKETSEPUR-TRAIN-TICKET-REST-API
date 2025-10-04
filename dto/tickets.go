package dto

type CreateTicketRequest struct {
	ScheduleID        int    `json:"schedule_id" binding:"required"`
	SeatNumber        string `json:"seat_number" binding:"required"`
	PassengerName     string `json:"passenger_name" binding:"required"`
	PassengerIDNumber string `json:"passenger_id_number" binding:"required"`
	PaymentMethod     string `json:"payment_method" binding:"required,oneof=bank_transfer e-wallet credit_card"`
}

type CancelTicketRequest struct {
	Reason string `json:"reason"`
}

type TicketResponse struct {
	ID                int     `json:"id"`
	BookingCode       string  `json:"booking_code"`
	ScheduleID        int     `json:"schedule_id"`
	SeatNumber        string  `json:"seat_number"`
	PassengerName     string  `json:"passenger_name"`
	PassengerIDNumber string  `json:"passenger_id_number"`
	Status            string  `json:"status"`
	TotalPrice        float64 `json:"total_price"`
	DepartureStation  string  `json:"departure_station"`
	ArrivalStation    string  `json:"arrival_station"`
	DepartureTime     string  `json:"departure_time"`
	TrainName         string  `json:"train_name"`
	PaymentStatus     string  `json:"payment_status"`
	PaymentCode       string  `json:"payment_code"`
}