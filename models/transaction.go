package models

import (
	"github.com/google/uuid"
)

const (
	AssignFaceTransactionType   string = "assign_face"
	UnassignFaceTransactionType string = "unassign_face"
	UpdateFaceTransactionType   string = "update_face"

	AddPermissionTransactionType    string = "add_permission"
	RemovePermissionTransactionType string = "remove_permission"
)

type RequestDeviceMessage struct {
	Operation          int
	RequestId          uuid.UUID
	SubRequestId       int
	IsController       bool
	CredentialId       uint32
	DeviceSerialNumber string
	DevicePermissionId int
	ChannelNo          int
	Configuration      int
	DeviceType         int
	UserId             uint32
	UserName           string
	KeyId              uint32
	ValidationKey      string
}

type RequestIotTransaction struct {
	MessageId          uuid.UUID
	DeviceSerialNumber string
	DevicePermissionId int
	TransactionType    string
	KeyId              int
}

type AckEvent struct {
	StatusCode   int `json:"statusCode"`
	SubRequestId int `json:"subRequestId,omitempty"`
}

type Event struct {
	MessageVersion int         `json:"messageVersion"`
	MessageData    MessageData `json:"messageData"`
}

type MessageData struct {
	RequestId   uuid.UUID `json:"requestId"`
	MessageType string    `json:"messageType"`
	DataVersion int       `json:"dataVersion"`
	Service     int       `json:"service"`
	Data        AckEvent  `json:"data"`
}

type RequestMsgToAccessPoint struct {
	CredentialId  uint32
	AccessorId    int
	OrgId         int
	RequestId     uuid.UUID
	SubRequestId  int
	Operation     int
	AccessPointId int
	FaceDataId    int
}

type RequestMsgToAccessPoints struct {
	CredentialId uint32
	Uuid         int32
	AccessorId   int
	OrgId        int
	RequestId    uuid.UUID
	SubRequestId int
	Operation    int
	Permissions  []AccessPointInfo
}

type RequestAckEvent struct {
	RequestId       uuid.UUID `json:"requestId"`
	SubRequestId    int       `json:"subRequestId"`
	KafkaMessageKey int       `json:"kafkaMessageKey"`
	StatusCode      int       `json:"statusCode"`
}
