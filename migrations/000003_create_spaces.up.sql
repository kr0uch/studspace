CREATE TABLE spaces (
    space_id      SMALLSERIAL PRIMARY KEY,
    code          TEXT    NOT NULL UNIQUE,
    name          TEXT    NOT NULL,
    kind          TEXT    NOT NULL,
    capacity      INTEGER NOT NULL CHECK (capacity > 0),
    photo_url     TEXT,
    color         TEXT    NOT NULL,
    open_time     TIME    NOT NULL DEFAULT '08:00',
    close_time    TIME    NOT NULL DEFAULT '22:00',
    slot_step_min INTEGER NOT NULL DEFAULT 30 CHECK (slot_step_min > 0),
    is_active     BOOLEAN NOT NULL DEFAULT true,
    CHECK (close_time > open_time)
);

CREATE TABLE amenities (
    amenity_id SMALLSERIAL PRIMARY KEY,
    name       TEXT NOT NULL UNIQUE,
    icon       TEXT NOT NULL
);

CREATE TABLE space_amenities (
    space_id   SMALLINT NOT NULL REFERENCES spaces (space_id) ON DELETE CASCADE,
    amenity_id SMALLINT NOT NULL REFERENCES amenities (amenity_id) ON DELETE CASCADE,
    PRIMARY KEY (space_id, amenity_id)
);
