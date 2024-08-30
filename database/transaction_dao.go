package database

import (
	"context"
	"face_management/logger"
	"face_management/models"
	"face_management/utils"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/lib/pq"
)

type txDb struct{}

type txDbInterface interface {
	CreateIotTxLog(data models.IotTransactionsSchema) (appError *models.ApplicationError)

	GetLatestIotTxn(serialNumber string) (iotRequest string, devicePermissionId int, appError *models.ApplicationError)
	GetIotTxLog(messageId uuid.UUID) (exist bool, data *models.IotTransactionsSchema, appError *models.ApplicationError)
	GetAllPendingIotTxns(ids []int) (iotTxns []models.IotTransactionsSchema, apiError *models.ApplicationError)

	UpdateIotTxLog(data models.IotTransactionsSchema) (appError *models.ApplicationError)

	DoesAcaasTransactionAssignOrAddExist(devicePermissionsTableId int) (exists bool, previousTransactionType string, appError *models.ApplicationError)
	DoesAcaasTransactionUnassignOrRemoveExist(devicePermissionsTableId int) (exists bool, previousTransactionType string, appError *models.ApplicationError)
	CreateAcaasTransaction(acaasTxn models.AcaasTransactionsSchema) *models.ApplicationError
	UpdateAcaasTransaction(acaasTxn models.AcaasTransactionsSchema) *models.ApplicationError
	UpdateAcaasTransactionWithoutTransactionType(acaasTxn models.AcaasTransactionsSchema) *models.ApplicationError
	UpdateAcaasTransactionOnlyTransactionType(acaasTxn models.AcaasTransactionsSchema) *models.ApplicationError
	DeleteAcaasTransaction(devicePermissionsTableId int, transactionType string) *models.ApplicationError

	GetDetailsOfAssignOrAddAcaasTxn(devicePermissionsTableId int) (exist bool, acaasTxn *models.AcaasTransactionsSchema, appError *models.ApplicationError)
	GetDetailsOfUnassignOrRemoveAcaasTxn(devicePermissionsTableId int) (acaasTxn *models.AcaasTransactionsSchema, appError *models.ApplicationError)
}

var TxDb txDbInterface

const (
	TRANSACTION_BYTE3_VALUE = 0x003c0000 //3 = dao, 6 = transaction
)

var (
	TRANSACTION_ERROR_BASE_CODE int
)

func init() {
	TxDb = &txDb{}
	TRANSACTION_ERROR_BASE_CODE = utils.ServiceErrorBaseCode + TRANSACTION_BYTE3_VALUE

}

func (d *txDb) CreateIotTxLog(data models.IotTransactionsSchema) (appError *models.ApplicationError) {

	sqlStatement := `INSERT INTO iot_transactions ("keyId","messageId", "serialNumber", "devicePermissionsTableId", "transactionType") VALUES ($1, $2, $3, $4, $5) RETURNING id`

	_, err := dbPool.Exec(context.Background(), sqlStatement, data.KeyId, data.MessageId, data.SerialNumber, data.DevicePermissionsTableId, data.TransactionType)

	if err != nil {

		logger.Log.Error(err.Error())
		appError := utils.RenderAppError(TRANSACTION_ERROR_BASE_CODE+1, fmt.Sprintf("Could not create iot transaction log. Error: %s", err.Error()))
		return appError
	}

	return nil
}

func (d *txDb) GetIotTxLog(messageId uuid.UUID) (exist bool, iotTxn *models.IotTransactionsSchema, appError *models.ApplicationError) {

	iotTxn = &models.IotTransactionsSchema{}

	sqlStatement := `SELECT
							it.id,
							it."devicePermissionsTableId",
							it."transactionType",
							it."keyId",
							it."createdAt"
						FROM
							iot_transactions it
						WHERE
							"messageId" = $1;`

	err := dbPool.QueryRow(context.Background(), sqlStatement, messageId).Scan(&iotTxn.Id, &iotTxn.DevicePermissionsTableId, &iotTxn.TransactionType, &iotTxn.KeyId, &iotTxn.CreatedAt)

	if err != nil {

		if err == pgx.ErrNoRows {
			return false, nil, nil
		}

		logger.Log.Error(err.Error())
		appError := utils.RenderAppError(TRANSACTION_ERROR_BASE_CODE+2, fmt.Sprintf("Could not get iot transaction log. Error: %s", err.Error()))
		return false, nil, appError
	}

	return true, iotTxn, nil
}

func (d *txDb) UpdateIotTxLog(data models.IotTransactionsSchema) (appError *models.ApplicationError) {

	sqlStatement := `UPDATE iot_transactions SET "responseReceived"=$1, "ackType"=$2, "gatewayTime"=$3, "cloudTime"=$4 , "roundTripTime"=$5 WHERE "id"=$6;`

	_, err := dbPool.Exec(context.Background(), sqlStatement, data.ResponseReceived, data.AckType, data.GatewayTime, data.CloudTime, data.RoundTripTime, data.Id)

	if err != nil {

		logger.Log.Error(err.Error())
		appError := utils.RenderAppError(TRANSACTION_ERROR_BASE_CODE+3, fmt.Sprintf("Could not update iot transaction log. Error: %s", err.Error()))
		return appError
	}

	return nil
}

func (d *txDb) GetAllPendingIotTxns(ids []int) (iotTxns []models.IotTransactionsSchema, apiError *models.ApplicationError) {

	if len(ids) == 0 {
		return nil, nil
	}

	sqlStatement := `SELECT "messageId" ,"serialNumber" ,"iotRespReceived" ,"ackType" , "roundTripTime", "createdAt" FROM iot_transactions WHERE "devicePermissionsTableId"= ANY($1) AND "iotRespReceived"=false AND "ackType" != 1;`

	rows, err := dbPool.Query(context.Background(), sqlStatement, pq.Array(ids))
	if err != nil {
		logger.Log.Error(err.Error())
		appError := utils.RenderAppError(TRANSACTION_ERROR_BASE_CODE+4, fmt.Sprintf("Could not verify if iot txns exist. Error: %s", err.Error()))
		return nil, appError
	}

	defer rows.Close()

	for rows.Next() {
		var iotTxn models.IotTransactionsSchema
		if err := rows.Scan(&iotTxn.MessageId, &iotTxn.SerialNumber, &iotTxn.ResponseReceived, &iotTxn.AckType, &iotTxn.RoundTripTime, &iotTxn.CreatedAt); err != nil {
			logger.Log.Error(err.Error())
			appError := utils.RenderAppError(TRANSACTION_ERROR_BASE_CODE+5, fmt.Sprintf("Error: %s", err.Error()))
			return nil, appError
		}
		iotTxns = append(iotTxns, iotTxn)
	}

	return iotTxns, nil
}

func (d *txDb) GetLatestIotTxn(serialNumber string) (iotRequest string, devicePermissionId int, appError *models.ApplicationError) {

	sqlStatement := `SELECT "transactionType", "devicePermissionsTableId" FROM iot_transactions WHERE "serialNumber" = $1 ORDER BY "createdAt" DESC;`

	err := dbPool.QueryRow(context.Background(), sqlStatement, serialNumber).Scan(&iotRequest, &devicePermissionId)

	if err != nil {

		logger.Log.Error(err.Error())
		appError := utils.RenderAppError(TRANSACTION_ERROR_BASE_CODE+6, fmt.Sprintf("Error: %s", err.Error()))
		return iotRequest, devicePermissionId, appError
	}

	return iotRequest, devicePermissionId, nil
}

func (d *txDb) DoesAcaasTransactionAssignOrAddExist(devicePermissionsTableId int) (exists bool, previousTransactionType string, appError *models.ApplicationError) {

	sqlStatement := `SELECT "transactionType" FROM acaas_transactions WHERE "transactionType" IN ('assign_card' , 'add_permission') AND "devicePermissionsTableId" = $1;`

	err := dbPool.QueryRow(context.Background(), sqlStatement, devicePermissionsTableId).Scan(&previousTransactionType)

	if err != nil {

		if err == pgx.ErrNoRows {

			return false, previousTransactionType, nil
		}

		logger.Log.Error(err.Error())
		appError := utils.RenderAppError(TRANSACTION_ERROR_BASE_CODE+6, fmt.Sprintf("Error: %s", err.Error()))
		return false, previousTransactionType, appError
	}

	return true, previousTransactionType, nil
}

func (d *txDb) DoesAcaasTransactionUnassignOrRemoveExist(devicePermissionsTableId int) (exists bool, previousTransactionType string, appError *models.ApplicationError) {

	sqlStatement := `SELECT "transactionType" FROM acaas_transactions WHERE "transactionType" IN ('unassign_face' , 'remove_permission') 
							AND "devicePermissionsTableId" = $1;`

	err := dbPool.QueryRow(context.Background(), sqlStatement, devicePermissionsTableId).Scan(&previousTransactionType)

	if err != nil {

		if err == pgx.ErrNoRows {

			return false, previousTransactionType, nil
		}

		logger.Log.Error(err.Error())
		appError := utils.RenderAppError(TRANSACTION_ERROR_BASE_CODE+6, fmt.Sprintf("Error: %s", err.Error()))
		return false, previousTransactionType, appError
	}

	return true, previousTransactionType, nil
}

func (d *txDb) CreateAcaasTransaction(acaasTxn models.AcaasTransactionsSchema) *models.ApplicationError {

	sqlStatement := `INSERT INTO acaas_transactions ("transactionType", "devicePermissionsTableId", "requestId", "subRequestId") VALUES ($1, $2, $3, $4);`

	_, err := dbPool.Exec(context.Background(), sqlStatement, acaasTxn.TransactionType, acaasTxn.DevicePermissionTableId, acaasTxn.RequestId, acaasTxn.SubRequestId)

	if err != nil {

		logger.Log.Error(err.Error())
		appError := utils.RenderAppError(TRANSACTION_ERROR_BASE_CODE+6, fmt.Sprintf("Error: %s", err.Error()))
		return appError
	}

	return nil
}

func (d *txDb) UpdateAcaasTransaction(acaasTxn models.AcaasTransactionsSchema) *models.ApplicationError {

	sqlStatement := `UPDATE acaas_transactions SET "transactionType" = $1, "requestId" = $2, "subRequestId" =$3 WHERE "transactionType" = $4 AND "devicePermissionsTableId" = $5;`

	_, err := dbPool.Exec(context.Background(), sqlStatement, acaasTxn.UpdatedTransactionType, acaasTxn.RequestId, acaasTxn.SubRequestId, acaasTxn.TransactionType, acaasTxn.DevicePermissionTableId)

	if err != nil {

		logger.Log.Error(err.Error())
		appError := utils.RenderAppError(TRANSACTION_ERROR_BASE_CODE+6, fmt.Sprintf("Error: %s", err.Error()))
		return appError
	}

	return nil
}

func (d *txDb) UpdateAcaasTransactionWithoutTransactionType(acaasTxn models.AcaasTransactionsSchema) *models.ApplicationError {

	sqlStatement := `UPDATE acaas_transactions SET "requestId" = $1, "subRequestId" =$2 WHERE "transactionType" = $3 AND "devicePermissionsTableId" = $4;`

	_, err := dbPool.Exec(context.Background(), sqlStatement, acaasTxn.RequestId, acaasTxn.SubRequestId, acaasTxn.TransactionType, acaasTxn.DevicePermissionTableId)

	if err != nil {

		logger.Log.Error(err.Error())
		appError := utils.RenderAppError(TRANSACTION_ERROR_BASE_CODE+6, fmt.Sprintf("Error: %s", err.Error()))
		return appError
	}

	return nil
}

func (d *txDb) UpdateAcaasTransactionOnlyTransactionType(acaasTxn models.AcaasTransactionsSchema) *models.ApplicationError {

	sqlStatement := `UPDATE acaas_transactions SET "transactionType" = $1 WHERE "transactionType" = $2 AND "devicePermissionsTableId" = $3;`

	_, err := dbPool.Exec(context.Background(), sqlStatement, acaasTxn.UpdatedTransactionType, acaasTxn.TransactionType, acaasTxn.DevicePermissionTableId)

	if err != nil {

		logger.Log.Error(err.Error())
		appError := utils.RenderAppError(TRANSACTION_ERROR_BASE_CODE+6, fmt.Sprintf("Error: %s", err.Error()))
		return appError
	}

	return nil
}

func (d *txDb) DeleteAcaasTransaction(devicePermissionsTableId int, transactionType string) *models.ApplicationError {

	sqlStatement := `DELETE FROM acaas_transactions WHERE "devicePermissionsTableId" = $1 AND "transactionType" = $2;`

	_, err := dbPool.Exec(context.Background(), sqlStatement, devicePermissionsTableId, transactionType)

	if err != nil {

		logger.Log.Error(err.Error())
		appError := utils.RenderAppError(TRANSACTION_ERROR_BASE_CODE+6, fmt.Sprintf("Error: %s", err.Error()))
		return appError
	}

	return nil
}

func (d *txDb) GetDetailsOfAssignOrAddAcaasTxn(devicePermissionsTableId int) (exist bool, acaasTxn *models.AcaasTransactionsSchema, appError *models.ApplicationError) {

	acaasTxn = &models.AcaasTransactionsSchema{}

	sqlStatement := `SELECT "requestId", "subRequestId" FROM acaas_transactions WHERE "devicePermissionsTableId" = $1 AND "transactionType" IN ('assign_face', 'add_permission');`

	err := dbPool.QueryRow(context.Background(), sqlStatement, devicePermissionsTableId).Scan(&acaasTxn.RequestId, &acaasTxn.SubRequestId)

	if err != nil {

		logger.Log.Error(err.Error())
		appError := utils.RenderAppError(TRANSACTION_ERROR_BASE_CODE+6, fmt.Sprintf("Error: %s", err.Error()))
		return false, nil, appError
	}

	return true, acaasTxn, nil
}

func (d *txDb) GetDetailsOfUnassignOrRemoveAcaasTxn(devicePermissionsTableId int) (acaasTxn *models.AcaasTransactionsSchema, appError *models.ApplicationError) {

	acaasTxn = &models.AcaasTransactionsSchema{}

	sqlStatement := `SELECT "transactionType", "requestId", "subRequestId" FROM acaas_transactions WHERE "devicePermissionsTableId" = $1 AND "transactionType" IN ('unassign_face', 'remove_permission');`

	err := dbPool.QueryRow(context.Background(), sqlStatement, devicePermissionsTableId).Scan(&acaasTxn.TransactionType, &acaasTxn.RequestId, &acaasTxn.SubRequestId)

	if err != nil {

		logger.Log.Error(err.Error())
		appError := utils.RenderAppError(TRANSACTION_ERROR_BASE_CODE+6, fmt.Sprintf("Error: %s", err.Error()))
		return nil, appError
	}

	return acaasTxn, nil
}
