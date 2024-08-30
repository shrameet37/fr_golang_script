BEGIN;

DROP TABLE IF EXISTS face_errors;
DROP TRIGGER IF EXISTS trigger_face_errors_set_updated_at_timestamp ON face_errors;

DROP TABLE IF EXISTS activity_logs_kafka;
DROP TRIGGER IF EXISTS trigger_activity_logs_kafka_set_updated_at_timestamp ON activity_logs_kafka;

DROP TABLE IF EXISTS kafka_topic_dropped_messages;
DROP TRIGGER IF EXISTS trigger_kafka_topic_dropped_messages_set_updated_at_timestamp ON kafka_topic_dropped_messages;

ALTER TABLE deleted_device_permissions DROP COLUMN  IF EXISTS "deletedAt";

COMMIT;
