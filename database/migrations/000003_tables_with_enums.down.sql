BEGIN;

DROP TABLE IF EXISTS acaas_transactions;
DROP TRIGGER IF EXISTS trigger_acaas_transactions_set_updated_at_timestamp ON acaas_transactions;

DROP TABLE IF EXISTS iot_transactions;
DROP TRIGGER IF EXISTS trigger_iot_transactions_set_updated_at_timestamp ON iot_transactions;

COMMIT;
