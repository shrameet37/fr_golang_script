package database

import (
	"context"
	"encoding/hex"
	"face_management/logger"
	"face_management/models"
	"face_management/utils"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
)

// TODO
type ofdiDb struct{}

type ofdiDbInterface interface {
	DoesFaceExistforOrg(faceData string, orgId int) (exists bool, faceDataId int, appError *models.ApplicationError)
	DeleteFaceFromOrgFacedataId(faceDataId int) (appError *models.ApplicationError)
	DoesFaceDataIdExistforOrg(faceDataId, orgId int) (exists bool, appError *models.ApplicationError)
	GetFaceDataIdForAccessorId(organisationId int, accessorId int) (bool, int, *models.ApplicationError)
	DoesFaceDataIdExistForAccessor(faceDataId int, orgId int, accessorId int) (exists bool, appError *models.ApplicationError)
	UpdateFaceforAccessor(orgId int, accessorId int, oldFace string, newFace string) (appError *models.ApplicationError)
	UpdateFaceForFaceDataId(newFace string, faceDataId int) (appError *models.ApplicationError)
	GetFaceDataIdDetails(faceDataId int) (bool, *models.OrganisationFaceIdsSchema, *models.ApplicationError)
	GetFaceDetailsFromFace(faceData string, orgId int) (bool, *models.OrganisationFaceIdsSchema, *models.ApplicationError)
	SetPendingUnassignOnDevicesFLag(faceDataId int, orgId int) (appError *models.ApplicationError)
	GetFaceDetailsForAccessorInOrg(orgId int, accessorId int) (bool, *models.OrganisationFaceIdsSchema, *models.ApplicationError)
	UpdateAccessorForFaceDataId(accessorId int, faceDataId int) (appError *models.ApplicationError)
	AddOrganisationFace(organisationId int, accessorId int, faceData string, userName string) (int, *models.ApplicationError)
	UpdateOrganisationFace(organisationId int, accessorId int, faceData string, userName string, faceDataId int) *models.ApplicationError
	DoesFaceExist(faceId int) (exists bool, appError *models.ApplicationError)
}

var OfdiDb ofdiDbInterface

const (
	ORGANISATION_FACE_DATA_IDS_DAO_BYTE3_VALUE = 0x003dffff //3 = dao, 7 = organisation
)

var (
	ORGANISATION_FACE_DATA_IDS_DAO_ERROR_BASE_CODE int
)

func init() {
	OfdiDb = &ofdiDb{}
	ORGANISATION_FACE_DATA_IDS_DAO_ERROR_BASE_CODE = utils.ServiceErrorBaseCode + ORGANISATION_FACE_DATA_IDS_DAO_BYTE3_VALUE
}

// numberExists checks if the given number exists in the database
func (d *ofdiDb) DoesFaceExistforOrg(faceData string, orgId int) (exists bool, faceDataId int, appError *models.ApplicationError) {

	query := `SELECT "faceDataId" FROM organisation_face_data_ids WHERE "faceData" = $1 and "organisationId" = $2`

	err := dbPool.QueryRow(context.Background(), query, faceData, orgId).Scan(&faceDataId)
	if err != nil {
		if err == pgx.ErrNoRows {
			return false, 0, nil
		}

		logger.Log.Error(err.Error())
		appError = utils.RenderAppError(ORGANISATION_FACE_DATA_IDS_DAO_ERROR_BASE_CODE+1, fmt.Sprintf("DoesFaceExistforOrg. Error:%s", err.Error()))
		return false, 0, appError
	}

	return true, faceDataId, nil
}

func (d *ofdiDb) DoesFaceDataIdExistforOrg(faceDataId, orgId int) (exists bool, appError *models.ApplicationError) {

	query := `SELECT EXISTS (SELECT 1 FROM organisation_face_data_ids WHERE "faceDataId" = $1 and "organisationId" = $2)`

	err := dbPool.QueryRow(context.Background(), query, faceDataId, orgId).Scan(&exists)
	if err != nil {
		if err == pgx.ErrNoRows {
			return false, nil
		}

		logger.Log.Error(err.Error())
		appError = utils.RenderAppError(ORGANISATION_FACE_DATA_IDS_DAO_ERROR_BASE_CODE+1, fmt.Sprintf("DoesFaceDataIdExistforOrg. Error:%s", err.Error()))
		return false, appError
	}

	return exists, nil
}

func (d *ofdiDb) GetFaceDataIdDetails(faceDataId int) (bool, *models.OrganisationFaceIdsSchema, *models.ApplicationError) {

	var response models.OrganisationFaceIdsSchema
	var faceData []byte
	query := `SELECT "faceDataId", "faceData", "accessorId", "organisationId", "userName", "pendingUnassignedOnDevices" FROM organisation_face_data_ids WHERE "faceDataId" = $1`

	err := dbPool.QueryRow(context.Background(), query, faceDataId).Scan(&response.FaceDataId, &faceData, &response.AccessorId, &response.OrganisationId, &response.UserName, &response.PendingUnassignedOnDevices)
	if err != nil {
		if err == pgx.ErrNoRows {
			return false, nil, nil
		}

		logger.Log.Error(err.Error())
		appError := utils.RenderAppError(ORGANISATION_FACE_DATA_IDS_DAO_ERROR_BASE_CODE+1, fmt.Sprintf("GetFaceDataIdDetails. Error:%s", err.Error()))
		return false, nil, appError
	}

	response.FaceData = string(faceData)

	return true, &response, nil
}

func (d *ofdiDb) DeleteFaceFromOrgFacedataId(faceDataId int) (appError *models.ApplicationError) {

	sqlStatement := `DELETE FROM organisation_face_data_ids WHERE "faceDataId"=$1;`

	_, err := dbPool.Exec(context.Background(), sqlStatement, faceDataId)

	if err != nil {

		logger.Log.Error(err.Error())
		appError = utils.RenderAppError(ACCESSOR_PERMISSION_DAO_ERROR_BASE_CODE+7, fmt.Sprintf("Couldnt DeleteFaceFromOrgFacedataId .Error:%s", err.Error()))
		return appError
	}

	return nil
}

func (d *ofdiDb) GetFaceDataIdForAccessorId(organisationId int, accessorId int) (bool, int, *models.ApplicationError) {

	var faceDataId int
	sqlStatement := `SELECT "faceDataId" FROM organisation_face_data_ids WHERE "organisationId" = $1 AND "accessorId" = $2;`

	err := dbPool.QueryRow(context.Background(), sqlStatement, organisationId, accessorId).Scan(&faceDataId)
	if err != nil {
		if err == pgx.ErrNoRows {
			return false, 0, nil
		}
		logger.Log.Error(err.Error())
		apiError := utils.RenderAppError(ACCESS_POINT_DAO_ERROR_BASE_CODE+7, fmt.Sprintf("Could not get access points under organisation.Error:%s", err.Error()))
		return false, 0, apiError
	}

	return true, faceDataId, nil
}

func (d *ofdiDb) DoesFaceDataIdExistForAccessor(faceDataId int, orgId int, accessorId int) (exists bool, appError *models.ApplicationError) {

	query := `SELECT EXISTS (SELECT 1 FROM organisation_face_data_ids WHERE "faceDataId" = $1 and "organisationId" = $2 and "accessorId" = $3)`

	err := dbPool.QueryRow(context.Background(), query, faceDataId, orgId, accessorId).Scan(&exists)
	if err != nil {
		if err == pgx.ErrNoRows {
			return false, nil
		}

		logger.Log.Error(err.Error())
		appError = utils.RenderAppError(ORGANISATION_FACE_DATA_IDS_DAO_ERROR_BASE_CODE+1, fmt.Sprintf("DoesFaceDataIdExistForAccessor. Error:%s", err.Error()))
		return false, appError
	}

	return exists, nil
}

func (d *ofdiDb) GetFaceDetailsForAccessorInOrg(orgId int, accessorId int) (bool, *models.OrganisationFaceIdsSchema, *models.ApplicationError) {

	var result models.OrganisationFaceIdsSchema
	var faceData []byte

	query := `SELECT "faceDataId", "faceData", "pendingUnassignedOnDevices" FROM organisation_face_data_ids WHERE "organisationId" = $1 and "accessorId" = $2`

	err := dbPool.QueryRow(context.Background(), query, orgId, accessorId).Scan(&result.FaceDataId, &faceData, &result.PendingUnassignedOnDevices)
	if err != nil {
		if err == pgx.ErrNoRows {
			return false, nil, nil
		}

		logger.Log.Error(err.Error())
		appError := utils.RenderAppError(ORGANISATION_FACE_DATA_IDS_DAO_ERROR_BASE_CODE+1, fmt.Sprintf("GetFaceDetailsForAccessorInOrg. Error:%s", err.Error()))
		return false, nil, appError
	}

	result.FaceData = hex.EncodeToString(faceData)

	return true, &result, nil
}

func (d *ofdiDb) UpdateFaceforAccessor(orgId int, accessorId int, oldFace string, newFace string) (appError *models.ApplicationError) {

	sqlStatement := `UPDATE organisation_face_data_ids SET "faceData" = $1 WHERE "faceData" = $2 AND "organisationId" = $3 AND "accessorId" = $4`

	_, err := dbPool.Exec(context.Background(), sqlStatement, newFace, oldFace, orgId, accessorId)
	if err != nil {
		logger.Log.Error(err.Error())
		apiError := utils.RenderAppError(ACCESS_POINT_DAO_ERROR_BASE_CODE+7, fmt.Sprintf("Could not get access points under organisation %d.Error:%s", orgId, err.Error()))
		return apiError
	}

	return nil
}

func (d *ofdiDb) UpdateFaceForFaceDataId(newFace string, faceDataId int) (appError *models.ApplicationError) {

	sqlStatement := `UPDATE organisation_face_data_ids SET "faceData" = $1 WHERE "faceDataId" = $2`

	_, err := dbPool.Exec(context.Background(), sqlStatement, newFace, faceDataId)
	if err != nil {
		logger.Log.Error(err.Error())
		appError := utils.RenderAppError(ACCESS_POINT_DAO_ERROR_BASE_CODE+7, fmt.Sprintf("Could not update faceData for FaceId %d.Error:%s", faceDataId, err.Error()))
		return appError
	}

	return nil
}

func (d *ofdiDb) UpdateAccessorForFaceDataId(accessorId int, faceDataId int) (appError *models.ApplicationError) {

	sqlStatement := `UPDATE organisation_face_data_ids SET "accessorId" = $1 WHERE "faceDataId" = $2`

	_, err := dbPool.Exec(context.Background(), sqlStatement, accessorId, faceDataId)
	if err != nil {
		logger.Log.Error(err.Error())
		appError := utils.RenderAppError(ACCESS_POINT_DAO_ERROR_BASE_CODE+7, fmt.Sprintf("Could not update faceData for FaceId %d.Error:%s", faceDataId, err.Error()))
		return appError
	}

	return nil
}

func (d *ofdiDb) GetFaceDetailsFromFace(faceData string, orgId int) (bool, *models.OrganisationFaceIdsSchema, *models.ApplicationError) {
	var response models.OrganisationFaceIdsSchema
	query := `SELECT "faceDataId", "faceData", "organisationId", "accessorId" , "pendingUnassignedOnDevices" FROM organisation_face_data_ids WHERE "faceData" = $1 and "organisationId" = $2`

	err := dbPool.QueryRow(context.Background(), query, faceData, orgId).Scan(&response.FaceDataId, &response.FaceData, &response.OrganisationId, &response.AccessorId, &response.PendingUnassignedOnDevices)
	if err != nil {
		if err == pgx.ErrNoRows {
			return false, nil, nil
		}

		logger.Log.Error(err.Error())
		appError := utils.RenderAppError(ORGANISATION_FACE_DATA_IDS_DAO_ERROR_BASE_CODE+1, fmt.Sprintf("GetFaceDetailsFromFace. Error:%s", err.Error()))
		return false, nil, appError
	}

	return true, &response, nil
}

func (d *ofdiDb) SetPendingUnassignOnDevicesFLag(faceDataId int, orgId int) *models.ApplicationError {

	sqlStatement := `UPDATE organisation_face_data_ids SET "pendingUnassignedOnDevices" = $1 WHERE "faceDataId" = $2 AND "organisationId"=$3`

	_, err := dbPool.Exec(context.Background(), sqlStatement, true, faceDataId, orgId)
	if err != nil {
		logger.Log.Error(err.Error())
		appError := utils.RenderAppError(ACCESS_POINT_DAO_ERROR_BASE_CODE+7, fmt.Sprintf("SetPendingUnassignOnDevicesFLag, faceDataId: %d.Error:%s", faceDataId, err.Error()))
		return appError
	}

	return nil
}

func (d *ofdiDb) GetFaceDetailsAssignedToAccessorOnAccessPoint(accessorId, accessPointId int) (exists bool, faceDataId int, credentialId uint32, keypadId int, appError *models.ApplicationError) {

	query := `SELECT op."faceDataId", oa."credentialId", pkm."keypadId"
			  FROM organisation_face_data_ids op 
			  INNER JOIN Face_data_id_keypad_id_mapping pkm on op."faceDataId" = pkm."faceDataId"
			  INNER JOIN organisation_accessors oa on oa."accessorId" = op."accessorId"
			  WHERE op."accessorId" = $1 and pkm."accessPointId" = $2`

	err := dbPool.QueryRow(context.Background(), query, accessorId, accessPointId).Scan(&faceDataId, &credentialId, &keypadId)
	if err != nil {
		if err == pgx.ErrNoRows {
			return false, 0, 0, 0, nil
		}

		logger.Log.Error(err.Error())
		appError := utils.RenderAppError(ORGANISATION_FACE_DATA_IDS_DAO_ERROR_BASE_CODE+1, fmt.Sprintf("GetFaceDetailsAssignedToAccessorOnAccessPoint. Error:%s", err.Error()))
		return false, 0, 0, 0, appError
	}

	return true, faceDataId, credentialId, keypadId, nil
}

func (d *ofdiDb) AddOrganisationFace(organisationId int, accessorId int, faceData string, userName string) (int, *models.ApplicationError) {

	log.Println("CP-F2")
	var faceDataId int
	sqlStatement := `INSERT INTO organisation_face_data_ids ("organisationId", "accessorId", "faceData", "userName", "pendingUnassignedOnDevices") VALUES ($1, $2, $3, $4, $5) returning "faceDataId"`

	err := dbPool.QueryRow(context.Background(), sqlStatement, organisationId, accessorId, faceData, userName, false).Scan(&faceDataId)

	if err != nil {
		errMsg := fmt.Sprintf("AddOrganisationFace:Could not create device permission.Error:%s!", err.Error())
		logger.Log.Error(errMsg)
		appError := utils.RenderAppError(ACCESS_POINT_DAO_ERROR_BASE_CODE+1, errMsg)
		return 0, appError
	}

	return faceDataId, nil
}

func (d *ofdiDb) UpdateOrganisationFace(organisationId int, accessorId int, faceData string, userName string, faceDataId int) *models.ApplicationError {

	sqlStatement := `UPDATE organisation_face_data_ids SET "faceData"=$1, "userName"=$2 WHERE "organisationId"=$3 AND "accessorId"=$4 AND "faceDataId" = $5;`

	_, err := dbPool.Exec(context.Background(), sqlStatement, faceData, userName, organisationId, accessorId, faceDataId)
	if err != nil {
		logger.Log.Error(err.Error())
		appError := utils.RenderAppError(ACCESS_POINT_DAO_ERROR_BASE_CODE+7, fmt.Sprintf("SetPendingUnassignOnDevicesFLag, faceDataId: %d.Error:%s", faceDataId, err.Error()))
		return appError
	}

	return nil
}

func (d *ofdiDb) DoesFaceExist(faceId int) (exists bool, appError *models.ApplicationError) {

	query := `SELECT EXISTS (SELECT 1 FROM organisation_face_data_ids WHERE "faceDataId" = $1)`

	err := dbPool.QueryRow(context.Background(), query, faceId).Scan(&exists)
	if err != nil {
		if err == pgx.ErrNoRows {
			return false, nil
		}

		logger.Log.Error(err.Error())
		appError = utils.RenderAppError(ACCESSOR_PERMISSION_DAO_ERROR_BASE_CODE+1, fmt.Sprintf("DoesTemplateExist. Error:%s", err.Error()))
		return false, appError
	}

	return exists, nil
}
