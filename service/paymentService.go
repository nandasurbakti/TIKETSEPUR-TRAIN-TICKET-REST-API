package service

import (
	"context"
	"errors"
	"log"
	"tiketsepur/models"
	"tiketsepur/repository"
	"tiketsepur/utils"
	"time"

	"github.com/jmoiron/sqlx"
)

type PaymentService interface {
	ConfirmPayment(ctx context.Context, paymentCode string) error
}

type paymentService struct {
	db           *sqlx.DB
	paymentRepo  repository.PaymentRepository
	ticketRepo   repository.TicketRepository
	scheduleRepo repository.ScheduleRepository
	userRepo     repository.UserRepository
	rabbitmq     *utils.RabbitMQ
}

func NewPaymentService(
	db *sqlx.DB,
	paymentRepo repository.PaymentRepository,
	ticketRepo repository.TicketRepository,
	scheduleRepo repository.ScheduleRepository,
	userRepo repository.UserRepository,
	rabbitmq *utils.RabbitMQ,
) PaymentService {
	return &paymentService{
		db:           db,
		paymentRepo:  paymentRepo,
		ticketRepo:   ticketRepo,
		scheduleRepo: scheduleRepo,
		userRepo:     userRepo,
		rabbitmq:     rabbitmq,
	}
}

func (s *paymentService) ConfirmPayment(ctx context.Context, paymentCode string) error {
	payment, err := s.paymentRepo.FindByPaymentCode(paymentCode)
	if err != nil {
		return errors.New("payment tidak ditemukan")
	}

	if payment.PaymentStatus == "success" {
		return errors.New("payment sudah dikonfirmasi")
	}

	ticket, err := s.ticketRepo.FindByID(payment.TicketID)
	if err != nil {
		return errors.New("tiket tidak ditemukan")
	}

	schedule, err := s.scheduleRepo.FindByID(ticket.ScheduleID)
	if err != nil {
		return errors.New("schedule tidak ditemukan")
	}

	user, err := s.userRepo.FindByID(ticket.UserID)
	if err != nil {
		return errors.New("user tidak ditemukan")
	}

	tx, err := s.db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	now := time.Now()
	if err := s.updatePaymentStatus(tx, payment.ID, "success", &now); err != nil {
		return err
	}

	if err := s.ticketRepo.UpdateStatus(ticket.ID, "confirmed", tx); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	go s.sendPaymentConfirmationNotification(user, ticket, schedule, payment)

	return nil
}

func (s *paymentService) updatePaymentStatus(tx *sqlx.Tx, paymentID int, status string, paidAt *time.Time) error {
	query := `UPDATE payments 
			  SET payment_status = $1, paid_at = $2, modified_at = NOW() 
			  WHERE id = $3`
	_, err := tx.Exec(query, status, paidAt, paymentID)
	return err
}

func (s *paymentService) sendPaymentConfirmationNotification(
	user *models.User,
	ticket *models.TicketWithDetails,
	schedule *models.Schedule,
	payment *models.Payment,
) {
	notification := utils.NotificationMessage{
		Type:          "payment",
		Email:         user.Email,
		BookingCode:   ticket.BookingCode,
		TrainName:     schedule.TrainName,
		Departure:     schedule.DepartureStation,
		Arrival:       schedule.ArrivalStation,
		SeatNumber:    ticket.SeatNumber,
		TotalPrice:    payment.PaymentAmount,
		PaymentCode:   payment.PaymentCode,
		PaymentMethod: payment.PaymentMethod,
		DepartureTime: schedule.DepartureTime.Format("2006-01-02 15:04"),
	}

	if err := s.rabbitmq.PublishNotification(notification); err != nil {
		log.Printf("gagal mengirim notifikasi pembayaran: %v", err)
	} else {
		log.Printf("pemberitahuan konfirmasi pembayaran dikirim ke %s", user.Email)
	}
}
