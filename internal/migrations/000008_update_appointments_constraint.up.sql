ALTER TABLE appointments 
DROP CONSTRAINT IF EXISTS no_overlapping_globally;

ALTER TABLE appointments 
ADD CONSTRAINT no_overlapping_active_appointments
EXCLUDE USING gist (
    tstzrange(start_time, end_time) WITH &&
) 
WHERE (deleted_at IS NULL);