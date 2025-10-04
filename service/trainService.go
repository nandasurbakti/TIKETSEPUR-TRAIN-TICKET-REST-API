package service

import (
	"errors"
	"tiketsepur/dto"
	"tiketsepur/models"
	"tiketsepur/repository"
)

type TrainService interface {
	Create(req dto.CreateTrainRequest) (*models.Train, error)
	GetByID(id int) (*models.Train, error)
	GetAll() ([]models.Train, error)
	Update(id int, req dto.UpdateTrainRequest) (*models.Train, error)
	Delete(id int) error
}

type trainService struct {
	trainRepo repository.TrainRepository
}

func NewTrainService(trainRepo repository.TrainRepository) TrainService {
	return &trainService{trainRepo: trainRepo}
}

func (s *trainService) Create(req dto.CreateTrainRequest) (*models.Train, error) {
	train := &models.Train{
		TrainCode:  req.TrainCode,
		TrainName:  req.TrainName,
		TrainType:  req.TrainType,
		TotalSeats: req.TotalSeats,
	}

	if err := s.trainRepo.Create(train); err != nil {
		return nil, err
	}

	return train, nil
}

func (s *trainService) GetByID(id int) (*models.Train, error) {
	return s.trainRepo.FindByID(id)
}

func (s *trainService) GetAll() ([]models.Train, error) {
	return s.trainRepo.FindAll()
}

func (s *trainService) Update(id int, req dto.UpdateTrainRequest) (*models.Train, error) {
	train, err := s.trainRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("kereta tidak ditemukan")
	}

	if req.TrainCode != nil {
		train.TrainCode = *req.TrainCode
	}

	if req.TrainName != nil {
		train.TrainName = *req.TrainName
	}

	if req.TrainType != nil {
		train.TrainType = *req.TrainType
	}

	if req.TotalSeats != nil {
		train.TotalSeats = *req.TotalSeats
	}

	if err := s.trainRepo.Update(id, train); err != nil {
		return nil, err
	}

	return train, nil
}

func (s *trainService) Delete(id int) error {
	_, err := s.trainRepo.FindByID(id)
	if err != nil {
		return errors.New("kereta tidak ditemukan")
	}

	return s.trainRepo.Delete(id)
}