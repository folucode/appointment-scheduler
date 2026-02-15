CREATE EXTENSION IF NOT EXISTS btree_gist;

ALTER TABLE appointments 
ADD CONSTRAINT no_overlapping_globally
EXCLUDE USING gist (
    tsrange(start_time, end_time) WITH &&
);