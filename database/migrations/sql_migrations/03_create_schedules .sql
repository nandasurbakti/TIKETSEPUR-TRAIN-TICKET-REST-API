-- +migrate Up
-- +migrate StatementBegin

create table schedules (
    id SERIAL PRIMARY KEY,
    train_id INT NOT NULL,
    departure_station VARCHAR(255) NOT NULL,
    arrival_station VARCHAR(255) NOT NULL,
    departure_time TIMESTAMP NOT NULL,
    arrival_time TIMESTAMP NOT NULL,
    price DECIMAL(13,2) NOT NULL,
    available_seats INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    modified_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_trains_schedules FOREIGN KEY (train_id) REFERENCES trains(id) ON DELETE CASCADE
)

-- +migrate StatementEnd