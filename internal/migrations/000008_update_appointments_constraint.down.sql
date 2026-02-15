ALTER TABLE appointments 
DROP CONSTRAINT IF EXISTS no_overlapping_active_appointments;

ALTER TABLE appointments 
ADD CONSTRAINT no_overlapping_globally
EXCLUDE USING gist (
    tstzrange(start_time, end_time) WITH &&
);