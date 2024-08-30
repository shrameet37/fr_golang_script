BEGIN;


CREATE TABLE IF NOT EXISTS face_errors (
    "id" SERIAL PRIMARY KEY,
    "priority" INT,
    "errorMessage" TEXT,
    "additionalInfo" TEXT,
    "createdAt" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "updatedAt" TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TRIGGER trigger_face_errors_set_updated_at_timestamp BEFORE UPDATE ON face_errors FOR EACH ROW EXECUTE FUNCTION set_updated_at_timestamp();


CREATE TABLE IF NOT EXISTS activity_logs_kafka (
    "id" SERIAL PRIMARY KEY,
    "queueId" INT,
    "requestId" UUID,
    "data" TEXT,
    "createdAt" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "updatedAt" TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TRIGGER trigger_activity_logs_kafka_set_updated_at_timestamp BEFORE UPDATE ON activity_logs_kafka FOR EACH ROW EXECUTE FUNCTION set_updated_at_timestamp();


CREATE TABLE IF NOT EXISTS kafka_topic_dropped_messages (
    "id" SERIAL PRIMARY KEY,
    "topicName" TEXT,
    "errorType" TEXT,
    "kafkaMessage" TEXT,
    "createdAt" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "updatedAt" TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TRIGGER trigger_kafka_topic_dropped_messages_set_updated_at_timestamp BEFORE UPDATE ON kafka_topic_dropped_messages FOR EACH ROW EXECUTE FUNCTION set_updated_at_timestamp();

ALTER TABLE deleted_device_permissions ADD COLUMN "deletedAt" TIMESTAMPTZ;

COMMIT;
