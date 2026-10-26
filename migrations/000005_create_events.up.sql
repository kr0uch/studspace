CREATE TABLE events (
    event_id    BIGSERIAL PRIMARY KEY,
    title       TEXT           NOT NULL,
    organizer   TEXT           NOT NULL,
    category    event_category NOT NULL,
    space_id    SMALLINT       NOT NULL REFERENCES spaces (space_id),
    starts_at   TIMESTAMPTZ    NOT NULL,
    ends_at     TIMESTAMPTZ    NOT NULL,
    period      TSTZRANGE      GENERATED ALWAYS AS (tstzrange(starts_at, ends_at, '[)')) STORED,
    capacity    INTEGER        NOT NULL CHECK (capacity > 0),
    description TEXT,
    status      event_status   NOT NULL DEFAULT 'published',
    created_by  BIGINT         NOT NULL REFERENCES users (user_id),
    created_at  TIMESTAMPTZ    NOT NULL DEFAULT now(),

    CHECK (ends_at > starts_at),

    -- Two published events cannot occupy the same space at the same time.
    -- Overlap with approved bookings is checked in the approve/create transaction.
    CONSTRAINT events_no_published_overlap EXCLUDE USING gist (
        space_id WITH =,
        period WITH &&
    ) WHERE (status = 'published')
);

CREATE INDEX events_starts_idx ON events (starts_at);
CREATE INDEX events_space_period_idx ON events USING gist (space_id, period);

CREATE TABLE event_registrations (
    event_id     BIGINT      NOT NULL REFERENCES events (event_id) ON DELETE CASCADE,
    user_id      BIGINT      NOT NULL REFERENCES users (user_id),
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    cancelled_at TIMESTAMPTZ,
    PRIMARY KEY (event_id, user_id)
);

CREATE INDEX event_registrations_active_idx
    ON event_registrations (event_id)
    WHERE cancelled_at IS NULL;

CREATE INDEX event_registrations_user_idx
    ON event_registrations (user_id)
    WHERE cancelled_at IS NULL;

CREATE TABLE event_waitlist (
    event_id    BIGINT      NOT NULL REFERENCES events (event_id) ON DELETE CASCADE,
    user_id     BIGINT      NOT NULL REFERENCES users (user_id),
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    notified_at TIMESTAMPTZ,
    PRIMARY KEY (event_id, user_id)
);

CREATE INDEX event_waitlist_pending_idx
    ON event_waitlist (event_id, created_at)
    WHERE notified_at IS NULL;
