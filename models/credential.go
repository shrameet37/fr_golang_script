package models

import (
	"github.com/google/uuid"
)

type CredentialAccessorAssignmentRequest struct {
	RequestId    uuid.UUID         `json:"requestId,omitempty"`
	OrgId        int               `json:"orgId,omitempty"`
	AccessorId   int               `json:"accessorId,omitempty"`
	CredentialId uint32            `json:"credentialId,omitempty"`
	Permissions  []AccessPointInfo `json:"permissions,omitempty"`
}

type CredentialAccessorUpdatePermissionRequest struct {
	RequestId           uuid.UUID
	OrgId               int
	AccessorId          int
	PermissionsToAdd    []AccessPointInfo
	PermissionsToRemove []AccessPointInfo
}

type CredentialAccessorPermissions struct {
	CredentialId  int
	AccessorId    int
	AccessPointId int
	SubRequestId  int
}

type CredentialResponse struct {
	Message string `json:"message"`
}

type CredentialAccessPointUpdatePermissionRequest struct {
	RequestId           uuid.UUID
	OrgId               int
	AccessPointId       int
	PermissionsToAdd    []AccessorInfo
	PermissionsToRemove []AccessorInfo
}
