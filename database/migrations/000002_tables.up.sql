BEGIN;


CREATE TABLE IF NOT EXISTS accessor_permissions (
    "id" SERIAL PRIMARY KEY,
    "accessorId" INT,
    "accessPointId" INT,
    "organisationId" INT,
    "faceDataId" INT,
    "createdAt" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "updatedAt" TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE UNIQUE INDEX IF NOT EXISTS "accessorId_accessPointId_faceDataId_index" ON accessor_permissions("accessorId", "accessPointId", "faceDataId");


CREATE TRIGGER trigger_accessor_permissions_set_updated_at_timestamp BEFORE UPDATE ON accessor_permissions FOR EACH ROW EXECUTE FUNCTION set_updated_at_timestamp();



CREATE TABLE IF NOT EXISTS device_permissions (
    "id" SERIAL PRIMARY KEY,
    "serialNumber" VARCHAR(20),
    "credentialId" INT,
    "faceDataId" INT,
    "accessPointId" INT,
    "channelNo" INT,
    "accessorId" INT,
    "organisationId" INT,
    "updatedOnDevice" INT,
    "toDelete" INT,
    "assignedAt" TIMESTAMPTZ,
    "createdAt" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "updatedAt" TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS "serialNumber_faceDataId_index" ON device_permissions("serialNumber", "faceDataId" );


CREATE TRIGGER trigger_device_permissions_set_updated_at_timestamp BEFORE UPDATE ON device_permissions FOR EACH ROW EXECUTE FUNCTION set_updated_at_timestamp();


CREATE TABLE IF NOT EXISTS deleted_device_permissions (
    "id" SERIAL PRIMARY KEY,
    "deletedDevicePermissionsTableId" INT,
    "serialNumber" VARCHAR(20),
    "faceDataId" INT,
    "credentialId" INT,
    "accessorId" INT,
    "keyId" INT,
    "accessPointId" INT,
    "channelNo" INT,
    "organisationId" INT,
    "assignedAt" INT,
    "unAssignedAt" INT,
    "createdAt" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "updatedAt" TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TRIGGER trigger_deleted_device_permissions_set_updated_at_timestamp BEFORE UPDATE ON deleted_device_permissions FOR EACH ROW EXECUTE FUNCTION set_updated_at_timestamp();


CREATE TABLE IF NOT EXISTS access_points (
    "id" SERIAL PRIMARY KEY,
    "accessPointId" INT,
    "organisationId" INT,
    "siteId" INT,
    "configuration" INT,
    "channelNo" INT NOT NULL DEFAULT 0,
    "createdAt" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "updatedAt" TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS "accessPointId_organisationId_index" ON access_points("accessPointId", "organisationId" );



CREATE TRIGGER trigger_access_points_set_updated_at_timestamp BEFORE UPDATE ON access_points FOR EACH ROW EXECUTE FUNCTION set_updated_at_timestamp();



CREATE TABLE IF NOT EXISTS devices (
    "id" SERIAL PRIMARY KEY,
    "serialNumber" VARCHAR(20),
    "deviceType" INT,
    "organisationId" INT,
    "createdAt" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "updatedAt" TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS "serialNumber_index" ON devices("serialNumber");


CREATE TRIGGER trigger_devices_set_updated_at_timestamp BEFORE UPDATE ON devices FOR EACH ROW EXECUTE FUNCTION set_updated_at_timestamp();


CREATE TABLE IF NOT EXISTS access_point_devices (
    "id" SERIAL PRIMARY KEY,
    "serialNumber" VARCHAR(20),
    "accessPointId" INT,
    "createdAt" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "updatedAt" TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS "serialNumber_accessPointId_index" ON access_point_devices("serialNumber", "accessPointId");


CREATE TRIGGER trigger_access_point_devices_set_updated_at_timestamp BEFORE UPDATE ON access_point_devices FOR EACH ROW EXECUTE FUNCTION set_updated_at_timestamp();


CREATE TABLE IF NOT EXISTS organisation_accessors (
    "id" SERIAL PRIMARY KEY,
    "accessorId" INT,
    "organisationId" INT,
    "credentialId" INT,
    "createdAt" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "updatedAt" TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS "accessorId_organisationId_index" ON organisation_accessors("accessorId", "organisationId");


CREATE TRIGGER trigger_organisation_accessors_set_updated_at_timestamp BEFORE UPDATE ON organisation_accessors FOR EACH ROW EXECUTE FUNCTION set_updated_at_timestamp();

CREATE TABLE IF NOT EXISTS face_data_id_mapping (
    "id" SERIAL PRIMARY KEY,
    "faceDataId" INT,
    "accessPointId" INT,
    "accessorId" INT,
    "assignedAt" INT,
    "createdAt" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "updatedAt" TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS "faceDataId_accessPoint_index" ON face_data_id_mapping("faceDataId", "accessPointId");


CREATE TRIGGER trigger_face_data_id_mapping_set_updated_at_timestamp BEFORE UPDATE ON face_data_id_mapping FOR EACH ROW EXECUTE FUNCTION set_updated_at_timestamp();

CREATE TABLE IF NOT EXISTS deleted_face_data_id_mapping (
    "id" SERIAL PRIMARY KEY,
    "faceDataId" INT,
    "accessPointId" INT,
    "accessorId" INT,
    "assignedAt" INT,
    "unAssignedAt" INT,
    "createdAt" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "updatedAt" TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TRIGGER trigger_deleted_face_data_id_mapping_set_updated_at_timestamp BEFORE UPDATE ON deleted_face_data_id_mapping FOR EACH ROW EXECUTE FUNCTION set_updated_at_timestamp();

COMMIT;