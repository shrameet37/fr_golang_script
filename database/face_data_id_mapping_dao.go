package database

import (
	"context"
	"database/sql"
	"face_management/logger"
	"face_management/models"
	"face_management/utils"
	"fmt"

	"log"

	"github.com/jackc/pgx/v5"
)

type fdimDb struct{}

type fdimDbInterface interface {
	DeleteFaceFromFaceDataIdMapping(faceDataId int, orgId int) (appError *models.ApplicationError)
	DeleteFaceFromFaceDataIdMappingForAccessor(faceDataId int, orgId int, accessorId int) (appError *models.ApplicationError)
	CreateFaceDataIdMapping(faceDataId *int, accessPointId int, accessorId int) (appError *models.ApplicationError)
	// GetAssignedAccessorIdFaceIdFromKeypadIdWithAssignedAt(keypadId int, accessPointId int, accessTime int) (faceAssigned bool, accessorId int, faceId int, appError *models.ApplicationError)
	// GetUnassignedAccessorIdFaceIdFromKeypadIdWithAssignedAt(keypadId int, accessPointId int, accessTime int) (faceAssigned bool, accessorId int, faceId int, appError *models.ApplicationError)
	GetFaceIdMappingDetailsFromFaceDataId(faceDataId int) (*[]models.FaceDataIdMappingSchema, *models.ApplicationError)
	SetAssignedAtForFaceDataIdAccessPointId(assignedAt int, keypadId int, accessPointId int) *models.ApplicationError
	DeleteFromFaceDataIdMappingGivenFaceDataIdAndAccessPointId(faceDataId int, accessPointId int) (appError *models.ApplicationError)
	GetFaceIdMappingDetailsFromFaceDataIdAndAccessPointId(faceDataId, accessPointId int) (*models.FaceDataIdMappingSchema, *models.ApplicationError)
	GetFaceDataIdFromAccessorIdAccessPointId(accessorId, accessPointId int) (bool, *int, *models.ApplicationError)
	DoesFaceExistForAccessorAccessPointId(accessorId, accessPoint, faceDataId int) (bool, *models.ApplicationError)
	DeleteFromFaceDataIdMappingGivenAccessPointId(accessPointId int) (appError *models.ApplicationError)
	GetFaceDataIdsFromAccessPointId(accessPointId int) (bool, *[]int, *models.ApplicationError)
	DoesFaceExist(faceDataId int) (exists bool, appError *models.ApplicationError)
}

var FdimDb fdimDbInterface

func init() {
	FdimDb = &fdimDb{}
	// ACCESS_POINT_DAO_ERROR_BASE_CODE = utils.ServiceErrorBaseCode + ACCESS_POINT_DAO_BYTE3_VALUE
}

func (d *fdimDb) DeleteFaceFromFaceDataIdMapping(faceDataId int, orgId int) (appError *models.ApplicationError) {

	sqlStatement := `DELETE FROM face_data_id_mapping WHERE "faceDataId"=$1;`

	_, err := dbPool.Exec(context.Background(), sqlStatement, faceDataId)

	if err != nil {

		logger.Log.Error(err.Error())
		appError = utils.RenderAppError(ACCESSOR_PERMISSION_DAO_ERROR_BASE_CODE+7, fmt.Sprintf("Couldnt DeleteFaceFromFaceDataIdKeypadIdMapping .Error:%s", err.Error()))
		//TODO error handling
		return appError
	}

	return nil
}

func (d *fdimDb) CreateFaceDataIdMapping(faceDataId *int, accessPointId int, accessorId int) (appError *models.ApplicationError) {

	sqlStatement := `INSERT INTO face_data_id_mapping ("faceDataId", "accessPointId", "accessorId")
							VALUES ($1, $2, $3)`

	_, err := dbPool.Exec(context.Background(), sqlStatement, faceDataId, accessPointId, accessorId)

	if err != nil {

		logger.Log.Error(err.Error())
		appError = utils.RenderAppError(ACCESSOR_PERMISSION_DAO_ERROR_BASE_CODE+7, fmt.Sprintf("Couldnt CreateFaceDataIdKeypadIdMapping .Error:%s", err.Error()))
		return appError
	}

	return nil
}

func (d *fdimDb) DeleteFaceFromFaceDataIdMappingForAccessor(faceDataId int, orgId int, accessorId int) (appError *models.ApplicationError) {

	sqlStatement := `DELETE FROM face_data_id_mapping WHERE "faceDataId"=$1 and "accessorId"=$2;`

	_, err := dbPool.Exec(context.Background(), sqlStatement, faceDataId, accessorId)

	if err != nil {

		logger.Log.Error(err.Error())
		appError = utils.RenderAppError(ACCESSOR_PERMISSION_DAO_ERROR_BASE_CODE+7, fmt.Sprintf("Couldnt DeleteFaceFromFaceDataIdMappingForAccessor .Error:%s", err.Error()))
		//TODO error handling
		return appError
	}

	return nil
}

func (d *fdimDb) GetAssignedAccessorIdFromFaceIdWithAssignedAt(faceId int, accessPointId int, accessTime int) (faceAssigned bool, accessorId, faceDataId int, appError *models.ApplicationError) {

	sqlStatement := `SELECT "accessorId", "faceDataId" FROM face_data_id_mapping WHERE "faceDataId"=$1 AND "accessPointId"=$2 AND $3 > "assignedAt";`

	row := dbPool.QueryRow(context.Background(), sqlStatement, faceId, accessPointId, accessTime)

	err := row.Scan(&accessorId, &faceDataId)

	if err != nil {

		if err == pgx.ErrNoRows {
			return false, accessorId, faceDataId, nil
		}

		logger.Log.Error(err.Error())
		appError = utils.RenderAppError(ACCESSOR_PERMISSION_DAO_ERROR_BASE_CODE+14, fmt.Sprintf("Could not get face from keypad Id with assigned At: %d.Error:%s", faceId, err.Error()))
		return false, accessorId, faceDataId, appError
	}

	return true, accessorId, faceDataId, nil
}

func (d *fdimDb) GetUnassignedAccessorIdFaceIdFromKeypadIdWithAssignedAt(keypadId int, accessPointId int, accessTime int) (faceAssigned bool, accessorId int, faceDataId int, appError *models.ApplicationError) {

	sqlStatement := `SELECT "accessorId", "faceDataId"
		FROM "deleted_face_data_id_mapping"
		WHERE "keypadId" = $1 AND "accessPointId"=$2
		AND $3 BETWEEN "assignedAt" AND "unassignedAt";`

	row := dbPool.QueryRow(context.Background(), sqlStatement, keypadId, accessPointId, accessTime)

	err := row.Scan(&accessorId, &faceDataId)

	if err != nil {

		if err == pgx.ErrNoRows {
			return false, accessorId, faceDataId, nil
		}

		logger.Log.Error(err.Error())
		appError = utils.RenderAppError(ACCESSOR_PERMISSION_DAO_ERROR_BASE_CODE+15, fmt.Sprintf("Could not get unassigned face from keypad Id with assigned At: %d.Error:%s", keypadId, err.Error()))
		return false, accessorId, faceDataId, appError
	}

	return true, accessorId, faceDataId, nil
}

func (d *fdimDb) GetFaceIdMappingDetailsFromFaceDataId(faceDataId int) (*[]models.FaceDataIdMappingSchema, *models.ApplicationError) {

	var response []models.FaceDataIdMappingSchema
	sqlStatement := `SELECT  "accessPointId", "accessorId" FROM face_data_id_mapping WHERE "faceDataId"=$1;`

	rows, err := dbPool.Query(context.Background(), sqlStatement, faceDataId)
	if err != nil {
		appError := utils.RenderAppError(ACCESSOR_PERMISSION_DAO_ERROR_BASE_CODE+7, fmt.Sprintf("GetFaceIdMappingDetailsFromFaceDataId .Error:%s", err.Error()))
		return nil, appError
	} else {
		defer rows.Close()
		for rows.Next() {
			var faceDetails models.FaceDataIdMappingSchema
			err = rows.Scan(&faceDetails.AccessPointId, &faceDetails.AccessorId)
			if err != nil {
				appError := utils.RenderAppError(ACCESSOR_PERMISSION_DAO_ERROR_BASE_CODE+7, fmt.Sprintf("GetFaceIdMappingDetailsFromFaceDataId .Error:%s", err.Error()))
				return nil, appError
			}

			faceDetails.FaceDataId = faceDataId

			response = append(response, faceDetails)
		}

		err = rows.Err()
		if err != nil {
			appError := utils.RenderAppError(ACCESSOR_PERMISSION_DAO_ERROR_BASE_CODE+7, fmt.Sprintf("GetFaceIdMappingDetailsFromFaceDataId .Error:%s", err.Error()))
			return nil, appError
		}
	}

	return &response, nil
}

func (d *fdimDb) SetAssignedAtForFaceDataIdAccessPointId(assignedAt int, faceDataId int, accessPointId int) *models.ApplicationError {

	sqlStatement := `UPDATE face_data_id_mapping
	SET "assignedAt"=$1
	WHERE "faceDataId"=$2 AND "accessPointId" = $3`

	_, err := dbPool.Exec(context.Background(), sqlStatement, assignedAt, faceDataId, accessPointId)
	if err != nil {

		logger.Log.Error(err.Error())
		appError := utils.RenderAppError(ACCESSOR_PERMISSION_DAO_ERROR_BASE_CODE+7, fmt.Sprintf("Could not update updatedOnDevice flag for device permission, Error:%s", err.Error()))
		return appError
	}

	return nil
}

func (d *fdimDb) DeleteFromFaceDataIdMappingGivenFaceDataIdAndAccessPointId(faceDataId int, accessPointId int) (appError *models.ApplicationError) {

	sqlStatement := `DELETE FROM face_data_id_mapping WHERE "faceDataId"=$1 AND "accessPointId" = $2;`

	_, err := dbPool.Exec(context.Background(), sqlStatement, faceDataId, accessPointId)

	if err != nil {

		logger.Log.Error(err.Error())
		appError = utils.RenderAppError(ACCESSOR_PERMISSION_DAO_ERROR_BASE_CODE+7, fmt.Sprintf("Couldnt DeleteFromFaceDataIdMappingGivenFaceDataIdAndAccessPointId .Error:%s", err.Error()))
		//TODO error handling
		return appError
	}

	return nil
}

func (d *fdimDb) DeleteFromFaceDataIdMappingGivenAccessPointId(accessPointId int) (appError *models.ApplicationError) {

	sqlStatement := `DELETE FROM face_data_id_mapping WHERE "accessPointId" = $1;`

	_, err := dbPool.Exec(context.Background(), sqlStatement, accessPointId)

	if err != nil {

		logger.Log.Error(err.Error())
		appError = utils.RenderAppError(ACCESSOR_PERMISSION_DAO_ERROR_BASE_CODE+7, fmt.Sprintf("Couldnt DeleteFromFaceDataIdMappingGivenAccessPointId .Error:%s", err.Error()))
		//TODO error handling
		return appError
	}

	return nil
}

func (d *fdimDb) GetFaceIdMappingDetailsFromFaceDataIdAndAccessPointId(faceDataId, accessPointId int) (*models.FaceDataIdMappingSchema, *models.ApplicationError) {

	var response models.FaceDataIdMappingSchema
	var assignedAt sql.NullInt64

	sqlStatement := `SELECT "assignedAt", "accessorId" FROM face_data_id_mapping WHERE "faceDataId"=$1 and "accessPointId" = $2;`

	err := dbPool.QueryRow(context.Background(), sqlStatement, faceDataId, accessPointId).Scan(&assignedAt, &response.AccessorId)
	if err != nil {
		if pgx.ErrNoRows == err {
			log.Println("GetFaceIdMappingDetailsFromFaceDataIdAndAccessPointId, no rows found for faceDataId:", faceDataId, " and accessPointId:", accessPointId)
		}
		logger.Log.Error(err.Error())
		appError := utils.RenderAppError(ACCESSOR_PERMISSION_DAO_ERROR_BASE_CODE+14, fmt.Sprintf("Could not GetFaceIdMappingDetailsFromFaceDataIdAndAccessPointId:.Error:%s", err.Error()))
		return nil, appError
	}

	if assignedAt.Valid {
		response.AssignedAt = int(assignedAt.Int64)
	} else {
		response.AssignedAt = 0
	}

	response.AccessPointId = accessPointId
	response.FaceDataId = faceDataId

	return &response, nil
}

func (d *fdimDb) GetFaceDataIdFromAccessorIdAccessPointId(accessorId, accessPointId int) (bool, *int, *models.ApplicationError) {

	var response int
	sqlStatement := `SELECT "faceDataId" FROM face_data_id_mapping WHERE "accessorId"=$1 and "accessPointId" = $2;`

	row := dbPool.QueryRow(context.Background(), sqlStatement, accessorId, accessPointId)

	err := row.Scan(&response)

	if err != nil {
		logger.Log.Error(err.Error())
		appError := utils.RenderAppError(ACCESSOR_PERMISSION_DAO_ERROR_BASE_CODE+14, fmt.Sprintf("Could not GetKeypadIdDetailsFromFaceDataIdAndAccessPointId:.Error:%s", err.Error()))
		return false, nil, appError
	}

	return true, &response, nil
}

func (d *fdimDb) GetFaceDataIdsFromAccessPointId(accessPointId int) (bool, *[]int, *models.ApplicationError) {

	var faceIds []int
	exists := false
	sqlStatement := `SELECT "faceDataId" FROM face_data_id_mapping WHERE "accessPointId" = $1;`

	rows, err := dbPool.Query(context.Background(), sqlStatement, accessPointId)

	if err != nil {
		appError := utils.RenderAppError(ACCESSOR_PERMISSION_DAO_ERROR_BASE_CODE+7, fmt.Sprintf("GetFaceDataIdsFromAccessPointId .Error:%s", err.Error()))
		return false, nil, appError
	} else {
		defer rows.Close()
		var faceId int
		for rows.Next() {
			err = rows.Scan(&faceId)
			if err != nil {
				appError := utils.RenderAppError(ACCESSOR_PERMISSION_DAO_ERROR_BASE_CODE+7, fmt.Sprintf("GetFaceDataIdsFromAccessPointId .Error:%s", err.Error()))
				return false, nil, appError
			}
			exists = true
			faceIds = append(faceIds, faceId)
		}

		err = rows.Err()
		if err != nil {
			appError := utils.RenderAppError(ACCESSOR_PERMISSION_DAO_ERROR_BASE_CODE+7, fmt.Sprintf("GetFaceDataIdsFromAccessPointId .Error:%s", err.Error()))
			return false, nil, appError
		}
	}

	return exists, &faceIds, nil
}

func (d *fdimDb) GetAllAccessPointsHavingFaceDataId(faceDataId int) (*[]models.FaceDataIdMappingSchema, *models.ApplicationError) {

	var response []models.FaceDataIdMappingSchema
	sqlStatement := `SELECT "accessPointId", "accessorId" FROM face_data_id_mapping WHERE "faceDataId"=$1;`

	rows, err := dbPool.Query(context.Background(), sqlStatement, faceDataId)
	if err != nil {
		appError := utils.RenderAppError(ACCESSOR_PERMISSION_DAO_ERROR_BASE_CODE+7, fmt.Sprintf("GetKeypadIdDetailsFromFaceDataId .Error:%s", err.Error()))
		return nil, appError
	} else {
		defer rows.Close()
		for rows.Next() {
			var faceDetails models.FaceDataIdMappingSchema
			err = rows.Scan(&faceDetails.AccessPointId, &faceDetails.AccessorId)
			if err != nil {
				appError := utils.RenderAppError(ACCESSOR_PERMISSION_DAO_ERROR_BASE_CODE+7, fmt.Sprintf("GetKeypadIdDetailsFromFaceDataId .Error:%s", err.Error()))
				return nil, appError
			}

			faceDetails.FaceDataId = faceDataId

			response = append(response, faceDetails)
		}

		err = rows.Err()
		if err != nil {
			appError := utils.RenderAppError(ACCESSOR_PERMISSION_DAO_ERROR_BASE_CODE+7, fmt.Sprintf("GetKeypadIdDetailsFromFaceDataId .Error:%s", err.Error()))
			return nil, appError
		}
	}

	return &response, nil
}

func (d *fdimDb) DoesFaceExistForAccessorAccessPointId(accessorId, accessPoint, faceDataId int) (exists bool, appError *models.ApplicationError) {

	query := `SELECT EXISTS (SELECT 1 FROM face_data_id_mapping WHERE "accessorId" = $1 and "accessPointId" = $2 and "faceDataId" = $3)`

	err := dbPool.QueryRow(context.Background(), query, accessorId, accessPoint, faceDataId).Scan(&exists)
	if err != nil {
		if err == pgx.ErrNoRows {
			return false, nil
		}

		logger.Log.Error(err.Error())
		appError = utils.RenderAppError(ACCESSOR_PERMISSION_DAO_ERROR_BASE_CODE+1, fmt.Sprintf("DoesFaceExistForAccessorAccessPointId. Error:%s", err.Error()))
		return false, appError
	}

	return exists, nil
}

func (d *fdimDb) DoesFaceExist(faceDataId int) (exists bool, appError *models.ApplicationError) {

	query := `SELECT EXISTS (SELECT 1 FROM face_data_id_mapping WHERE "faceDataId" = $1)`

	err := dbPool.QueryRow(context.Background(), query, faceDataId).Scan(&exists)
	if err != nil {
		if err == pgx.ErrNoRows {
			return false, nil
		}

		logger.Log.Error(err.Error())
		appError = utils.RenderAppError(ACCESSOR_PERMISSION_DAO_ERROR_BASE_CODE+1, fmt.Sprintf("DoesFaceExistForAccessorAccessPointId. Error:%s", err.Error()))
		return false, appError
	}

	return exists, nil
}
