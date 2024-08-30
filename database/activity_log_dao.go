package database

import (
	"context"
	"face_management/logger"
	"face_management/models"
	"face_management/utils"
	"fmt"

	"github.com/google/uuid"
)

type alDb struct{}

type alDbInterface interface {
	SaveKafkaActivityLog(queueId int, requestId uuid.UUID, data string) *models.ApplicationError
	SaveApiActivityLog(credentialId, orgId, accessorIdIfPresent int, credentialTarget string, credentialOperation string, additionalInfoId int) *models.ApplicationError
	SaveApiActivityLogAdditionalInfo(additionalInfo string) (int, *models.ApplicationError)
}

var AlDb alDbInterface

const (
	ACTIVITY_LOG_DAO_BYTE3_VALUE = 0x003d0000 //3 = dao, d = activity log
)

var (
	ACTIVITY_LOG_DAO_ERROR_BASE_CODE int
)

func init() {
	AlDb = &alDb{}
	ACTIVITY_LOG_DAO_ERROR_BASE_CODE = utils.ServiceErrorBaseCode + ACTIVITY_LOG_DAO_BYTE3_VALUE
}

func (d *alDb) SaveKafkaActivityLog(queueId int, requestId uuid.UUID, data string) *models.ApplicationError {

	sqlStatement := `insert into activity_logs_kafka ("queueId","requestId","data") values ($1,$2,$3)`

	_, err := dbPool.Exec(context.Background(), sqlStatement, queueId, requestId, data)

	if err != nil {
		errMsg := fmt.Sprintf("SaveKafkaActivityLog:Could not write to kafka activity logs database.Error:%s!", err.Error())
		logger.Log.Error(errMsg)
		appError := utils.RenderAppError(ACTIVITY_LOG_DAO_ERROR_BASE_CODE+1, errMsg)
		return appError
	}

	return nil

}

func (d *alDb) SaveApiActivityLog(credentialId, orgId, accessorIdIfPresent int, isAssignedTo string, operation string, additionalInfoId int) *models.ApplicationError {

	sqlStatement := `insert into activity_logs ("credentialId","accessorIdIfPresent","isAssignedTo","operation","additionalInfoId") values ($1,$2,$3,$4,$5)`
	_, err := dbPool.Exec(context.Background(), sqlStatement, credentialId, orgId, accessorIdIfPresent, isAssignedTo, operation, additionalInfoId)

	if err != nil {
		errMsg := fmt.Sprintf("SaveApiActivityLog:Could not write to api kafka activity logs database table.Error:%s!", err.Error())
		logger.Log.Error(errMsg)
		appError := utils.RenderAppError(ACTIVITY_LOG_DAO_ERROR_BASE_CODE+2, errMsg)
		return appError
	}

	return nil

}

func (d *alDb) SaveApiActivityLogAdditionalInfo(additionalInfo string) (int, *models.ApplicationError) {

	var id int

	sqlStatement := `insert into activity_logs_additional_info ("additionalInfo") values ($1) returning id`
	err := dbPool.QueryRow(context.Background(), sqlStatement, additionalInfo).Scan(&id)

	if err != nil {
		errMsg := fmt.Sprintf("SaveApiActivityLogAdditionalInfo:Could not write to api kafka activity logs additional info database table.Error:%s!", err.Error())
		logger.Log.Error(errMsg)
		appError := utils.RenderAppError(ACTIVITY_LOG_DAO_ERROR_BASE_CODE+3, errMsg)
		return 0, appError
	}

	return id, nil

}
