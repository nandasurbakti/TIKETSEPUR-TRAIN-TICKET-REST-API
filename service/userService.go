package service

import (
	"errors"
	"tiketsepur/dto"
	"tiketsepur/models"
	"tiketsepur/repository"

	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Create(req dto.CreateUserRequest) (*models.User, error)
	GetByID(id int) (*models.User, error)
	GetAll() ([]models.User, error)
	Update(id int, req dto.UpdateUserRequest) (*models.User, error)
	Delete(id int) error
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

func (s *userService) Create(req dto.CreateUserRequest) (*models.User, error) {
	existingUser, _ := s.userRepo.FindByEmail(req.Email)
	if existingUser != nil {
		return nil, errors.New("email sudah digunakan")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Email:    req.Email,
		Password: string(hashedPassword),
		FullName: req.FullName,
		Phone:    req.Phone,
		Role:     req.Role,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) GetByID(id int) (*models.User, error) {
	return s.userRepo.FindByID(id)
}

func (s *userService) GetAll() ([]models.User, error) {
	return s.userRepo.FindAll()
}

func (s *userService) Update(id int, req dto.UpdateUserRequest) (*models.User, error) {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("user tidak ditemukan")
	}

	if req.Email != nil {
		if *req.Email != user.Email {
			existingUser, _ := s.userRepo.FindByEmail(*req.Email)
			if existingUser != nil {
				return nil, errors.New("email sudah digunakan")
			}
		}
		user.Email = *req.Email
	}

	if req.FullName != nil {
		user.FullName = *req.FullName
	}

	if req.Phone != nil {
		user.Phone = *req.Phone
	}

	if req.Role != nil {
		user.Role = *req.Role
	}

	if err := s.userRepo.Update(id, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) Delete(id int) error {
	_, err := s.userRepo.FindByID(id)
	if err != nil {
		return errors.New("user tidak ditemukan")
	}

	return s.userRepo.Delete(id)
}