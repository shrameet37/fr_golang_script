BEGIN;

ALTER TABLE device_permissions DROP COLUMN IF EXISTS "assignedAt";
ALTER TABLE device_permissions ADD COLUMN "assignedAt" TIMESTAMPTZ;

COMMIT;