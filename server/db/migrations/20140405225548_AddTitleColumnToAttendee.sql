
-- +goose Up
ALTER TABLE Attendee
  ADD COLUMN Title text not null DEFAULT 'General';

-- +goose Down
ALTER TABLE Attendee
  DROP COLUMN Title;
