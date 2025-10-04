-- +migrate Up
-- +migrate StatementBegin

create table tickets (
    id SERIAL PRIMARY KEY,
    user_id INT,
    schedule_id INT UNIQUE NOT NULL,
    seat_number VARCHAR(10) UNIQUE NOT NULL,
    passenger_name VARCHAR(255) NOT NULL,
    passenger_id_number VARCHAR(50) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    booking_code VARCHAR(50) UNIQUE NOT NULL,
    total_price DECIMAL(13,2) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    modified_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_tickets_users FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET NULL,
    CONSTRAINT fk_tickets_schedules FOREIGN KEY (schedule_id) REFERENCES schedules(id) ON DELETE CASCADE
)

-- +migrate StatementEnd