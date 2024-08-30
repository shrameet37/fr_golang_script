package models

type SuccessResponse struct {
	Type    string      `json:"type"`
	Message interface{} `json:"message"`
}

const (
	IsAssignedToAccessorOnAccessPoint = 3
	IsAssignedToAccessPoint           = 2
	IsAssignedToAccessor              = 1
	IsAssignedToFree                  = 0
)

const (
	AckNotRecieved             = 0
	AckRecievedFromIot         = 1
	AckRecievedFromSdk         = 2
	AckRecievedFromIotAfterSdk = 3
)

type KeypadIdAccessPointId struct {
	AccessPointId int `json:"accessPointId"`
	KeypadId      int `json:"keypadId"`
}
