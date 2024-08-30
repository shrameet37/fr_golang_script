package models

import (
	"time"
)

type DevicePermissionStatusResponse struct {
	Status Status `json:"status"`
}

type Status struct {
	AckRecievedFromDevice   bool      `json:"ackReceivedFromDevice"`
	SentToDeviceAt          time.Time `json:"sentToDeviceAt"`
	RecievedAckFromDeviceAt time.Time `json:"receivedAckFromDeviceAt"`
}
