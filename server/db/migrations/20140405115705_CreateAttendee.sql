
-- +goose Up
CREATE TABLE Attendee (
  Id serial not null primary key,
  EventbriteId integer not null unique,
  FirstName text not null,
  LastName text not null,
  Gender text,
  Age text,
  Email text,
  CheckedIn bool,
  Barcode text not null
);


-- +goose Down
DROP TABLE Attendee;
