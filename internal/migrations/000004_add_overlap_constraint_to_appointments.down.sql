ALTER TABLE appointments DROP CONSTRAINT IF EXISTS no_overlapping_appointments;

DROP EXTENSION IF EXISTS btree_gist;