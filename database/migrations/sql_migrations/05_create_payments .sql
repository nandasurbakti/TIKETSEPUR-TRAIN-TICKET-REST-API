-- +migrate Up
-- +migrate StatementBegin

create table payments (
    id SERIAL PRIMARY KEY,
    ticket_id INT NOT NULL,
    payment_method VARCHAR(50) NOT NULL,
    payment_amount DECIMAL(13,2) NOT NULL,
    payment_status VARCHAR(20) NOT NULL DEFAULT 'pending',
    payment_code VARCHAR(100) UNIQUE NOT NULL,
    paid_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    modified_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_payments_tickets FOREIGN KEY (ticket_id) REFERENCES tickets(id) ON DELETE RESTRICT
)

-- +migrate StatementEnd