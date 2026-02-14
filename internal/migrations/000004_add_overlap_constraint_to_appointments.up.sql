CREATE EXTENSION IF NOT EXISTS btree_gist;

ALTER TABLE appointments 
ADD CONSTRAINT no_overlapping_appointments 
EXCLUDE USING gist (
    user_id WITH =,
    tstzrange(start_time, end_time) WITH &&
);