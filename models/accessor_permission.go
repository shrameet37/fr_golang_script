package models

import (
	"time"

	"github.com/google/uuid"
)

type AccessorPermission struct {
	Id            int       `json:"id,omitempty"`
	AccessorId    int       `json:"accessorId,omitempty"`
	AccessPointId int       `json:"accessPointId,omitempty"`
	OrgId         int       `json:"orgId,omitempty"`
	CreatedAt     time.Time `json:"createdAt,omitempty"`
}

type AccessorPermissionStatus struct {
	Permission int
	InStore    bool
	InRequest  bool
	InBoth     bool
}

type GetPermissionsResponse struct {
	Count       int
	Permissions []AccessorPermission
}

type PendingPermissionsOnDeviceResponse struct {
	Count              int
	PendingPermissions []IotTransactionsSchema
}

type AddAccessorToOrg struct {
	RequestId    uuid.UUID
	OrgId        int
	AccessorId   int
	CredentialId uint32
}
