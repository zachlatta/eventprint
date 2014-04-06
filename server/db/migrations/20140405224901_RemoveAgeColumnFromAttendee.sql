
-- +goose Up
ALTER TABLE Attendee
  DROP COLUMN Age;


-- +goose Down
ALTER TABLE Attendee
  ADD COLUMN Age text;
