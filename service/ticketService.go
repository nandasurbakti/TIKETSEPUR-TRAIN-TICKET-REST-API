package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"tiketsepur/dto"
	"tiketsepur/models"
	"tiketsepur/repository"
	"tiketsepur/utils"
	"time"

	"github.com/jmoiron/sqlx"
)

type TicketService interface {
	Create(ctx context.Context, userID int, req dto.CreateTicketRequest) (*models.Ticket, error)
	GetByID(id int) (*models.TicketWithDetails, error)
	GetByBookingCode(code string) (*models.TicketWithDetails, error)
	GetByUserID(userID int) ([]models.TicketWithDetails, error)
	GetAll() ([]models.TicketWithDetails, error)
	Cancel(ctx context.Context, id int, userID int, role string) error
}

type ticketService struct {
	db           *sqlx.DB
	ticketRepo   repository.TicketRepository
	scheduleRepo repository.ScheduleRepository
	paymentRepo  repository.PaymentRepository
	userRepo     repository.UserRepository
	redis        *utils.RedisClient
	rabbitmq     *utils.RabbitMQ
}

func NewTicketService(
	db *sqlx.DB,
	ticketRepo repository.TicketRepository,
	scheduleRepo repository.ScheduleRepository,
	userRepo repository.UserRepository,
	paymentRepo repository.PaymentRepository,
	redis *utils.RedisClient,
	rabbitmq *utils.RabbitMQ,
) TicketService {
	return &ticketService{
		db:           db,
		ticketRepo:   ticketRepo,
		scheduleRepo: scheduleRepo,
		paymentRepo:  paymentRepo,
		userRepo:     userRepo,
		redis:        redis,
		rabbitmq:     rabbitmq,
	}
}

func (s *ticketService) Create(ctx context.Context, userID int, req dto.CreateTicketRequest) (*models.Ticket, error) {
	lockKey := fmt.Sprintf("lock:seat:%d:%s", req.ScheduleID, req.SeatNumber)
	lockValue := fmt.Sprintf("%d-%d", userID, time.Now().Unix())
	
	locked, err := s.acquireLock(ctx, lockKey, lockValue, 10*time.Second)
	if err != nil || !locked {
		return nil, errors.New("kursi sudah dibeli oleh pengguna lain, silahkan coba lagi")
	}
	defer s.releaseLock(ctx, lockKey, lockValue)

	schedule, err := s.scheduleRepo.FindByID(req.ScheduleID)
	if err != nil {
		return nil, errors.New("schedule tidak ditemukan")
	}

	available, err := s.ticketRepo.CheckSeatAvailability(req.ScheduleID, req.SeatNumber)
	if err != nil {
		return nil, err
	}
	if !available {
		return nil, errors.New("kursi sudah dipesan, silahkan pilih kursi lain")
	}

	tx, err := s.db.Beginx()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	if err := s.scheduleRepo.DecrementSeat(req.ScheduleID, tx); err != nil {
		return nil, errors.New("tidak ada kursi yang tersedia")
	}

	bookingCode := s.generateBookingCode()

	ticket := &models.Ticket{
		UserID:            userID,
		ScheduleID:        req.ScheduleID,
		SeatNumber:        req.SeatNumber,
		PassengerName:     req.PassengerName,
		PassengerIDNumber: req.PassengerIDNumber,
		Status:            "pending",
		BookingCode:       bookingCode,
		TotalPrice:        schedule.Price,
	}

	if err := s.ticketRepo.Create(ticket, tx); err != nil {
		return nil, err
	}

	paymentCode := s.generatePaymentCode()
	payment := &models.Payment{
		TicketID:      ticket.ID,
		PaymentMethod: req.PaymentMethod,
		PaymentAmount: schedule.Price,
		PaymentStatus: "pending",
		PaymentCode:   paymentCode,
	}

	if err := s.paymentRepo.Create(payment, tx); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	user, err := s.userRepo.FindByID(ticket.UserID)
	if err != nil {
		log.Printf("Warning: tidak bisa mendapatkan user untuk notifikasi: %v", err)
	} else {
		go s.sendBookingNotification(user, ticket, schedule)
	}

	return ticket, nil
}

func (s *ticketService) GetByID(id int) (*models.TicketWithDetails, error) {
	return s.ticketRepo.FindByID(id)
}

func (s *ticketService) GetByBookingCode(code string) (*models.TicketWithDetails, error) {
	return s.ticketRepo.FindByBookingCode(code)
}

func (s *ticketService) GetByUserID(userID int) ([]models.TicketWithDetails, error) {
	return s.ticketRepo.FindByUserID(userID)
}

func (s *ticketService) GetAll() ([]models.TicketWithDetails, error) {
	return s.ticketRepo.FindAll()
}

func (s *ticketService) Cancel(ctx context.Context, id int, userID int, role string) error {
	ticket, err := s.ticketRepo.FindByID(id)
	if err != nil {
		return errors.New("tiket tidak ditemukan")
	}

	if role != "admin" && ticket.UserID != userID {
		return errors.New("tidak ada wewenang untuk membatalkan tiket ini")
	}

	if ticket.Status == "cancelled" {
		return errors.New("tiket sudah dibatalkan")
	}

	lockKey := fmt.Sprintf("lock:cancel:%d", id)
	lockValue := fmt.Sprintf("%d-%d", userID, time.Now().Unix())
	
	locked, err := s.acquireLock(ctx, lockKey, lockValue, 10*time.Second)
	if err != nil || !locked {
		return errors.New("pembatalan sedang diproses, silakan coba lagi")
	}
	defer s.releaseLock(ctx, lockKey, lockValue)

	tx, err := s.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := s.ticketRepo.UpdateStatus(id, "cancelled", tx); err != nil {
		return err
	}

	if err := s.scheduleRepo.IncrementSeat(ticket.ScheduleID, tx); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	user, err := s.userRepo.FindByID(ticket.UserID)
if err == nil {
    go func() {
        notification := utils.NotificationMessage{
            Type:        "cancellation",
            Email:       user.Email,
            BookingCode: ticket.BookingCode,
            TrainName:   ticket.TrainName,
        }
        if err := s.rabbitmq.PublishNotification(notification); err != nil {
            log.Printf("gagal untuk mempublish pembatalan: %v", err)
        }
    }()
}

	return nil
}



func (s *ticketService) acquireLock(ctx context.Context, key, value string, expiry time.Duration) (bool, error) {
	result, err := s.redis.Get(ctx, key)
	if err == nil && result != "" {
		return false, nil 
	}

	err = s.redis.Set(ctx, key, value, expiry)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (s *ticketService) releaseLock(ctx context.Context, key, value string) error {
	currentValue, err := s.redis.Get(ctx, key)
	if err != nil {
		return err
	}

	if currentValue == value {
		return s.redis.Delete(ctx, key)
	}

	return nil
}

func (s *ticketService) generateBookingCode() string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	
	code := make([]byte, 8)
	for i := range code {
		code[i] = charset[rand.Intn(len(charset))]
	}
	
	return "TRN" + string(code)
}

func (s *ticketService) generatePaymentCode() string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	
	code := make([]byte, 12)
	for i := range code {
		code[i] = charset[rand.Intn(len(charset))]
	}
	
	return "PAY" + string(code)
}

func (s *ticketService) sendBookingNotification(user *models.User,ticket *models.Ticket, schedule *models.Schedule) {
	notification := utils.NotificationMessage{
		Type:        "booking",
		Email:       user.Email,
		BookingCode: ticket.BookingCode,
		TrainName:   schedule.TrainName,
		Departure:   schedule.DepartureStation,
		Arrival:     schedule.ArrivalStation,
		SeatNumber:  ticket.SeatNumber,
		TotalPrice:  ticket.TotalPrice,
	}

if err := s.rabbitmq.PublishNotification(notification); err != nil {
	} else {
		log.Printf("[BOOKING] Sending notification to %s", user.Email)
		log.Printf("Booking Code: %s", ticket.BookingCode)
		log.Printf("Train: %s", schedule.TrainName)
	}
}