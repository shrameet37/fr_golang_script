package models

import "github.com/google/uuid"

type AccessorPermissions struct {
	AccessorId  int
	OrgId       int
	Operation   int
	Permissions []AccessPointInfo
}

type AccessPointPermissions struct {
	AccessPointId int
	OrgId         int
	Operation     int
	Permissions   []AccessorInfo
}

type AccessPointInfo struct {
	AccessPointId int
	SubRequestId  int
}

type AccessorInfo struct {
	AccessorId   int
	SubRequestId int
}

type UpdateCredential struct {
	RequestId    uuid.UUID
	OrgId        int
	AccessorId   int
	CredentialId uint32
}

type RemoveAccessorFromOrg struct {
	RequestId  uuid.UUID
	OrgId      int
	AccessorId int
}

type PendingPermissionFaceStatusRequestBody struct {
	Operation      string `json:"operation" binding:"required"`
	IsFaceAssigned bool   `json:"isFaceAssigned" binding:"required"`
	AccessPoints   []int  `json:"accessPoints" binding:"required"`
}

type PendingPermissionFaceStatusPathParams struct {
	OrganisationId int `uri:"organisationId"`
	AccessorId     int `uri:"accessorId"`
}

type PendingPermissionFaceStatusResponse struct {
	Type    string                                     `json:"type"`
	Message PendingPermissionFaceStatusResponseMessage `json:"message"`
}

type PendingPermissionFaceStatusResponseMessage struct {
	RequestId   uuid.UUID                                `json:"requestId"`
	Permissions []PendingPermissionFaceStatusPermissions `json:"permissions"`
}
type PendingPermissionFaceStatusPermissions struct {
	AccessPointId int `json:"accessPointId"`
	SubRequestId  int `json:"subRequestId"`
}
