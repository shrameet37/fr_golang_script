BEGIN;

CREATE FUNCTION set_updated_at_timestamp()
RETURNS TRIGGER AS $$
BEGIN
  NEW."updatedAt" = NOW();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TABLE IF NOT EXISTS organisation_face_data_ids (
    "faceDataId" SERIAL PRIMARY KEY,
    "faceData" BYTEA,
    "organisationId" INT,
    "accessorId" INT,
    "userName" VARCHAR(64),
    "pendingUnassignedOnDevices" BOOLEAN DEFAULT false,
    "createdAt" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "updatedAt" TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE UNIQUE INDEX IF NOT EXISTS "organisation_face_data_ids_accessorId_organisationId_index" ON organisation_face_data_ids("accessorId", "organisationId");

CREATE TRIGGER trigger_organisation_face_data_ids_set_updated_at_timestamp BEFORE UPDATE ON organisation_face_data_ids FOR EACH ROW EXECUTE FUNCTION set_updated_at_timestamp();

COMMIT;