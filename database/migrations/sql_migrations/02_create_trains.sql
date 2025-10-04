-- +migrate Up
-- +migrate StatementBegin

create table trains (
    id SERIAL PRIMARY KEY,
    train_code VARCHAR(50) UNIQUE NOT NULL,
    train_name VARCHAR(255) NOT NULL,
    train_type VARCHAR(50) NOT NULL,
    total_seats INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    modified_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
)

-- +migrate StatementEnd