package service

import (
	"errors"
	"tiketsepur/dto"
	"tiketsepur/models"
	"tiketsepur/repository"
)

type ScheduleService interface {
	Create(req dto.CreateScheduleRequest) (*models.Schedule, error)
	GetByID(id int) (*models.Schedule, error)
	GetAll() ([]models.Schedule, error)
	Search(req dto.SearchScheduleRequest) ([]models.Schedule, error)
	Update(id int, req dto.UpdateScheduleRequest) (*models.Schedule, error)
	Delete(id int) error
}

type scheduleService struct {
	scheduleRepo repository.ScheduleRepository
	trainRepo    repository.TrainRepository
}

func NewScheduleService(scheduleRepo repository.ScheduleRepository, trainRepo repository.TrainRepository) ScheduleService {
	return &scheduleService{
		scheduleRepo: scheduleRepo,
		trainRepo:    trainRepo,
	}
}

func (s *scheduleService) Create(req dto.CreateScheduleRequest) (*models.Schedule, error) {
	_, err := s.trainRepo.FindByID(req.TrainID)
	if err != nil {
		return nil, errors.New("kereta tidak ditemukan")
	}

	schedule := &models.Schedule{
		TrainID:          req.TrainID,
		DepartureStation: req.DepartureStation,
		ArrivalStation:   req.ArrivalStation,
		DepartureTime:    req.DepartureTime,
		ArrivalTime:      req.ArrivalTime,
		Price:            req.Price,
		AvailableSeats:   req.AvailableSeats,
	}

	if err := s.scheduleRepo.Create(schedule); err != nil {
		return nil, err
	}

	return schedule, nil
}

func (s *scheduleService) GetByID(id int) (*models.Schedule, error) {
	return s.scheduleRepo.FindByID(id)
}

func (s *scheduleService) GetAll() ([]models.Schedule, error) {
	return s.scheduleRepo.FindAll()
}

func (s *scheduleService) Search(req dto.SearchScheduleRequest) ([]models.Schedule, error) {
	return s.scheduleRepo.Search(req.DepartureStation, req.ArrivalStation, req.Date)
}

func (s *scheduleService) Update(id int, req dto.UpdateScheduleRequest) (*models.Schedule, error) {
	schedule, err := s.scheduleRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("schedule not found")
	}

	if req.TrainID != nil {
		_, err := s.trainRepo.FindByID(*req.TrainID)
		if err != nil {
			return nil, errors.New("train not found")
		}
		schedule.TrainID = *req.TrainID
	}

	if req.DepartureStation != nil {
		schedule.DepartureStation = *req.DepartureStation
	}

	if req.ArrivalStation != nil {
		schedule.ArrivalStation = *req.ArrivalStation
	}

	if req.DepartureTime != nil {
		schedule.DepartureTime = *req.DepartureTime
	}

	if req.ArrivalTime != nil {
		schedule.ArrivalTime = *req.ArrivalTime
	}

	if req.Price != nil {
		schedule.Price = *req.Price
	}

	if req.AvailableSeats != nil {
		schedule.AvailableSeats = *req.AvailableSeats
	}

	if err := s.scheduleRepo.Update(id, schedule); err != nil {
		return nil, err
	}

	return schedule, nil
}

func (s *scheduleService) Delete(id int) error {
	_, err := s.scheduleRepo.FindByID(id)
	if err != nil {
		return errors.New("schedule tidak ditemukan")
	}

	return s.scheduleRepo.Delete(id)
}