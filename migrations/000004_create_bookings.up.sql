CREATE TABLE bookings (
    booking_id          BIGSERIAL PRIMARY KEY,
    user_id             BIGINT         NOT NULL REFERENCES users (user_id),
    space_id            SMALLINT       NOT NULL REFERENCES spaces (space_id),
    first_name          TEXT           NOT NULL,
    last_name           TEXT           NOT NULL,
    starts_at           TIMESTAMPTZ    NOT NULL,
    ends_at             TIMESTAMPTZ    NOT NULL,
    period              TSTZRANGE      GENERATED ALWAYS AS (tstzrange(starts_at, ends_at, '[)')) STORED,
    purpose             VARCHAR(150)   NOT NULL CHECK (char_length(purpose) > 0),
    people_count        INTEGER        NOT NULL CHECK (people_count > 0),
    status              booking_status NOT NULL DEFAULT 'pending',
    reject_reason_title TEXT,
    reject_reason_text  VARCHAR(500),
    decided_by          BIGINT REFERENCES users (user_id),
    decided_at          TIMESTAMPTZ,
    cancelled_by        cancel_actor,
    cancelled_at        TIMESTAMPTZ,
    rescheduled_from_id BIGINT REFERENCES bookings (booking_id),
    created_at          TIMESTAMPTZ    NOT NULL DEFAULT now(),
    updated_at          TIMESTAMPTZ    NOT NULL DEFAULT now(),

    CHECK (ends_at > starts_at),
    CHECK (status <> 'rejected' OR reject_reason_text IS NOT NULL),
    CHECK (status <> 'cancelled' OR (cancelled_by IS NOT NULL AND cancelled_at IS NOT NULL)),

    -- Two admins cannot approve overlapping bookings for the same space.
    CONSTRAINT bookings_no_approved_overlap EXCLUDE USING gist (
        space_id WITH =,
        period WITH &&
    ) WHERE (status = 'approved')
);

-- «Мои заявки»
CREATE INDEX bookings_user_status_starts_idx ON bookings (user_id, status, starts_at DESC);
-- админские вкладки Новые / Одобренные / Отклонённые
CREATE INDEX bookings_status_created_idx ON bookings (status, created_at);
-- расписание и поиск конфликтов, в том числе среди pending
CREATE INDEX bookings_space_period_idx ON bookings USING gist (space_id, period);

CREATE TRIGGER bookings_set_updated_at
    BEFORE UPDATE ON bookings
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

CREATE TABLE booking_status_log (
    booking_status_log_id BIGSERIAL PRIMARY KEY,
    booking_id            BIGINT         NOT NULL REFERENCES bookings (booking_id) ON DELETE CASCADE,
    from_status           booking_status,
    to_status             booking_status NOT NULL,
    actor_user_id         BIGINT         NOT NULL REFERENCES users (user_id),
    comment       TEXT,
    created_at    TIMESTAMPTZ    NOT NULL DEFAULT now()
);

CREATE INDEX booking_status_log_booking_idx ON booking_status_log (booking_id, created_at);
