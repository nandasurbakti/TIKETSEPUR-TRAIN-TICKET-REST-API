package service

import (
	"context"
	"errors"
	config "tiketsepur/configs"
	"tiketsepur/dto"
	"tiketsepur/models"
	"tiketsepur/repository"
	"tiketsepur/utils"

	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Register(req dto.RegisterRequest) (*models.User, error)
	Login(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error)
	Logout(ctx context.Context, token string) error
	ValidateToken(ctx context.Context, token string) (*utils.Claims, error)
}

type authService struct {
	userRepo repository.UserRepository
	redis    *utils.RedisClient
	config   *config.Config
}

func NewAuthService(userRepo repository.UserRepository, redis *utils.RedisClient, cfg *config.Config) AuthService {
	return &authService{
		userRepo: userRepo,
		redis:    redis,
		config:   cfg,
	}
}

func (s *authService) Register(req dto.RegisterRequest) (*models.User, error) {
	existingUser, _ := s.userRepo.FindByEmail(req.Email)
	if existingUser != nil {
		return nil, errors.New("email already registered")
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
		Role:     "user",
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *authService) Login(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error) {
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		return nil, errors.New("credentials tidak valid ")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("credentials tidak valid ")
	}

	token, err := utils.GenerateJWT(user.ID, user.Email, user.Role, s.config.JWT.Secret, s.config.JWT.Exp)
	if err != nil {
		return nil, err
	}

	tokenKey := "token:" + token
	if err := s.redis.Set(ctx, tokenKey, user.ID, s.config.Redis.SessionExpiry); err != nil {
		return nil, err
	}

	sessionKey := "session:" + string(rune(user.ID))
	if err := s.redis.Set(ctx, sessionKey, "active", s.config.Redis.SessionExpiry); err != nil {
		return nil, err
	}

	return &dto.LoginResponse{
		Token: token,
		User: map[string]interface{}{
			"id":        user.ID,
			"email":     user.Email,
			"full_name": user.FullName,
			"role":      user.Role,
		},
	}, nil
}

func (s *authService) Logout(ctx context.Context, token string) error {
	tokenKey := "token:" + token
	return s.redis.Delete(ctx, tokenKey)
}

func (s *authService) ValidateToken(ctx context.Context, token string) (*utils.Claims, error) {
	tokenKey := "token:" + token
	exists, err := s.redis.Exists(ctx, tokenKey)
	if err != nil || !exists {
		return nil, errors.New("token tidak ditemukan atau expired")
	}

	claims, err := utils.ValidateJWT(token, s.config.JWT.Secret)
	if err != nil {
		return nil, err
	}

	return claims, nil
}