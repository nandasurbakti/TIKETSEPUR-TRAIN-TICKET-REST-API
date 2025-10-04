-- +migrate Up
ALTER TABLE tickets DROP CONSTRAINT IF EXISTS tickets_schedule_id_key;
ALTER TABLE tickets DROP CONSTRAINT IF EXISTS tickets_seat_number_key;

ALTER TABLE tickets 
ADD CONSTRAINT unique_schedule_seat UNIQUE (schedule_id, seat_number);

-- +migrate Down
ALTER TABLE tickets DROP CONSTRAINT IF EXISTS unique_schedule_seat;

ALTER TABLE tickets 
ADD CONSTRAINT tickets_schedule_id_key UNIQUE (schedule_id);

ALTER TABLE tickets 
ADD CONSTRAINT tickets_seat_number_key UNIQUE (seat_number);
