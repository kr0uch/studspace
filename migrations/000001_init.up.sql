-- Needed for EXCLUDE constraints that mix scalar equality (space_id =) with range overlap (&&).
CREATE EXTENSION IF NOT EXISTS btree_gist;

CREATE OR REPLACE FUNCTION set_updated_at() RETURNS trigger AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TYPE user_role AS ENUM ('student', 'admin');
CREATE TYPE booking_status AS ENUM ('pending', 'approved', 'rejected', 'cancelled', 'withdrawn', 'completed');
CREATE TYPE cancel_actor AS ENUM ('user', 'admin');
CREATE TYPE event_category AS ENUM ('study', 'career', 'leisure');
CREATE TYPE event_status AS ENUM ('published', 'cancelled');
