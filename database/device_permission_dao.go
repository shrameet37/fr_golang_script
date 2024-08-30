package database

import (
	"context"
	"face_management/logger"
	"face_management/models"
	"face_management/utils"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/lib/pq"
)

type devPermDb struct{}

type devPermDbInterface interface {
	CreateDevicePermission(devicePermission models.DevicePermissionsSchema) (devicePermissionId int, appError *models.ApplicationError)
	DeleteDevicePermissionsOnAccessPoint(accessPointId, organisationId int) *models.ApplicationError
	GetDevicePermissionId(devicePermission models.DevicePermissionsSchema) (devicePermissionId int, appError *models.ApplicationError)
	SetToDeleteFlagForDevicePermission(devicePermissionId int) *models.ApplicationError
	SetUpdatedOnDeviceFlagForDevicePermission(devicePermissionId int) *models.ApplicationError
	DeleteDevicePermissionFromId(devicePermissionId int) (appError *models.ApplicationError)
	GetDevicePermissionIds(credentialId uint32, serialNumber string) (ids []int, appError *models.ApplicationError)
	GetDevicePermissionsOfCredential(credentialId uint32) (accessPointIds []int, appError *models.ApplicationError)
	GetDevicePermissionIdsOnAccessPoint(credentialId uint32, accessorId int, accessPointIds []int) (ids []int, appError *models.ApplicationError)
	GetLatestDevicePermission(devicePermissionId, channelNo int) (exists bool, devicePermission models.DevicePermissionsSchema, appError *models.ApplicationError)

	GetDevicePermissionDetailsFromId(devicePermissionId int) (devicePermission *models.DevicePermissionsSchema, appError *models.ApplicationError)
	GetAccessPointIdFromId(devicePermissionId int) (accessPointId int, appError *models.ApplicationError)
	GetDevicePermissionFromId(devicePermissionId int) (devicePermission *models.DevicePermissionsSchema, appError *models.ApplicationError)

	CreateDeletedDevicePermission(devicePermission models.DeleteDevicePermissionsSchema) (appError *models.ApplicationError)
	GetDeletedDevicePermission(deletedDevicePermissionId int) (exists bool, deletedDevicePermission models.DevicePermissionsSchema, appError *models.ApplicationError)

	CheckIfAccessorIdAccessPointIdExistForTemplateIdAddPermission(templateId int, accessorId, accessPointId int) (bool, *models.ApplicationError)
	CheckIfAccessorIdAccessPointIdTemplateIdExistForRemovePermission(templateId int, accessorId, accessPointId int) (bool, *models.ApplicationError)
	CheckIfFaceIdPresentInDevicePermission(templateId int) (bool, *models.ApplicationError)
	GetTruncatedDevicePermission(serialNumber string, channelNo int, credentialId uint32, templateId int) (bool, *models.TruncatedDevicePermission, *models.ApplicationError)
	UpdateToDeleteFlag(devicePermissionsTableId, updateValue int) *models.ApplicationError
	UpdateToDeleteAndUpdatedOnDeviceFlag(devicePermissionsTableId, setToDelete, setUpdatedOnDevice int) *models.ApplicationError
	GetDevicePermissionIdForTemplateNumber(templateNumber int, serialNumber string) (int, *models.ApplicationError)
	CreateDevicePermissionWherePermissionIsSynced(devicePermission models.DevicePermissionsSchema) (devicePermissionId int, appError *models.ApplicationError)
	GetTemplateIdFromDevicePermissions(organisationId int, accessorId int, serialNumber string) (exists bool, templateId int, appError *models.ApplicationError)
	GetDevicePermissionDetailsForTemplateId(credentialId uint32, serialNumber string, templateId int) (*models.DevicePermissionsSchema, *models.ApplicationError)
	CheckIfAccessorIdAccessPointIdExistForCredentialIdAddPermission(credentialId uint32, accessorId, accessPointId int) (bool, *models.ApplicationError)
	GetDevicePendingSyncPermissionsfromDeviceSerialNumber(serialNumber string) ([]models.DevicePermissionsSchema, *models.ApplicationError)
	UpdateCredentialIdforAccessorId(organisationId int, accessorId int, credentialId uint32) *models.ApplicationError
	// GetDevicePermissionsForOrganisation(organisationId int) (*[]models.SyncPendingPermDetails, *models.ApplicationError)
	UpdateSubIndexOfTemplateInDevicePermission(subIndex int, templateId int, credentialId uint32, serialNumber string) *models.ApplicationError

	GetDevicePermissionFromAccessorIdAndDeviceSerialNumber(accessorId int, serialNumber string) (*models.DevicePermissionsSchema, *models.ApplicationError)

	ExcelrGetDevicePermissionFromAccessPointId(accessPointId int) ([]models.ExclerDevicePermission, *models.ApplicationError)
}

var DevPermDb devPermDbInterface

const (
	DEVICE_PERMISSION_DAO_BYTE3_VALUE = 0x00350000 //3 = dao, 5 = device permission
)

var (
	DEVICE_PERMISSION_DAO_ERROR_BASE_CODE int
)

func init() {
	DevPermDb = &devPermDb{}
	DEVICE_PERMISSION_DAO_ERROR_BASE_CODE = utils.ServiceErrorBaseCode + DEVICE_PERMISSION_DAO_BYTE3_VALUE

}

func (d *devPermDb) CreateDevicePermission(devicePermission models.DevicePermissionsSchema) (devicePermissionId int, appError *models.ApplicationError) {

	log.Println("CP-F7")
	sqlStatement := `INSERT INTO device_permissions ("accessorId", "accessPointId", "channelNo", "toDelete", "updatedOnDevice", "organisationId","serialNumber","credentialId", "faceDataId", "assignedAt") VALUES ($1, $2, $3, $4, $5,$6 ,$7, $8, $9, $10) RETURNING id `

	err := dbPool.QueryRow(context.Background(), sqlStatement, devicePermission.AccessorId, devicePermission.AccessPointId, devicePermission.ChannelNo, 0, devicePermission.UpdatedOnDevice, devicePermission.OrganisationId, devicePermission.SerialNumber, devicePermission.CredentialId, devicePermission.FaceId, devicePermission.AssignedAt).Scan(&devicePermissionId)
	if err != nil {
		errMsg := fmt.Sprintf("CreateDevicePermission:Could not create device permission.AccessorId:%d.CredentialId:%d.AccessPointId:%d.SerialNumber:%s.Error:%s!", devicePermission.AccessorId, devicePermission.CredentialId, devicePermission.AccessPointId, devicePermission.SerialNumber, err.Error())
		logger.Log.Error(errMsg)
		appError := utils.RenderAppError(DEVICE_PERMISSION_DAO_ERROR_BASE_CODE+1, errMsg)
		return 0, appError
	}

	return devicePermissionId, nil
}

func (d *devPermDb) CreateDevicePermissions(permissions []models.DevicePermissionsSchema) (devicePermissionIds []int, appError *models.ApplicationError) {
	batch := &pgx.Batch{}

	sqlStatement := `INSERT INTO device_permissions ("accessorId", "accessPointId", "channelNo", "toDelete", "updatedOnDevice", "organisationId","serialNumber","credentialId") VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`

	for _, permission := range permissions {
		batch.Queue(sqlStatement, permission.AccessorId, permission.AccessPointId, permission.ChannelNo, 0, 0, permission.OrganisationId, permission.SerialNumber, permission.CredentialId)
	}

	conn, err := dbPool.Acquire(context.Background())
	if err != nil {
		logger.Log.Error(err.Error())
		return nil, utils.RenderAppError(DEVICE_PERMISSION_DAO_ERROR_BASE_CODE+2, fmt.Sprintf("Couldn't acquire database connection. Error:%s", err.Error()))
	}
	defer conn.Release()

	results := conn.SendBatch(context.Background(), batch)
	defer results.Close()

	for range permissions {

		var id int

		if err := results.QueryRow().Scan(&id); err != nil {
			logger.Log.Error(err.Error())
			return nil, utils.RenderAppError(DEVICE_PERMISSION_DAO_ERROR_BASE_CODE+3, fmt.Sprintf("Couldn't create device permission. Error:%s", err.Error()))
		}
		devicePermissionIds = append(devicePermissionIds, id)
	}

	return devicePermissionIds, nil
}

func (d *devPermDb) DeleteDevicePermissionsOnAccessPoint(accessPointId, organisationId int) *models.ApplicationError {

	sqlStatement := `DELETE FROM device_permissions WHERE "accessPointId"=$1 AND "organisationId"=$2;`

	_, err := dbPool.Exec(context.Background(), sqlStatement, accessPointId, organisationId)

	if err != nil {

		logger.Log.Error(err.Error())
		appError := utils.RenderAppError(DEVICE_PERMISSION_DAO_ERROR_BASE_CODE+4, fmt.Sprintf("Could not delete device permission of access point: %d, Error:%s", accessPointId, err.Error()))
		return appError
	}

	return nil
}

func (d *devPermDb) GetDevicePermissionId(devicePermission models.DevicePermissionsSchema) (devicePermissionId int, appError *models.ApplicationError) {

	sqlStatement := `
        SELECT id FROM device_permissions
        WHERE "organisationId"=$1 AND "accessPointId"=$2 AND "serialNumber"=$3 AND "accessorId"=$4 AND "credentialId"=$6 AND "channelNo"=$7
    `

	row := dbPool.QueryRow(context.Background(), sqlStatement,
		devicePermission.OrganisationId,
		devicePermission.AccessPointId,
		devicePermission.SerialNumber,
		devicePermission.AccessorId,
		devicePermission.CredentialId,
		devicePermission.ChannelNo,
	)

	err := row.Scan(&devicePermissionId)

	if err != nil {

		if err == pgx.ErrNoRows {

			return 0, nil
		}

		logger.Log.Error(err.Error())
		appError := utils.RenderAppError(DEVICE_PERMISSION_DAO_ERROR_BASE_CODE+5, fmt.Sprintf("Could not get device permission of access point: %d, Error:%s", devicePermission.AccessPointId, err.Error()))
		return 0, appError
	}

	return devicePermissionId, nil
}

func (d *devPermDb) SetToDeleteFlagForDevicePermission(devicePermissionId int) *models.ApplicationError {

	sqlStatement := `UPDATE device_permissions
					SET "toDelete"=$1
					WHERE "id"=$2`

	_, err := dbPool.Exec(context.Background(), sqlStatement, 1, devicePermissionId)

	if err != nil {

		logger.Log.Error(err.Error())
		appError := utils.RenderAppError(DEVICE_PERMISSION_DAO_ERROR_BASE_CODE+6, fmt.Sprintf("Could not update toDelete flag for device permission bearing id: %d, Error:%s", devicePermissionId, err.Error()))
		return appError
	}

	return nil
}

func (d *devPermDb) SetUpdatedOnDeviceFlagForDevicePermission(devicePermissionId int) *models.ApplicationError {

	sqlStatement := `UPDATE device_permissions
					SET "updatedOnDevice"=$1
					WHERE "id"=$2`

	_, err := dbPool.Exec(context.Background(), sqlStatement, 1, devicePermissionId)

	if err != nil {

		logger.Log.Error(err.Error())
		appError := utils.RenderAppError(DEVICE_PERMISSION_DAO_ERROR_BASE_CODE+7, fmt.Sprintf("Could not update updatedOnDevice flag for device permission bearing id: %d, Error:%s", devicePermissionId, err.Error()))
		return appError
	}

	return nil
}

func (d *devPermDb) GetDevicePermissionDetailsFromId(devicePermissionId int) (devicePermission *models.DevicePermissionsSchema, appError *models.ApplicationError) {

	devicePermission = &models.DevicePermissionsSchema{}

	sqlStatement := `SELECT "credentialId", "organisationId", "accessPointId", "accessorId", "serialNumber", "updatedOnDevice", "toDelete" ,"faceDataId" FROM device_permissions WHERE id=$1;`

	err := dbPool.QueryRow(context.Background(), sqlStatement, devicePermissionId).Scan(&devicePermission.CredentialId, &devicePermission.OrganisationId, &devicePermission.AccessPointId, &devicePermission.AccessorId, &devicePermission.SerialNumber, &devicePermission.UpdatedOnDevice, &devicePermission.ToDelete, &devicePermission.FaceId)

	if err != nil {

		logger.Log.Error(err.Error())
		appError := utils.RenderAppError(DEVICE_PERMISSION_DAO_ERROR_BASE_CODE+8, fmt.Sprintf("Could not retrieve credential Id, accessorId from device permission bearing id: %d, Error:%s", devicePermissionId, err.Error()))
		return devicePermission, appError
	}
	return devicePermission, nil
}

func (d *devPermDb) GetAccessPointIdFromId(devicePermissionId int) (accessPointId int, appError *models.ApplicationError) {

	sqlStatement := `SELECT "accessPointId" FROM device_permissions WHERE id=$1;`

	err := dbPool.QueryRow(context.Background(), sqlStatement, devicePermissionId).Scan(&accessPointId)

	if err != nil {

		logger.Log.Error(err.Error())
		appError := utils.RenderAppError(DEVICE_PERMISSION_DAO_ERROR_BASE_CODE+8, fmt.Sprintf("Could not retrieve credential Id, accessorId from device permission bearing id: %d, Error:%s", devicePermissionId, err.Error()))
		return accessPointId, appError
	}
	return accessPointId, nil
}

func (d *devPermDb) GetDevicePermissionFromId(devicePermissionId int) (devicePermission *models.DevicePermissionsSchema, appError *models.ApplicationError) {

	devicePermission = &models.DevicePermissionsSchema{}

	sqlStatement := `
        SELECT "organisationId", "serialNumber","credentialId", "accessorId" , "accessPointId", "updatedAt", "faceDataId", "updatedOnDevice", "toDelete" FROM device_permissions
        WHERE "id" =$1;
    `

	err := dbPool.QueryRow(context.Background(), sqlStatement, devicePermissionId).Scan(&devicePermission.OrganisationId, &devicePermission.SerialNumber, &devicePermission.CredentialId, &devicePermission.AccessorId, &devicePermission.AccessPointId, &devicePermission.UpdatedAt, &devicePermission.FaceId, &devicePermission.UpdatedOnDevice, &devicePermission.ToDelete)

	if err != nil {

		logger.Log.Error(err.Error())
		appError := utils.RenderAppError(DEVICE_PERMISSION_DAO_ERROR_BASE_CODE+9, fmt.Sprintf("Could not get device permission details from device permission bearing id: %d, Error:%s", devicePermissionId, err.Error()))
		return devicePermission, appError
	}

	return devicePermission, nil
}

func (d *devPermDb) CreateDeletedDevicePermission(devicePermission models.DeleteDevicePermissionsSchema) (appError *models.ApplicationError) {

	sqlStatement := `INSERT INTO deleted_device_permissions ("deletedDevicePermissionsTableId", "serialNumber", "channelNo", "credentialId", "accessorId", "accessPointId", "organisationId", "faceDataId", "assignedAt", "unAssignedAt") VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

	_, err := dbPool.Exec(context.Background(), sqlStatement, devicePermission.Id, devicePermission.SerialNumber, devicePermission.ChannelNo, devicePermission.CredentialId, devicePermission.AccessorId, devicePermission.AccessPointId, devicePermission.OrganisationId, devicePermission.FaceId, devicePermission.AssignedAt, devicePermission.UnAssignedAt)

	if err != nil {

		logger.Log.Error(err.Error())
		appError := utils.RenderAppError(DEVICE_PERMISSION_DAO_ERROR_BASE_CODE+10, fmt.Sprintf("Could not create deleted device permission of accessorId: %d and organisationId: %d, Error:%s", devicePermission.AccessorId, devicePermission.OrganisationId, err.Error()))
		return appError
	}

	return nil
}

func (d *devPermDb) DeleteDevicePermissionFromId(devicePermissionId int) (appError *models.ApplicationError) {

	sqlStatement := `DELETE FROM device_permissions WHERE "id"=$1 AND "toDelete"=$2;`

	_, err := dbPool.Exec(context.Background(), sqlStatement, devicePermissionId, 1)

	if err != nil {

		logger.Log.Error(err.Error())
		appError := utils.RenderAppError(DEVICE_PERMISSION_DAO_ERROR_BASE_CODE+11, fmt.Sprintf("Could not delete device permission bearing Id: %d , Error:%s", devicePermissionId, err.Error()))
		return appError
	}

	return nil
}

func (d *devPermDb) GetDevicePermissionIds(credentialId uint32, serialNumber string) (ids []int, appError *models.ApplicationError) {

	sqlStatement := `SELECT id FROM device_permissions WHERE "serialNumber"=$1 AND "credentialId"=$2;`

	rows, err := dbPool.Query(context.Background(), sqlStatement, serialNumber, credentialId)
	if err != nil {
		appError = utils.RenderAppError(DEVICE_PERMISSION_DAO_ERROR_BASE_CODE+12, fmt.Sprintf("Couldnt get dev perm Ids of credentialId Id: %d on device: %s.Error:%s", credentialId, serialNumber, err.Error()))
		return ids, appError
	}

	defer rows.Close()

	for rows.Next() {

		var id int
		err = rows.Scan(&id)
		if err != nil {
			logger.Log.Error(err.Error())
			appError = utils.RenderAppError(DEVICE_PERMISSION_DAO_ERROR_BASE_CODE+13, fmt.Sprintf("Error:%s", err.Error()))
			return ids, appError

		}

		ids = append(ids, id)
	}

	return ids, nil
}

func (d *devPermDb) GetDevicePermissionIdsOnAccessPoint(credentialId uint32, accessorId int, accessPointIds []int) (ids []int, apiError *models.ApplicationError) {

	if len(ids) == 0 {
		return nil, nil
	}

	sqlStatement := `SELECT id FROM device_permissions
	WHERE "accessorId"=$1 AND "credentialId"=$2 AND "accessPointId" = ANY($3);`

	rows, err := dbPool.Query(context.Background(), sqlStatement, accessorId, credentialId, pq.Array(accessPointIds))
	if err != nil {
		logger.Log.Error(err.Error())
		appError := utils.RenderAppError(DEVICE_PERMISSION_DAO_ERROR_BASE_CODE+14, fmt.Sprintf("Could not verify if device permissions exist. Error: %s", err.Error()))
		return nil, appError
	}

	defer rows.Close()

	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			logger.Log.Error(err.Error())
			appError := utils.RenderAppError(DEVICE_PERMISSION_DAO_ERROR_BASE_CODE+15, fmt.Sprintf("Error: %s", err.Error()))
			return nil, appError
		}
		ids = append(ids, id)
	}

	return ids, nil
}

func (d *devPermDb) GetLatestDevicePermission(devicePermissionId, channelNo int) (exists bool, devicePermission models.DevicePermissionsSchema, appError *models.ApplicationError) {

	sqlStatement := `
        SELECT "updatedOnDevice", "toDelete", "createdAt", "updatedAt" FROM device_permissions
        WHERE "id"=$1 AND "channelNo"=$2;
    `

	err := dbPool.QueryRow(context.Background(), sqlStatement, devicePermissionId, channelNo).Scan(&devicePermission.UpdatedOnDevice, &devicePermission.ToDelete, &devicePermission.CreatedAt, &devicePermission.UpdatedAt)

	if err != nil {

		if err == pgx.ErrNoRows {
			return false, devicePermission, nil
		}

		logger.Log.Error(err.Error())
		appError := utils.RenderAppError(DEVICE_PERMISSION_DAO_ERROR_BASE_CODE+16, fmt.Sprintf("Error: %s", err.Error()))
		return true, devicePermission, appError
	}

	return true, devicePermission, nil
}

func (d *devPermDb) GetDeletedDevicePermission(deletedDevicePermissionId int) (exists bool, deletedDevicePermission models.DevicePermissionsSchema, appError *models.ApplicationError) {

	sqlStatement := `
        SELECT "serialNumber","credentialId", "accessorId" , "accessPointId" , "organisationId", "createdAt", "deletedAt" FROM deleted_device_permissions
        WHERE "deletedDevicePermissionsId"=$1;
    `

	err := dbPool.QueryRow(context.Background(), sqlStatement, deletedDevicePermissionId).Scan(&deletedDevicePermission.SerialNumber, &deletedDevicePermission.CredentialId, &deletedDevicePermission.AccessorId, &deletedDevicePermission.AccessPointId, &deletedDevicePermission.OrganisationId, &deletedDevicePermission.CreatedAt, &deletedDevicePermission.DeletedAt)

	if err != nil {

		if err == pgx.ErrNoRows {
			return false, deletedDevicePermission, nil
		}

		logger.Log.Error(err.Error())
		appError := utils.RenderAppError(DEVICE_PERMISSION_DAO_ERROR_BASE_CODE+17, fmt.Sprintf("Error: %s", err.Error()))
		return false, deletedDevicePermission, appError
	}

	return true, deletedDevicePermission, nil
}

func (d *devPermDb) CheckIfAccessorIdAccessPointIdExistForCredentialIdAddPermission(credentialId uint32, accessorId, accessPointId int) (bool, *models.ApplicationError) {

	var count int

	sqlStatement := `
        SELECT COUNT(*) FROM device_permissions
        WHERE "credentialId" = $1 AND "accessorId" = $2 AND "accessPointId" = $3 AND "updatedOnDevice" = 0;
    `

	err := dbPool.QueryRow(context.Background(), sqlStatement, credentialId, accessorId, accessPointId).Scan(&count)
	if err != nil {
		if err == pgx.ErrNoRows {
			return false, nil
		}

		logger.Log.Error(err.Error())
		appError := utils.RenderAppError(DEVICE_PERMISSION_DAO_ERROR_BASE_CODE+18, fmt.Sprintf("Error: %s", err.Error()))
		return false, appError
	}

	return count > 0, nil
}

func (d *devPermDb) CheckIfAccessorIdAccessPointIdExistForCredentialIdRemovePermission(credentialId uint32, accessorId, accessPointId int) (bool, *models.ApplicationError) {

	var count int

	sqlStatement := `
        SELECT COUNT(*) FROM device_permissions
        WHERE "credentialId" = $1 AND "accessorId" = $2 AND "accessPointId" = $3;
    `

	err := dbPool.QueryRow(context.Background(), sqlStatement, credentialId, accessorId, accessPointId).Scan(&count)
	if err != nil {
		if err == pgx.ErrNoRows {
			return false, nil
		}

		logger.Log.Error(err.Error())
		appError := utils.RenderAppError(DEVICE_PERMISSION_DAO_ERROR_BASE_CODE+19, fmt.Sprintf("Error: %s", err.Error()))
		return false, appError
	}

	return count > 0, nil
}

func (d *devPermDb) CheckIfCredentialIdPresentInDevicePermission(credentialId uint32) (bool, *models.ApplicationError) {

	var count int

	sqlStatement := `
        SELECT COUNT(*) FROM device_permissions
        WHERE "credentialId" = $1;
    `

	err := dbPool.QueryRow(context.Background(), sqlStatement, credentialId).Scan(&count)
	if err != nil {
		if err == pgx.ErrNoRows {
			return false, nil
		}

		logger.Log.Error(err.Error())
		appError := utils.RenderAppError(DEVICE_PERMISSION_DAO_ERROR_BASE_CODE+20, fmt.Sprintf("Error: %s", err.Error()))
		return false, appError
	}

	return count > 0, nil
}

func (d *devPermDb) GetDevicePermissionsOfCredential(credentialId uint32) (accessPointIds []int, appError *models.ApplicationError) {

	sqlStatement := `SELECT DISTINCT "accessPointId" FROM device_permissions WHERE "credentialId"=$1;`

	rows, err := dbPool.Query(context.Background(), sqlStatement, credentialId)
	if err != nil {
		appError = utils.RenderAppError(DEVICE_PERMISSION_DAO_ERROR_BASE_CODE+12, fmt.Sprintf("Couldnt get dev perm access point id of credentialId Id: %d. Error:%s", credentialId, err.Error()))
		return accessPointIds, appError
	}

	defer rows.Close()

	for rows.Next() {

		var accessPointId int
		err = rows.Scan(&accessPointId)
		if err != nil {
			logger.Log.Error(err.Error())
			appError = utils.RenderAppError(DEVICE_PERMISSION_DAO_ERROR_BASE_CODE+21, fmt.Sprintf("Error:%s", err.Error()))
			return accessPointIds, appError

		}

		accessPointIds = append(accessPointIds, accessPointId)
	}

	return accessPointIds, nil
}

func (d *devPermDb) GetDevicePermissionsOnDevice(credentialId uint32) (accessPointIds []int, appError *models.ApplicationError) {

	sqlStatement := `SELECT "z" FROM device_permissions WHERE "accessPointId"=$1;`

	rows, err := dbPool.Query(context.Background(), sqlStatement, credentialId)
	if err != nil {
		appError = utils.RenderAppError(DEVICE_PERMISSION_DAO_ERROR_BASE_CODE+12, fmt.Sprintf("Couldnt get dev perm access point id of credentialId Id: %d. Error:%s", credentialId, err.Error()))
		return accessPointIds, appError
	}

	defer rows.Close()

	for rows.Next() {

		var accessPointId int
		err = rows.Scan(&accessPointId)
		if err != nil {
			logger.Log.Error(err.Error())
			appError = utils.RenderAppError(DEVICE_PERMISSION_DAO_ERROR_BASE_CODE+21, fmt.Sprintf("Error:%s", err.Error()))
			return accessPointIds, appError

		}

		accessPointIds = append(accessPointIds, accessPointId)
	}

	return accessPointIds, nil
}

func (d *devPermDb) GetTruncatedDevicePermission(serialNumber string, channelNo int, credentialId uint32, templateId int) (bool, *models.TruncatedDevicePermission, *models.ApplicationError) {

	truncatedDevicePermission := models.TruncatedDevicePermission{}

	sqlStatement := `SELECT id, "updatedOnDevice", "toDelete" FROM device_permissions
        WHERE "serialNumber"=$1 AND "channelNo"=$2 AND "credentialId"=$3 AND "faceDataId"=$4;`

	err := dbPool.QueryRow(context.Background(), sqlStatement, serialNumber, channelNo, credentialId, templateId).Scan(&truncatedDevicePermission.Id, &truncatedDevicePermission.UpdatedOnDevice, &truncatedDevicePermission.ToDelete)

	if err != nil {

		if err == pgx.ErrNoRows {
			return false, nil, nil
		}

		errMsg := fmt.Sprintf("GetTruncatedDevicePermission:Could not get device permissions. serialNumber: %s,channelNo: %d,credentialId: %d.Error:%s!", serialNumber, channelNo, credentialId, err.Error())
		logger.Log.Error(errMsg)
		appError := utils.RenderAppError(DEVICE_PERMISSION_DAO_ERROR_BASE_CODE+1, errMsg)
		return true, nil, appError
	}

	return true, &truncatedDevicePermission, nil
}

func (d *devPermDb) UpdateToDeleteFlag(devicePermissionsTableId, updateValue int) *models.ApplicationError {

	sqlStatement := `UPDATE device_permissions
					SET "toDelete"=$1
					WHERE "id"=$2`

	_, err := dbPool.Exec(context.Background(), sqlStatement, updateValue, devicePermissionsTableId)

	if err != nil {

		logger.Log.Error(err.Error())
		appError := utils.RenderAppError(DEVICE_PERMISSION_DAO_ERROR_BASE_CODE+7, fmt.Sprintf("Could not update updatedOnDevice flag for device permission bearing id: %d, Error:%s", devicePermissionsTableId, err.Error()))
		return appError
	}

	return nil
}

func (d *devPermDb) UpdateToDeleteAndUpdatedOnDeviceFlag(devicePermissionsTableId, setToDelete, setUpdatedOnDevice int) *models.ApplicationError {

	sqlStatement := `UPDATE device_permissions
					SET "toDelete"=$1, "updatedOnDevice"=$2
					WHERE "id"=$3`

	_, err := dbPool.Exec(context.Background(), sqlStatement, setToDelete, setUpdatedOnDevice, devicePermissionsTableId)

	if err != nil {

		logger.Log.Error(err.Error())
		appError := utils.RenderAppError(DEVICE_PERMISSION_DAO_ERROR_BASE_CODE+7, fmt.Sprintf("Could not update updatedOnDevice flag for device permission bearing id: %d, Error:%s", devicePermissionsTableId, err.Error()))
		return appError
	}

	return nil
}

func (d *devPermDb) CreateDevicePermissionWherePermissionIsSynced(devicePermission models.DevicePermissionsSchema) (devicePermissionId int, appError *models.ApplicationError) {

	sqlStatement := `INSERT INTO device_permissions ("accessorId", "accessPointId", "channelNo", "toDelete", "updatedOnDevice", "organisationId","serialNumber","credentialId") VALUES ($1, $2, $3, $4, $5,$6 ,$7, $8) RETURNING id `

	err := dbPool.QueryRow(context.Background(), sqlStatement, devicePermission.AccessorId, devicePermission.AccessPointId, devicePermission.ChannelNo, 0, 1, devicePermission.OrganisationId, devicePermission.SerialNumber, devicePermission.CredentialId).Scan(&devicePermissionId)
	if err != nil {
		errMsg := fmt.Sprintf("CreateDevicePermission:Could not create device permission.AccessorId:%d.CredentialId:%d.AccessPointId:%d.SerialNumber:%s.Error:%s!", devicePermission.AccessorId, devicePermission.CredentialId, devicePermission.AccessPointId, devicePermission.SerialNumber, err.Error())
		logger.Log.Error(errMsg)
		appError := utils.RenderAppError(DEVICE_PERMISSION_DAO_ERROR_BASE_CODE+1, errMsg)
		return 0, appError
	}

	return devicePermissionId, nil
}

func (d *devPermDb) GetDevicePermissionIdForTemplateNumber(templateNumber int, serialNumber string) (int, *models.ApplicationError) {

	var dpId int
	sqlStatement := `SELECT id FROM device_permissions
        			WHERE "templateNumber"=$1 AND "serialNumber"=$2`

	err := dbPool.QueryRow(context.Background(), sqlStatement, templateNumber, serialNumber).Scan(&dpId)
	if err != nil {
		errMsg := fmt.Sprintf("GetDevicePermissionIdForTemplateNumber:Could not get device permission. Error:%s!", err.Error())
		logger.Log.Error(errMsg)
		appError := utils.RenderAppError(DEVICE_PERMISSION_DAO_ERROR_BASE_CODE+1, errMsg)
		return 0, appError
	}

	return dpId, nil
}

func (d *devPermDb) GetTemplateIdFromDevicePermissions(organisationId int, accessorId int, serialNumber string) (exists bool, templateId int, appError *models.ApplicationError) {

	sqlStatement := `SELECT "faceDataId" FROM device_permissions WHERE "organisationId"=$1 AND "accessorId"=$2 AND "serialNumber"=$3`

	err := dbPool.QueryRow(context.Background(), sqlStatement, organisationId, accessorId, serialNumber).Scan(&templateId)
	if err != nil {
		errMsg := fmt.Sprintf("GetTemplateIdFromDevicePermissions:Could not get device permission. Error:%s!", err.Error())
		logger.Log.Error(errMsg)
		appError := utils.RenderAppError(DEVICE_PERMISSION_DAO_ERROR_BASE_CODE+1, errMsg)
		return false, 0, appError
	}

	return true, templateId, nil
}

// func (d *devPermDb) GetTemplateDetailsFromDevicePermissions(organisationId int, accessorId int, serialNumber string) (exists bool, templateId int, appError *models.ApplicationError) {

// 	sqlStatement := `SELECT "faceDataId" FROM device_permissions WHERE "organisationId"=$1 AND "accessorId"=$2 AND "serialNumber"=$3`

// 	err := dbPool.QueryRow(context.Background(), sqlStatement, organisationId, accessorId, serialNumber).Scan(&templateId)
// 	if err != nil {
// 		errMsg := fmt.Sprintf("CreateDevicePermission:Could not get device permission. Error:%s!", err.Error())
// 		logger.Log.Error(errMsg)
// 		appError := utils.RenderAppError(DEVICE_PERMISSION_DAO_ERROR_BASE_CODE+1, errMsg)
// 		return false, 0, appError
// 	}

// 	return true, templateId, nil
// }

func (d *devPermDb) GetDevicePermissionDetailsForTemplateId(credentialId uint32, serialNumber string, templateId int) (*models.DevicePermissionsSchema, *models.ApplicationError) {

	var devicePermission models.DevicePermissionsSchema
	sqlStatement := `SELECT "id", "credentialId", "organisationId", "accessPointId", "accessorId", "serialNumber", "updatedOnDevice", "toDelete" FROM device_permissions
        			WHERE "credentialId"=$1 AND "serialNumber"=$2 AND "faceDataId"=$3`

	err := dbPool.QueryRow(context.Background(), sqlStatement, credentialId, serialNumber, templateId).Scan(&devicePermission.Id, &devicePermission.CredentialId, &devicePermission.OrganisationId, &devicePermission.AccessPointId, &devicePermission.AccessorId, &devicePermission.SerialNumber, &devicePermission.UpdatedOnDevice, &devicePermission.ToDelete)
	if err != nil {
		errMsg := fmt.Sprintf("GetDevicePermissionDetailsForTemplateId:Could not get device permission. Error:%s!", err.Error())
		logger.Log.Error(errMsg)
		appError := utils.RenderAppError(DEVICE_PERMISSION_DAO_ERROR_BASE_CODE+1, errMsg)
		return nil, appError
	}

	return &devicePermission, nil
}

func (d *devPermDb) CheckIfAccessorIdAccessPointIdExistForTemplateIdAddPermission(templateId int, accessorId, accessPointId int) (bool, *models.ApplicationError) {

	var count int

	sqlStatement := `
        SELECT COUNT(*) FROM device_permissions
        WHERE "faceDataId" = $1 AND "accessorId" = $2 AND "accessPointId" = $3 AND "updatedOnDevice" = 0;
    `

	err := dbPool.QueryRow(context.Background(), sqlStatement, templateId, accessorId, accessPointId).Scan(&count)
	if err != nil {
		if err == pgx.ErrNoRows {
			return false, nil
		}

		logger.Log.Error(err.Error())
		appError := utils.RenderAppError(DEVICE_PERMISSION_DAO_ERROR_BASE_CODE+18, fmt.Sprintf("Error: %s", err.Error()))
		return false, appError
	}

	return count > 0, nil
}

func (d *devPermDb) CheckIfAccessorIdAccessPointIdTemplateIdExistForRemovePermission(templateId int, accessorId, accessPointId int) (bool, *models.ApplicationError) {

	var count int

	sqlStatement := `
        SELECT COUNT(*) FROM device_permissions
        WHERE "faceDataId" = $1 AND "accessorId" = $2 AND "accessPointId" = $3;
    `

	err := dbPool.QueryRow(context.Background(), sqlStatement, templateId, accessorId, accessPointId).Scan(&count)
	if err != nil {
		if err == pgx.ErrNoRows {
			return false, nil
		}

		logger.Log.Error(err.Error())
		appError := utils.RenderAppError(DEVICE_PERMISSION_DAO_ERROR_BASE_CODE+19, fmt.Sprintf("Error: %s", err.Error()))
		return false, appError
	}

	return count > 0, nil
}

func (d *devPermDb) CheckIfFaceIdPresentInDevicePermission(faceId int) (bool, *models.ApplicationError) {

	var count int

	sqlStatement := `
        SELECT COUNT(*) FROM device_permissions
        WHERE "faceDataId" = $1;
    `

	err := dbPool.QueryRow(context.Background(), sqlStatement, faceId).Scan(&count)
	if err != nil {
		if err == pgx.ErrNoRows {
			return false, nil
		}

		logger.Log.Error(err.Error())
		appError := utils.RenderAppError(DEVICE_PERMISSION_DAO_ERROR_BASE_CODE+20, fmt.Sprintf("Error: %s", err.Error()))
		return false, appError
	}

	return count > 0, nil
}

func (d *devPermDb) GetDevicePendingSyncPermissionsfromDeviceSerialNumber(serialNumber string) ([]models.DevicePermissionsSchema, *models.ApplicationError) {

	var devicePermissions []models.DevicePermissionsSchema
	sqlStatement := `SELECT "credentialId", "faceDataId" FROM device_permissions WHERE "serialNumber"=$1 AND "updatedOnDevice"=0 AND "toDelete"=0;`

	rows, err := dbPool.Query(context.Background(), sqlStatement, serialNumber)
	if err != nil {
		appError := utils.RenderAppError(DEVICE_PERMISSION_DAO_ERROR_BASE_CODE+12, fmt.Sprintf("Couldnt get devicePermissions on device: %s.Error:%s", serialNumber, err.Error()))
		return nil, appError
	}

	defer rows.Close()

	for rows.Next() {

		var devicePermission models.DevicePermissionsSchema
		err = rows.Scan(&devicePermission.CredentialId, &devicePermission.FaceId)
		if err != nil {
			logger.Log.Error(err.Error())
			appError := utils.RenderAppError(DEVICE_PERMISSION_DAO_ERROR_BASE_CODE+13, fmt.Sprintf("Error:%s", err.Error()))
			return nil, appError

		}

		devicePermissions = append(devicePermissions, devicePermission)
	}

	return devicePermissions, nil
}

func (d *devPermDb) UpdateCredentialIdforAccessorId(organisationId int, accessorId int, credentialId uint32) *models.ApplicationError {

	sqlStatement := `UPDATE device_permissions
	SET "credentialId"=$1
	WHERE "organisationId"=$2 AND "accessorId"=$3`

	_, err := dbPool.Exec(context.Background(), sqlStatement, credentialId, organisationId, accessorId)

	if err != nil {

		logger.Log.Error(err.Error())
		appError := utils.RenderAppError(DEVICE_PERMISSION_DAO_ERROR_BASE_CODE+6, fmt.Sprintf("Could not update UpdateCredentialIdforAccessorId for device permission bearing credentialId: %d, Error:%s", credentialId, err.Error()))
		return appError
	}

	return nil
}

// func (d *devPermDb) GetDevicePermissionsForOrganisation(organisationId int) (*[]models.SyncPendingPermDetails, *models.ApplicationError) {
// 	var syncPendingPermDetails []models.SyncPendingPermDetails
// 	sqlStatement := `SELECT DISTINCT "serialNumber", id FROM device_permissions WHERE "organisationId" = $1 AND ("updatedOnDevice"=0 AND "toDelete"=0);`

// 	rows, err := dbPool.Query(context.Background(), sqlStatement, organisationId)
// 	if err != nil {
// 		appError := utils.RenderAppError(DEVICE_PERMISSION_DAO_ERROR_BASE_CODE+12, fmt.Sprintf("Couldnt get devicePermissions on org: %s.Error:%s", organisationId, err.Error()))
// 		return nil, appError
// 	}

// 	defer rows.Close()

// 	for rows.Next() {

// 		var syncPendingPermDetail models.SyncPendingPermDetails
// 		err = rows.Scan(&syncPendingPermDetail.SerialNumber, &syncPendingPermDetail.DevicePermissionId)
// 		if err != nil {
// 			logger.Log.Error(err.Error())
// 			appError := utils.RenderAppError(DEVICE_PERMISSION_DAO_ERROR_BASE_CODE+13, fmt.Sprintf("Error:%s", err.Error()))
// 			return nil, appError

// 		}

// 		syncPendingPermDetails = append(syncPendingPermDetails, syncPendingPermDetail)
// 	}

// 	return &syncPendingPermDetails, nil
// }

func (d *devPermDb) UpdateSubIndexOfTemplateInDevicePermission(subIndex int, templateId int, credentialId uint32, serialNumber string) *models.ApplicationError {

	sqlStatement := `UPDATE device_permissions
	SET "subIndex"=$1
	WHERE "faceDataId"=$2 AND "credentialId"=$3 AND "serialNumber"=$4`

	_, err := dbPool.Exec(context.Background(), sqlStatement, subIndex, templateId, credentialId, serialNumber)

	if err != nil {

		logger.Log.Error(err.Error())
		appError := utils.RenderAppError(DEVICE_PERMISSION_DAO_ERROR_BASE_CODE+6, fmt.Sprintf("Could not update UpdateCredentialIdforAccessorId for device permission bearing credentialId: %d, Error:%s", credentialId, err.Error()))
		return appError
	}

	return nil
}

func (d *devPermDb) GetDevicePermissionFromAccessorIdAndDeviceSerialNumber(accessorId int, serialNumber string) (*models.DevicePermissionsSchema, *models.ApplicationError) {

	var dp models.DevicePermissionsSchema
	sqlStatement := `
    SELECT "organisationId", "serialNumber","credentialId", "faceDataId", "accessPointId", "updatedAt" FROM device_permissions
    WHERE "accessorId"=$1 AND "serialNumber" = $2;
    `

	err := dbPool.QueryRow(context.Background(), sqlStatement, accessorId, serialNumber).Scan(&dp.OrganisationId, &dp.SerialNumber, &dp.CredentialId, &dp.FaceId, &dp.AccessPointId, &dp.UpdatedAt)
	if err != nil {
		appError := utils.RenderAppError(DEVICE_PERMISSION_DAO_ERROR_BASE_CODE+12, fmt.Sprintf("Couldnt get dev perm Ids of credentialId Id: %d.Error:%s", accessorId, err.Error()))
		return nil, appError
	}

	return &dp, nil
}

func (d *devPermDb) ExcelrGetDevicePermissionFromAccessPointId(accessPointId int) ([]models.ExclerDevicePermission, *models.ApplicationError) {

	var devicePermissions []models.ExclerDevicePermission
	sqlStatement := `select dp."serialNumber", dp."accessorId", ofdi."faceDataId", ofdi."faceData", ofdi."userName" 
	from device_permissions dp 
	inner join organisation_face_data_ids ofdi 
	on dp."faceDataId" = ofdi."faceDataId" 
	where dp."accessPointId" = $1 and dp."updatedOnDevice"= 1 and dp."toDelete" = 0;`

	rows, err := dbPool.Query(context.Background(), sqlStatement, accessPointId)
	if err != nil {
		appError := utils.RenderAppError(DEVICE_PERMISSION_DAO_ERROR_BASE_CODE+12, fmt.Sprintf("Couldnt get devicePermissions on accessPointId: %s.Error:%s", accessPointId, err.Error()))
		return nil, appError
	}

	defer rows.Close()

	for rows.Next() {

		var devicePermission models.ExclerDevicePermission
		err = rows.Scan(&devicePermission.SerialNumber, &devicePermission.AccessorId, &devicePermission.FaceDataId, &devicePermission.FaceData, &devicePermission.Username)
		if err != nil {
			logger.Log.Error(err.Error())
			appError := utils.RenderAppError(DEVICE_PERMISSION_DAO_ERROR_BASE_CODE+13, fmt.Sprintf("Error:%s", err.Error()))
			return nil, appError

		}

		devicePermissions = append(devicePermissions, devicePermission)
	}

	return devicePermissions, nil
}
