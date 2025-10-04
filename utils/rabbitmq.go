package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQ struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	queue   amqp.Queue
}

type NotificationMessage struct {
	Type        string  `json:"type"`
	Email       string  `json:"email"`
	BookingCode string  `json:"booking_code"`
	TrainName   string  `json:"train_name"`
	Departure   string  `json:"departure"`
	Arrival     string  `json:"arrival"`
	SeatNumber  string  `json:"seat_number"`
	TotalPrice  float64 `json:"total_price"`
	PaymentCode string  `json:"payment_code,omitempty"`
	PaymentMethod string `json:"payment_method,omitempty"`
	DepartureTime string  `json:"departure_time,omitempty"`
}

func NewRabbitMQ(url, queueName string) (*RabbitMQ, error) {
	rabbitURL := os.Getenv("RABBITMQ_URL")
	if rabbitURL == "" {
		return nil, fmt.Errorf("RABBITMQ_URL environment variable belum diatur")
	}
	
	conn, err := amqp.Dial(rabbitURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect: %w", err)
	}
	channel, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	queue, err := channel.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}
	return &RabbitMQ{
		conn:    conn,
		channel: channel,
		queue: queue,
	}, nil
}

func (r *RabbitMQ) PublishNotification(msg NotificationMessage) error {
	body, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	return r.channel.Publish(
		"",
		r.queue.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body: body,
			DeliveryMode: amqp.Persistent,
		},
	)
}

func (r *RabbitMQ) ConsumeNotifications() {
	msgs, err := r.channel.Consume(
		r.queue.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Print("gagal mendaftarkan konsumen: ", err)
		return
	}
	go func() {
		for msg := range msgs {
			var notification NotificationMessage
			if err := json.Unmarshal(msg.Body, &notification); err != nil {
				log.Print("error: ", err)
				continue
			}
			switch notification.Type {
			case "booking":
				r.handleBookingNotification(notification)
			case "payment":
				r.handlePaymentNotification(notification)
			case "cancellation":
				r.handleCancellationNotification(notification)
			default:
				log.Printf("Unknown notification type: %s", notification.Type)
			}
		}
	}()
}

func (r *RabbitMQ) handleBookingNotification(n NotificationMessage) {
	log.Printf("[BOOKING] Sending notification to %s", n.Email)
	log.Printf("Booking Code: %s", n.BookingCode)
	log.Printf("Train: %s", n.TrainName)
	log.Printf("Route: %s → %s", n.Departure, n.Arrival)
	log.Printf("Seat: %s | Price: Rp %.0f", n.SeatNumber, n.TotalPrice)
}

func (r *RabbitMQ) handlePaymentNotification(n NotificationMessage) {
	log.Printf("[PAYMENT] Sending payment confirmation to %s", n.Email)
	log.Printf("Booking Code: %s", n.BookingCode)
	log.Printf("Payment Code: %s", n.PaymentCode)
	log.Printf("Amount: Rp %.0f", n.TotalPrice)
	log.Printf("Train: %s (%s → %s)", n.TrainName, n.Departure, n.Arrival)
	log.Printf("Departure: %s", n.DepartureTime)
	log.Printf("Seat: %s", n.SeatNumber)
	log.Printf("Payment Confirmed - Ticket is now ACTIVE")
}

func (r *RabbitMQ) handleCancellationNotification(n NotificationMessage) {
	log.Printf("[CANCELLATION] Sending cancellation notification to %s", n.Email)
	log.Printf("Booking Code: %s", n.BookingCode)
	log.Printf("Train: %s", n.TrainName)
}

func (r *RabbitMQ) Close() {
	r.channel.Close()
	r.conn.Close()
}

