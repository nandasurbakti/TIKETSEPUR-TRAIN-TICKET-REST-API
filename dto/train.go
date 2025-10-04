package dto

type CreateTrainRequest struct {
	TrainCode  string `json:"train_code" binding:"required"`
	TrainName  string `json:"train_name" binding:"required"`
	TrainType  string `json:"train_type" binding:"required"`
	TotalSeats int    `json:"total_seats" binding:"required,min=1"`
}

type UpdateTrainRequest struct {
	TrainCode  *string `json:"train_code"`
	TrainName  *string `json:"train_name"`
	TrainType  *string `json:"train_type"`
	TotalSeats *int    `json:"total_seats" binding:"omitempty,min=1"`
}