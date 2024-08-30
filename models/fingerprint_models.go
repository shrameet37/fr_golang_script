package models

import (
	"time"

	"github.com/google/uuid"
)

type AcaasTransactionsSchema struct {
	Id                      int       `json:"id,omitempty"`
	DevicePermissionTableId int       `json:"devicePermissionTableId"`
	TransactionType         string    `json:"transactionType"`
	UpdatedTransactionType  string    `json:"-"`
	RequestId               uuid.UUID `json:"requestId"`
	SubRequestId            int       `json:"subRequestId,omitempty"`
	CreatedAt               time.Time `json:"createdAt,omitempty" binding:"omitempty"`
	UpdatedAt               time.Time `json:"updatedAt,omitempty"`
}

type IotTransactionsSchema struct {
	Id                       int       `json:"id,omitempty"`
	MessageId                uuid.UUID `json:"messageId,omitempty"`
	SerialNumber             string    `json:"serialNumber,omitempty"`
	DevicePermissionsTableId int       `json:"devicePermissionsTableId,omitempty"`
	KeyId                    int       `json:"keyId"`
	ResponseReceived         int       `json:"responseReceived,omitempty"`
	AckType                  int       `json:"ackType,omitempty"`
	TransactionType          string    `json:"transactionType,omitempty"`
	CloudTime                int       `json:"cloudTime,omitempty"`
	GatewayTime              int       `json:"gatewayTime,omitempty"`
	RoundTripTime            float64   `json:"roundTripTime"`
	CreatedAt                time.Time `json:"createdAt,omitempty" binding:"omitempty"`
	UpdatedAt                time.Time `json:"updatedAt,omitempty"`
}

type AccessorPermissionsSchema struct {
	Id             int       `json:"id,omitempty"`
	AccessorId     int       `json:"accessorId"`
	OrganisationId int       `json:"organisationId"`
	TemplateId     int       `json:"templateId,omitempty"`
	CreatedAt      time.Time `json:"createdAt,omitempty" binding:"omitempty"`
	UpdatedAt      time.Time `json:"updatedAt,omitempty"`
}

type DevicePermissionsSchema struct {
	Id              int       `json:"id,omitempty"`
	SerialNumber    string    `json:"serialNumber,omitempty"`
	CredentialId    uint32    `json:"credentialId" binding:"required,min=0"`
	FaceId          int       `json:"faceDataId,omitempty"`
	AccessPointId   int       `json:"accessPointId"`
	ChannelNo       int       `json:"channelNo"`
	AccessorId      int       `json:"accessorId"`
	OrganisationId  int       `json:"organisationId"`
	UpdatedOnDevice int       `json:"updatedOnDevice"`
	ToDelete        int       `json:"toDelete"`
	AssignedAt      int       `json:"assignedAt"`
	CreatedAt       time.Time `json:"createdAt,omitempty" binding:"omitempty"`
	UpdatedAt       time.Time `json:"updatedAt,omitempty"`
	DeletedAt       time.Time `json:"deletedAt"`
}

type TruncatedDevicePermission struct {
	Id                  int   `json:"id"`
	CredentialId        int32 `json:"credentialId"`
	OrgId               int   `json:"organisationId"`
	AccessorIdIfPresent int   `json:"accessorIdIfPresent"`
	AccessPointId       int   `json:"accessPointId"`
	ToDelete            int   `json:"toDelete"`
	UpdatedOnDevice     int   `json:"updatedOnDevice"`
}

type DeleteDevicePermissionsSchema struct {
	Id                              int       `json:"id,omitempty"`
	DeletedDevicePermissionsTableId int       `json:"deletedDevicePermissionsTableId"`
	SerialNumber                    string    `json:"serialNumber,omitempty"`
	FaceId                          int       `json:"faceId,omitempty"`
	CredentialId                    uint32    `json:"credentialId" binding:"required,min=0"`
	AccessorId                      int       `json:"accessorId"`
	AccessPointId                   int       `json:"accessPointId"`
	ChannelNo                       int       `json:"channelNo"`
	OrganisationId                  int       `json:"organisationId"`
	AssignedAt                      int       `json:"assignedAt"`
	UnAssignedAt                    int       `json:"unAssignedAt"`
	CreatedAt                       time.Time `json:"createdAt,omitempty" binding:"omitempty"`
	UpdatedAt                       time.Time `json:"updatedAt,omitempty"`
}

type AccessPointsSchema struct {
	Id             int       `json:"id,omitempty"`
	AccessPointId  int       `json:"accessPointId"`
	OrganisationId int       `json:"organisationId"`
	SiteId         int       `json:"siteId"`
	Configuration  int       `json:"configuration"`
	CreatedAt      time.Time `json:"createdAt,omitempty" binding:"omitempty"`
	UpdatedAt      time.Time `json:"updatedAt,omitempty"`
}

type DevicesSchema struct {
	Id             int       `json:"id,omitempty"`
	SerialNumber   string    `json:"serialNumber,omitempty"`
	DeviceType     int       `json:"deviceType"`
	OrganisationId int       `json:"organisationId"`
	CreatedAt      time.Time `json:"createdAt,omitempty" binding:"omitempty"`
	UpdatedAt      time.Time `json:"updatedAt,omitempty"`
}

type AccessPointDevicesSchema struct {
	Id            int       `json:"id,omitempty"`
	SerialNumber  string    `json:"serialNumber,omitempty"`
	AccessPointId int       `json:"accessPointId"`
	CreatedAt     time.Time `json:"createdAt,omitempty" binding:"omitempty"`
	UpdatedAt     time.Time `json:"updatedAt,omitempty"`
}

type OrganisationAccessorsSchema struct {
	Id             int       `json:"id,omitempty"`
	AccessorId     int       `json:"accessorId"`
	OrganisationId int       `json:"organisationId"`
	CredentialId   uint32    `json:"credentialId" binding:"required,min=0"`
	CreatedAt      time.Time `json:"createdAt,omitempty" binding:"omitempty"`
	UpdatedAt      time.Time `json:"updatedAt,omitempty"`
}
