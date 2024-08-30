package database

import (
	"context"
	"face_management/logger"
	"face_management/models"
	"face_management/utils"
	"fmt"
)

type dfdimDb struct{}

type dfdimDbInterface interface {
	CreateDeletedFaceDataIdMapping(faceDataId *int, accessPointId int, accessorId int, assignedAt int, unassignedAt int) (appError *models.ApplicationError)
}

var DfdimDb dfdimDbInterface

func init() {
	DfdimDb = &dfdimDb{}
}

func (d *dfdimDb) CreateDeletedFaceDataIdMapping(faceDataId *int, accessPointId int, accessorId int, assignedAt int, unassignedAt int) (appError *models.ApplicationError) {

	sqlStatement := `INSERT INTO deleted_face_data_id_mapping ("faceDataId", "accessPointId", "accessorId", "assignedAt", "unAssignedAt")
							VALUES ($1, $2, $3, $4, $5)`

	_, err := dbPool.Exec(context.Background(), sqlStatement, faceDataId, accessPointId, accessorId, assignedAt, unassignedAt)

	if err != nil {

		logger.Log.Error(err.Error())
		appError = utils.RenderAppError(ACCESSOR_PERMISSION_DAO_ERROR_BASE_CODE+7, fmt.Sprintf("Couldnt CreateDeletedFaceDataIdMapping .Error:%s", err.Error()))
		return appError
	}

	return nil
}
