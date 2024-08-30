package database

import (
	"context"
	"face_management/logger"
	"face_management/models"
	"face_management/utils"

	"fmt"

	"github.com/jackc/pgx/v5"
)

type orgAccDb struct{}

type orgAccDbInterface interface {
	GetCredentialIdForAccessorId(orgId int, accessorId int) (credentilId *uint32, appError *models.ApplicationError)
	DoesAccessorExistInOrg(accessorId, orgId int) (exists bool, appError *models.ApplicationError)
	CreateAccessorInOrg(organisationId, accessorId int, credentialId uint32) *models.ApplicationError
	UpdateAccessorCredential(organisationId, accessorId int, credentialId uint32) *models.ApplicationError
	DeleteAccessorFromOrg(organisationId, accessorId int) *models.ApplicationError
	GetAccessorIdForCredntialId(organisationId int, credentialId uint32) (int, *models.ApplicationError)
}

var OrgAccDb orgAccDbInterface

func init() {
	OrgAccDb = &orgAccDb{}
}

func (d *orgAccDb) GetCredentialIdForAccessorId(orgId int, accessorId int) (*uint32, *models.ApplicationError) {

	var credId uint32
	sqlStatement := `SELECT "credentialId" FROM organisation_accessors WHERE "organisationId"=$1 AND "accessorId"=$2;`

	err := dbPool.QueryRow(context.Background(), sqlStatement, orgId, accessorId).Scan(&credId)
	if err != nil {
		logger.Log.Error(err.Error())
		appError := utils.RenderAppError(ACCESS_POINT_DAO_ERROR_BASE_CODE+1, fmt.Sprintf("GetCredentialIdForAccessorId. Error:%s", err.Error()))
		return nil, appError
	}

	return &credId, nil
}

// numberExists checks if the given number exists in the database
func (d *orgAccDb) DoesAccessorExistInOrg(accessorId, orgId int) (exists bool, appError *models.ApplicationError) {

	query := `SELECT EXISTS (SELECT 1 FROM organisation_accessors WHERE "accessorId" = $1 and "organisationId" = $2)`

	err := dbPool.QueryRow(context.Background(), query, accessorId, orgId).Scan(&exists)
	if err != nil {
		if err == pgx.ErrNoRows {
			return false, nil
		}

		logger.Log.Error(err.Error())
		appError = utils.RenderAppError(ACCESSOR_PERMISSION_DAO_ERROR_BASE_CODE+1, fmt.Sprintf("DoesPasscodeExistforOrg. Error:%s", err.Error()))
		return false, appError
	}

	return exists, nil
}

func (d *orgAccDb) CreateAccessorInOrg(organisationId, accessorId int, credentialId uint32) *models.ApplicationError {

	sqlStatement := `INSERT INTO organisation_accessors ("organisationId", "accessorId", "credentialId") VALUES ($1, $2, $3)`

	_, err := dbPool.Exec(context.Background(), sqlStatement, organisationId, accessorId, credentialId)

	if err != nil {

		logger.Log.Error(err.Error())
		appError := utils.RenderAppError(DEVICE_PERMISSION_DAO_ERROR_BASE_CODE+10, fmt.Sprintf("Could not CreateAccessorInOrg of accessorId: %d and organisationId: %d, Error:%s", accessorId, organisationId, err.Error()))
		return appError
	}

	return nil
}

func (d *orgAccDb) UpdateAccessorCredential(organisationId, accessorId int, credentialId uint32) *models.ApplicationError {

	sqlStatement := `UPDATE organisation_accessors SET "credentialId"=$1 WHERE "organisationId"=$2 AND "accessorId" = $3;`

	_, err := dbPool.Exec(context.Background(), sqlStatement, credentialId, organisationId, accessorId)

	if err != nil {

		logger.Log.Error(err.Error())
		appError := utils.RenderAppError(DEVICE_PERMISSION_DAO_ERROR_BASE_CODE+10, fmt.Sprintf("Could not CreateAccessorInOrg of accessorId: %d and organisationId: %d, Error:%s", accessorId, organisationId, err.Error()))
		return appError
	}

	return nil
}

func (d *orgAccDb) DeleteAccessorFromOrg(organisationId, accessorId int) *models.ApplicationError {

	sqlStatement := `DELETE FROM organisation_accessors WHERE "organisationId"=$1 AND "accessorId"=$2`

	_, err := dbPool.Exec(context.Background(), sqlStatement, organisationId, accessorId)

	if err != nil {
		logger.Log.Error(err.Error())
		appError := utils.RenderAppError(DEVICE_PERMISSION_DAO_ERROR_BASE_CODE+4, fmt.Sprintf("Could not delete accessorId %d.Error:%s", accessorId, err.Error()))
		return appError
	}

	return nil
}

func (d *orgAccDb) GetAccessorIdForCredntialId(organisationId int, credentialId uint32) (int, *models.ApplicationError) {

	var accessorId int
	sqlStatement := `SELECT "accessorId" FROM organisation_accessors WHERE "organisationId"=$1 AND "credentialId"=$2;`

	err := dbPool.QueryRow(context.Background(), sqlStatement, organisationId, credentialId).Scan(&accessorId)
	if err != nil {
		logger.Log.Error(err.Error())
		appError := utils.RenderAppError(ACCESS_POINT_DAO_ERROR_BASE_CODE+1, fmt.Sprintf("GetCredentialIdForAccessorId. Error:%s", err.Error()))
		return 0, appError
	}

	return accessorId, nil
}
