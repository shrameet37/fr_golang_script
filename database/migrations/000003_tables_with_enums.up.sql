BEGIN;


CREATE TABLE IF NOT EXISTS acaas_transactions (
    "id" SERIAL PRIMARY KEY,
    "devicePermissionsTableId" INT,
    "transactionType" TEXT,
    "requestId" UUID,
    "subRequestId" INT,
    "createdAt" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "updatedAt" TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TRIGGER trigger_acaas_transactions_set_updated_at_timestamp BEFORE UPDATE ON acaas_transactions FOR EACH ROW EXECUTE FUNCTION set_updated_at_timestamp();

CREATE TABLE IF NOT EXISTS iot_transactions (
    "id" SERIAL PRIMARY KEY,
    "messageId" UUID,
    "serialNumber" VARCHAR(20),
    "devicePermissionsTableId" INT,
    "responseReceived" INT,
    "ackType" INT,
    "transactionType" TEXT,
    "keyId" INT,
    "gatewayTime" BIGINT NOT NULL DEFAULT 0,
    "cloudTime" BIGINT NOT NULL DEFAULT 0,
    "roundTripTime" REAL NOT NULL DEFAULT 0,
    "createdAt" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "updatedAt" TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TRIGGER trigger_iot_transactions_set_updated_at_timestamp BEFORE UPDATE ON iot_transactions FOR EACH ROW EXECUTE FUNCTION set_updated_at_timestamp();


COMMIT;