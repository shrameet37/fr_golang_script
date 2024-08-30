package models

import (
	"time"

	"github.com/google/uuid"
)

type OrganisationFaceIdsSchema struct {
	FaceDataId                 int       `json:"faceDataId,omitempty"`
	FaceData                   string    `json:"faceData,omitempty"`
	AccessorId                 int       `json:"accessorId"`
	OrganisationId             int       `json:"organisationId"`
	UserName                   string    `json:"userName,omitempty"`
	PendingUnassignedOnDevices bool      `json:"PendingUnassignedOnDevices"`
	CreatedAt                  time.Time `json:"createdAt,omitempty" binding:"omitempty"`
	UpdatedAt                  time.Time `json:"updatedAt,omitempty"`
}

type IotEngineCommandMsg struct {
	Version        int       `json:"version"`            // 1
	SrcAppID       uint8     `json:"srcAppId"`           // 0x88
	DestAppID      uint8     `json:"dstAppId"`           // 0x04
	IsLiveMsg      bool      `json:"isLiveMsg"`          // false
	Target         string    `json:"target"`             // device
	TargetSerialNo string    `json:"targetSerialNumber"` //
	MsgType        string    `json:"msgType"`
	MsgTypeVer     uint8     `json:"msgTypeVer"` // add_fingerprint_msg
	UserId         uint32    `json:"userId"`     // specify here users accessor Id
	UserName       string    `json:"userName"`   // max 24 characters
	KeyId          uint32    `json:"keyId"`
	ValidationKey  string    `json:"validationKey"`
	MessageId      uuid.UUID `json:"messageId"`
}

type AddAccessorFacedataPathParams struct {
	OrganisationId int `uri:"organisationId" binding:"required,min=1"`
	AccessorId     int `uri:"accessorId" binding:"required,min=1"`
}

type AddAccessorFacedataRequestBody struct {
	FaceData string `json:"faceData"`
	UserName string `json:"userName"`
}

type AddAccessorFacedataResponse struct {
	Type    string                             `json:"type"`
	Message AddAccessorFacedataResponseMessage `json:"message"`
}

type AddAccessorFacedataResponseMessage struct {
	FaceDataId int `json:"faceDataId"`
}

type UpdateAccessorFacedataPathParams struct {
	OrganisationId int `uri:"organisationId" binding:"required,min=1"`
	AccessorId     int `uri:"accessorId" binding:"required,min=1"`
}

type UpdateAccessorFacedataRequestBody struct {
	FaceData string `json:"faceData"`
	UserName string `json:"userName"`
}

type UpdateAccessorFacedataResponse struct {
	Type    string                                `json:"type"`
	Message UpdateAccessorFacedataResponseMessage `json:"message"`
}

type UpdateAccessorFacedataResponseMessage struct {
	FaceDataId int `json:"faceDataId"`
}

type DeleteAccessorFacedataPathParams struct {
	OrganisationId int `uri:"organisationId" binding:"required,min=1"`
	AccessorId     int `uri:"accessorId" binding:"required,min=1"`
}

type DeleteAccessorFacedataResponse struct {
	Type string `json:"type"`
}
type FaceDataIdMappingSchema struct {
	Id            int       `json:"id,omitempty"`
	FaceDataId    int       `json:"faceDataId,omitempty"`
	AccessPointId int       `json:"accessPointId"`
	AccessorId    int       `json:"accessorId"`
	AssignedAt    int       `json:"assignedAt"`
	CreatedAt     time.Time `json:"createdAt,omitempty" binding:"omitempty"`
	UpdatedAt     time.Time `json:"updatedAt,omitempty"`
}

type ResponseSuccess struct {
	Type string `json:"type"`
}
