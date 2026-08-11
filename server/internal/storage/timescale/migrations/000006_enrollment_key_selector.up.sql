ALTER TABLE enrollment_keys ADD COLUMN selector TEXT NOT NULL DEFAULT '';

ALTER TABLE enrollment_keys ALTER COLUMN selector DROP DEFAULT;