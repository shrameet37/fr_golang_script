package database

import (
	"context"
	"face_management/logger"
	"face_management/models"
	"face_management/utils"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
)

type permDb struct{}

type permDbInterface interface {
	CreateAccessorPermission(permission models.AccessorPermission) (permissionId int, appError *models.ApplicationError)
	CreateAccessorPermissions(permissions []models.AccessorPermission) *models.ApplicationError
	DeleteAccessorPermission(permission models.AccessorPermission) *models.ApplicationError
	RemoveAccessorPermissions(permission models.AccessorPermission) *models.ApplicationError
	DeleteAccessorPermissions(permissions []models.AccessorPermission) *models.ApplicationError
	DeleteAccessorPermissionsOnAccessPoint(accessPointId, orgId int) *models.ApplicationError
	DeleteAccessorPermissionsForAccessPointDelete(accessPointId int) *models.ApplicationError
	GetAccessorPermissionsOfAccessor(accessorId, orgId int) (permissions []int, appError *models.ApplicationError)
	GetAccessorPermissionsOfAccessorWithCreatedAt(accessorId, orgId int) (permissions []models.AccessorPermission, appError *models.ApplicationError)
	GetAccessorPermissionsOnAccessPoint(accessPointId, orgId int) (permissions []int, appError *models.ApplicationError)
	GetAccessorPermissionsOnAccessPointWithCreatedAt(accessPointId, orgId int) (permissions []models.AccessorPermission, appError *models.ApplicationError)
	GetAccessorPermissionCountOfAccessor(accessorId int, orgId int) (count int, appError *models.ApplicationError)
	GetAccessorPermissionCountOfAccessPoint(accessPointId int, orgId int) (count int, appError *models.ApplicationError)
	DoesAccessorHavePermissionOnAccessPoint(accessPointId, accessorId, orgId int) (bool, *models.ApplicationError)
	DoesAccessorHavePermissionInOrganisation(accessorId, orgId int) (bool, *models.ApplicationError)
}

var PermDb permDbInterface

const (
	ACCESSOR_PERMISSION_DAO_BYTE3_VALUE = 0x00320000 //3 = dao, 2 = accessor permission
)

var (
	ACCESSOR_PERMISSION_DAO_ERROR_BASE_CODE int
)

func init() {
	PermDb = &permDb{}
	ACCESSOR_PERMISSION_DAO_ERROR_BASE_CODE = utils.ServiceErrorBaseCode + ACCESSOR_PERMISSION_DAO_BYTE3_VALUE

}

func (d *permDb) DeleteAccessorPermissionsOnAccessPoint(accessPointId, orgId int) *models.ApplicationError {

	sqlStatement := `DELETE FROM accessor_permissions WHERE "accessPointId"=$1 AND "organisationId"=$2;`

	_, err := dbPool.Exec(context.Background(), sqlStatement, accessPointId, orgId)

	if err != nil {

		logger.Log.Error(err.Error())
		appError := utils.RenderAppError(ACCESSOR_PERMISSION_DAO_ERROR_BASE_CODE+1, fmt.Sprintf("Could not delete permissions of access point: %d .Error:%s", accessPointId, err.Error()))
		return appError
	}

	return nil
}

func (d *permDb) DeleteAccessorPermissionsForAccessPointDelete(accessPointId int) *models.ApplicationError {

	sqlStatement := `DELETE FROM accessor_permissions WHERE "accessPointId"=$1;`

	_, err := dbPool.Exec(context.Background(), sqlStatement, accessPointId)

	if err != nil {

		logger.Log.Error(err.Error())
		appError := utils.RenderAppError(ACCESSOR_PERMISSION_DAO_ERROR_BASE_CODE+1, fmt.Sprintf("Could not delete permissions of access point: %d .Error:%s", accessPointId, err.Error()))
		return appError
	}

	return nil
}

func (d *permDb) GetAccessorPermissionsOfAccessor(accessorId int, orgId int) (apIds []int, appError *models.ApplicationError) {

	log.Println("CP-F4")
	sqlStatement := `SELECT "accessPointId"  FROM accessor_permissions WHERE "accessorId"=$1 AND "organisationId"=$2`

	rows, err := dbPool.Query(context.Background(), sqlStatement, accessorId, orgId)
	if err != nil {
		appError = utils.RenderAppError(ACCESSOR_PERMISSION_DAO_ERROR_BASE_CODE+2, fmt.Sprintf("Couldnt get permission of accessor Id: %d in orgId: %d.Error:%s", accessorId, orgId, err.Error()))
		return apIds, appError
	}

	defer rows.Close()

	for rows.Next() {

		var apId int

		err = rows.Scan(&apId)
		if err != nil {
			logger.Log.Error(err.Error())
			appError = utils.RenderAppError(ACCESSOR_PERMISSION_DAO_ERROR_BASE_CODE+3, fmt.Sprintf("Error:%s", err.Error()))
			return apIds, appError

		}

		apIds = append(apIds, apId)
	}

	return apIds, nil
}

func (d *permDb) GetAccessorPermissionCountOfAccessor(accessorId int, orgId int) (count int, appError *models.ApplicationError) {

	sqlStatement := `SELECT COUNT(*) FROM accessor_permissions WHERE "accessorId" = $1 AND "organisationId" = $2;`

	err := dbPool.QueryRow(context.Background(), sqlStatement, accessorId, orgId).Scan(&count)

	if err != nil {

		if err == pgx.ErrNoRows {

			count = 0
			return count, nil
		}

		appError := utils.RenderAppError(ACCESSOR_PERMISSION_DAO_ERROR_BASE_CODE+4, fmt.Sprintf("Couldnt get permission count of accessor Id: %d in orgId: %d.Error:%s", accessorId, orgId, err.Error()))
		return count, appError

	}

	return count, nil
}

func (d *permDb) GetAccessorPermissionCountOfAccessPoint(accessPointId int, orgId int) (count int, appError *models.ApplicationError) {

	sqlStatement := `SELECT COUNT(*) FROM accessor_permissions WHERE "accessPointId" = $1 AND "organisationId" = $2;`

	err := dbPool.QueryRow(context.Background(), sqlStatement, accessPointId, orgId).Scan(&count)

	if err != nil {

		if err == pgx.ErrNoRows {

			count = -1
			return count, nil
		}

		appError := utils.RenderAppError(ACCESSOR_PERMISSION_DAO_ERROR_BASE_CODE+5, fmt.Sprintf("Couldnt get permission count of accessPointId: %d in orgId: %d.Error:%s", accessPointId, orgId, err.Error()))
		return count, appError

	}

	return count, nil
}

func (d *permDb) GetAccessorPermissionsOnAccessPoint(accessPointId, orgId int) (permissions []int, appError *models.ApplicationError) {

	sqlStatement := `SELECT "accessorId" FROM accessor_permissions WHERE "accessPointId"=$1 and "organisationId"=$2;`

	rows, err := dbPool.Query(context.Background(), sqlStatement, accessPointId, orgId)

	if err != nil {

		logger.Log.Error(err.Error())
		appError := utils.RenderAppError(ACCESSOR_PERMISSION_DAO_ERROR_BASE_CODE+6, fmt.Sprintf("Couldnt get permissions on access point Id:%d.Error:%s", accessPointId, err.Error()))
		return permissions, appError
	}

	defer rows.Close()

	for rows.Next() {

		var permission int

		err = rows.Scan(&permission)
		if err != nil {
			logger.Log.Error(err.Error())
			appError = utils.RenderAppError(ACCESSOR_PERMISSION_DAO_ERROR_BASE_CODE+17, fmt.Sprintf("Error:%s", err.Error()))
			return permissions, appError

		}

		permissions = append(permissions, permission)
	}

	return permissions, nil
}

func (d *permDb) CreateAccessorPermission(permission models.AccessorPermission) (permissionId int, appError *models.ApplicationError) {

	sqlStatement := `INSERT INTO accessor_permissions ("accessorId", "accessPointId", "organisationId") VALUES ($1, $2, $3) RETURNING id`

	err := dbPool.QueryRow(context.Background(), sqlStatement, permission.AccessorId, permission.AccessPointId, permission.OrgId).Scan(&permissionId)

	if err != nil {

		logger.Log.Error(err.Error())
		appError = utils.RenderAppError(ACCESSOR_PERMISSION_DAO_ERROR_BASE_CODE+7, fmt.Sprintf("Couldnt create permission for accessor Id:%d on access point Id:%d.Error:%s", permission.AccessorId, permission.AccessPointId, err.Error()))
		return 0, appError
	}

	return permissionId, nil
}

func (d *permDb) CreateAccessorPermissions(permissions []models.AccessorPermission) *models.ApplicationError {

	batch := &pgx.Batch{}

	for _, permission := range permissions {
		sqlStatement := `INSERT INTO accessor_permissions ("accessorId", "accessPointId", "organisationId") VALUES ($1, $2, $3)`
		batch.Queue(sqlStatement, permission.AccessorId, permission.AccessPointId, permission.OrgId)
	}

	conn, err := dbPool.Acquire(context.Background())
	if err != nil {
		logger.Log.Error(err.Error())
		return utils.RenderAppError(ACCESSOR_PERMISSION_DAO_ERROR_BASE_CODE+8, fmt.Sprintf("Couldn't acquire database connection. Error:%s", err.Error()))
	}

	defer conn.Release()

	results := conn.SendBatch(context.Background(), batch)

	defer results.Close()

	for range permissions {
		if _, err := results.Exec(); err != nil {
			logger.Log.Error(err.Error())
			return utils.RenderAppError(ACCESSOR_PERMISSION_DAO_ERROR_BASE_CODE+9, fmt.Sprintf("Couldn't create permission. Error:%s", err.Error()))
		}
	}

	return nil
}

func (d *permDb) DeleteAccessorPermission(permission models.AccessorPermission) *models.ApplicationError {

	sqlStatement := `DELETE FROM accessor_permissions WHERE "accessorId"=$1 AND "accessPointId"=$2 AND "organisationId"=$3;`

	_, err := dbPool.Exec(context.Background(), sqlStatement, permission.AccessorId, permission.AccessPointId, permission.OrgId)

	if err != nil {

		logger.Log.Error(err.Error())
		appError := utils.RenderAppError(ACCESSOR_PERMISSION_DAO_ERROR_BASE_CODE+10, fmt.Sprintf("Could not delete permissions of accessor %d on access point: %d .Error:%s", permission.AccessorId, permission.AccessPointId, err.Error()))
		return appError
	}

	return nil
}

func (d *permDb) RemoveAccessorPermissions(permission models.AccessorPermission) *models.ApplicationError {

	sqlStatement := `DELETE FROM accessor_permissions WHERE "accessorId"=$1 AND "organisationId"=$2;`

	_, err := dbPool.Exec(context.Background(), sqlStatement, permission.AccessorId, permission.OrgId)

	if err != nil {

		logger.Log.Error(err.Error())
		appError := utils.RenderAppError(ACCESSOR_PERMISSION_DAO_ERROR_BASE_CODE+10, fmt.Sprintf("Could not delete permissions of accessor %d on access points.Error:%s", permission.AccessorId, err.Error()))
		return appError
	}

	return nil
}

func (d *permDb) DeleteAccessorPermissions(permissions []models.AccessorPermission) *models.ApplicationError {

	batch := &pgx.Batch{}

	for _, permission := range permissions {
		sqlStatement := `DELETE FROM accessor_permissions WHERE "accessorId"=$1 AND "accessPointId"=$2 AND "organisationId"=$3;`
		batch.Queue(sqlStatement, permission.AccessorId, permission.AccessPointId, permission.OrgId)
	}

	conn, err := dbPool.Acquire(context.Background())
	if err != nil {
		logger.Log.Error(err.Error())
		return utils.RenderAppError(ACCESSOR_PERMISSION_DAO_ERROR_BASE_CODE+11, fmt.Sprintf("Couldn't acquire database connection. Error:%s", err.Error()))
	}

	defer conn.Release()

	results := conn.SendBatch(context.Background(), batch)

	defer results.Close()

	for range permissions {
		if _, err := results.Exec(); err != nil {
			logger.Log.Error(err.Error())
			return utils.RenderAppError(ACCESSOR_PERMISSION_DAO_ERROR_BASE_CODE+12, fmt.Sprintf("Couldn't create permission. Error:%s", err.Error()))
		}
	}

	return nil
}

func (d *permDb) GetAccessorPermissionsOfAccessorWithCreatedAt(accessorId, orgId int) (permissions []models.AccessorPermission, appError *models.ApplicationError) {

	sqlStatement := `SELECT DISTINCT "accessPointId", "createdAt" FROM accessor_permissions WHERE "accessorId"=$1 AND "organisationId"=$2`

	rows, err := dbPool.Query(context.Background(), sqlStatement, accessorId, orgId)
	if err != nil {
		appError = utils.RenderAppError(ACCESSOR_PERMISSION_DAO_ERROR_BASE_CODE+13, fmt.Sprintf("Couldnt get permission of accessor Id: %d in orgId: %d.Error:%s", accessorId, orgId, err.Error()))
		return permissions, appError
	}

	defer rows.Close()

	for rows.Next() {

		var permission models.AccessorPermission

		err = rows.Scan(&permission.AccessPointId, &permission.CreatedAt)
		if err != nil {
			logger.Log.Error(err.Error())
			appError = utils.RenderAppError(ACCESSOR_PERMISSION_DAO_ERROR_BASE_CODE+14, fmt.Sprintf("Error:%s", err.Error()))
			return permissions, appError

		}

		permissions = append(permissions, permission)
	}

	return permissions, nil
}

func (d *permDb) GetAccessorPermissionsOnAccessPointWithCreatedAt(accessPointId, orgId int) (permissions []models.AccessorPermission, appError *models.ApplicationError) {

	sqlStatement := `SELECT "accessorId", "createdAt" FROM accessor_permissions WHERE "accessPointId"=$1 and "organisationId"=$2;`

	rows, err := dbPool.Query(context.Background(), sqlStatement, accessPointId, orgId)

	if err != nil {

		logger.Log.Error(err.Error())
		appError := utils.RenderAppError(ACCESSOR_PERMISSION_DAO_ERROR_BASE_CODE+15, fmt.Sprintf("Couldnt get permissions on access point Id:%d.Error:%s", accessPointId, err.Error()))
		return permissions, appError
	}

	defer rows.Close()

	for rows.Next() {

		var permission models.AccessorPermission

		err = rows.Scan(&permission.AccessorId, &permission.CreatedAt)
		if err != nil {
			logger.Log.Error(err.Error())
			appError = utils.RenderAppError(ACCESSOR_PERMISSION_DAO_ERROR_BASE_CODE+16, fmt.Sprintf("Error:%s", err.Error()))
			return permissions, appError

		}

		permissions = append(permissions, permission)
	}

	return permissions, nil
}

func (d *permDb) DoesAccessorHavePermissionOnAccessPoint(accessPointId, accessorId, orgId int) (bool, *models.ApplicationError) {

	var exists bool
	sqlStatement := `SELECT EXISTS (SELECT 1 FROM accessor_permissions WHERE "accessPointId" = $1 and "accessorId" = $2 and "organisationId" = $3)`

	err := dbPool.QueryRow(context.Background(), sqlStatement, accessPointId, accessorId, orgId).Scan(&exists)
	if err != nil {
		if err == pgx.ErrNoRows {
			return false, nil
		}

		logger.Log.Error(err.Error())
		appError := utils.RenderAppError(ACCESSOR_PERMISSION_DAO_ERROR_BASE_CODE+1, fmt.Sprintf("DoesAccessorHavePermissionOnAccessPoint. Error:%s", err.Error()))
		return false, appError
	}

	return exists, nil
}

func (d *permDb) DoesAccessorHavePermissionInOrganisation(accessorId, orgId int) (bool, *models.ApplicationError) {

	log.Println("CP-F3")
	var exists bool
	sqlStatement := `SELECT EXISTS (SELECT 1 FROM accessor_permissions WHERE "accessorId" = $1 and "organisationId" = $2)`

	err := dbPool.QueryRow(context.Background(), sqlStatement, accessorId, orgId).Scan(&exists)
	if err != nil {
		if err == pgx.ErrNoRows {
			return false, nil
		}

		logger.Log.Error(err.Error())
		appError := utils.RenderAppError(ACCESSOR_PERMISSION_DAO_ERROR_BASE_CODE+1, fmt.Sprintf("DoesAccessorHavePermissionOnAccessPoint. Error:%s", err.Error()))
		return false, appError
	}

	return exists, nil
}
