-- +goose Up
CREATE TABLE halls (
    id UUID PRIMARY KEY,
    venue_id UUID NOT NULL REFERENCES venues (id) ON DELETE RESTRICT,
    name TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,
    CONSTRAINT halls_venue_id_name_unique UNIQUE (venue_id, name)
);

CREATE TABLE hall_sections (
    hall_id UUID NOT NULL REFERENCES halls (id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    rows_count INT NOT NULL CHECK (rows_count > 0),
    seats_per_row INT NOT NULL CHECK (seats_per_row > 0),
    PRIMARY KEY (hall_id, name)
);

CREATE TABLE seats (
    id UUID PRIMARY KEY,
    hall_id UUID NOT NULL REFERENCES halls (id) ON DELETE CASCADE,
    section TEXT NOT NULL,
    row_number INT NOT NULL CHECK (row_number > 0),
    number INT NOT NULL CHECK (number > 0),
    CONSTRAINT seats_position_unique UNIQUE (hall_id, section, row_number, number)
);

-- +goose Down
DROP TABLE seats;

DROP TABLE hall_sections;

DROP TABLE halls;
