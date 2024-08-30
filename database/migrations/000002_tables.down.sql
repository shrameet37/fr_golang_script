BEGIN;



DROP TABLE IF EXISTS accessor_permissions;
DROP TRIGGER IF EXISTS trigger_accessor_permissions_set_updated_at_timestamp ON accessor_permissions;

DROP TABLE IF EXISTS device_permissions;
DROP TRIGGER IF EXISTS trigger_device_permissions_set_updated_at_timestamp ON device_permissions;

DROP TABLE IF EXISTS deleted_device_permissions;
DROP TRIGGER IF EXISTS trigger_deleted_device_permissions_set_updated_at_timestamp ON deleted_device_permissions;

DROP TABLE IF EXISTS access_points;
DROP TRIGGER IF EXISTS trigger_access_points_set_updated_at_timestamp ON access_points;

DROP TABLE IF EXISTS devices;
DROP TRIGGER IF EXISTS trigger_devices_set_updated_at_timestamp ON devices;

DROP TABLE IF EXISTS access_point_devices;
DROP TRIGGER IF EXISTS trigger_access_point_devices_set_updated_at_timestamp ON access_point_devices;

DROP TABLE IF EXISTS organisation_accessors;
DROP TRIGGER IF EXISTS trigger_organisation_accessors_set_updated_at_timestamp ON organisation_accessors;

DROP TABLE IF EXISTS face_data_id_mapping;
DROP TRIGGER IF EXISTS trigger_face_data_id_mapping_set_updated_at_timestamp ON face_data_id_mapping;

DROP TABLE IF EXISTS deleted_face_data_id_mapping;
DROP TRIGGER IF EXISTS trigger_deleted_face_data_id_mapping_set_updated_at_timestamp ON deleted_face_data_id_mapping;

DROP INDEX IF EXISTS "accessorId_accessPointId_faceDataId_index";
DROP INDEX IF EXISTS "serialNumber_faceDataId_index";
DROP INDEX IF EXISTS "accessPointId_organisationId_index";
DROP INDEX IF EXISTS "serialNumber_index";
DROP INDEX IF EXISTS "serialNumber_accessPointId_index";
DROP INDEX IF EXISTS "accessorId_organisationId_index";
DROP INDEX IF EXISTS "faceDataId_accessPoint_index";

COMMIT;