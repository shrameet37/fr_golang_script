BEGIN;

DROP TABLE IF EXISTS organisation_face_data_ids;
DROP TRIGGER IF EXISTS trigger_organisation_face_data_ids_set_updated_at_timestamp ON organisation_face_data_ids;

DROP FUNCTION IF EXISTS set_updated_at_timestamp() CASCADE;

DROP INDEX IF EXISTS "organisation_face_data_ids_accessorId_organisationId_index";
COMMIT;