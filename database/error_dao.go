package database

import (
	"context"
	"face_management/logger"
	"face_management/models"
	"face_management/utils"
	"fmt"
)

type errorDb struct{}

type errorDbInterface interface {
	SaveDroppedMessage(droppedAckEvent models.DroppedMessage) *models.ApplicationError
	ProcessErrorMessages(priority int, errorMessage string, moreInfo string) *models.ApplicationError
}

var ErDb errorDbInterface

const (
	ERROR_DAO_BYTE3_VALUE = 0x00340000 //3 = dao, 4 = transaction
)

var (
	ERROR_DAO_ERROR_BASE_CODE int
)

func init() {
	ErDb = &errorDb{}
	ERROR_DAO_ERROR_BASE_CODE = utils.ServiceErrorBaseCode + ERROR_DAO_BYTE3_VALUE

}

func (d *errorDb) SaveDroppedMessage(droppedMessage models.DroppedMessage) *models.ApplicationError {

	sqlStatement := `insert into kafka_topic_dropped_messages ("topicName","errorType","kafkaMessage") values ($1,$2,$3)`
	_, err := dbPool.Exec(context.Background(), sqlStatement, droppedMessage.TopicName, droppedMessage.ErrorType, droppedMessage.KafkaMessage)

	if err != nil {
		logger.Log.Error(err.Error())
		appError := utils.RenderAppError(ERROR_DAO_ERROR_BASE_CODE+1, fmt.Sprintf("Could not write to topic dropped message database table.Error:%s", err))
		return appError
	}

	return nil

}

func (d *errorDb) ProcessErrorMessages(priority int, errorMessage string, additionalInfo string) *models.ApplicationError {

	sqlStatement := `insert into face_errors ("priority","errorMessage","additionalInfo") values ($1,$2,$3)`
	_, err := dbPool.Exec(context.Background(), sqlStatement, priority, errorMessage, additionalInfo)

	if err != nil {
		logger.Log.Error(err.Error())
		appError := utils.RenderAppError(ERROR_DAO_ERROR_BASE_CODE+2, fmt.Sprintf("Could not write to card_errors database table.Error:%s", err))
		return appError
	}

	return nil

}
