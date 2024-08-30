package methods

import (
	"encoding/json"
	"face_management/clients"
	"face_management/models"
)

func SendMessageToEventAgrregatorTopic(message interface{}, key string) *models.ApplicationError {
	txByteArray, err := json.Marshal(message)
	if err != nil {
	} else {

		sqsMessage := string(txByteArray)

		clients.ToSqsChEventAggregatorQueue <- sqsMessage
	}
	return nil

}
