package database

import (
	"context"
	"face_management/logger"
	"face_management/models"
	"face_management/utils"
	"fmt"
)

type devDb struct{}

type devDbInterface interface {
	CreateDevice(device models.Device) *models.ApplicationError
	UpdateDevice(deviceType, accessPointId int, serialNumber string) (appError *models.ApplicationError)
	DeleteDevice(serialNumber string, orgId int) *models.ApplicationError
	DeleteDeviceV2(serialNumber string) *models.ApplicationError
	DeleteDeviceExceptController(serialNumber string, orgId int) (appError *models.ApplicationError)

	GetConfigurationAndDeviceTypeOfDevice(serialNumber string) (int, int, *models.ApplicationError)
}

var DevDb devDbInterface

const (
	DEVICE_DAO_BYTE3_VALUE = 0x00330000 //3 = dao, 3 = device
)

var (
	DEVICE_DAO_ERROR_BASE_CODE int
)

func init() {
	DevDb = &devDb{}
	DEVICE_DAO_ERROR_BASE_CODE = utils.ServiceErrorBaseCode + DEVICE_DAO_BYTE3_VALUE

}

func (d *devDb) CreateDevice(device models.Device) *models.ApplicationError {

	sqlStatement := `INSERT INTO devices ("deviceType", "serialNumber", "organisationId") VALUES ($1, $2, $3);`

	_, err := dbPool.Exec(context.Background(), sqlStatement, device.DeviceType, device.SerialNumber, device.OrgId)

	if err != nil {
		logger.Log.Error(err.Error())
		appError := utils.RenderAppError(DEVICE_DAO_ERROR_BASE_CODE+1, fmt.Sprintf("Could not create device.Error:%s", err.Error()))
		return appError
	}

	return nil

}

func (d *devDb) DeleteDevice(serialNumber string, orgId int) *models.ApplicationError {

	sqlStatement := `delete from devices where "serialNumber"=$1;`

	_, err := dbPool.Exec(context.Background(), sqlStatement, serialNumber)

	if err != nil {
		logger.Log.Error(err.Error())
		apiError := utils.RenderAppError(DEVICE_DAO_ERROR_BASE_CODE+2, fmt.Sprintf("Could not delete device %s.Error:%s", serialNumber, err.Error()))
		return apiError
	}

	return nil

}

func (d *devDb) DeleteDeviceExceptController(serialNumber string, orgId int) (appError *models.ApplicationError) {

	sqlStatement := `DELETE FROM devices WHERE "serialNumber"=$1 AND "deviceType" NOT IN (3,4)`

	_, err := dbPool.Exec(context.Background(), sqlStatement, serialNumber)

	if err != nil {

		logger.Log.Error(err.Error())
		appError := utils.RenderAppError(DEVICE_DAO_ERROR_BASE_CODE+3, fmt.Sprintf("Couldnt delete device except controller for device serialNumber: %s, Error :%s", serialNumber, err.Error()))
		return appError
	}

	return nil
}

func (d *devDb) DeleteDeviceV2(serialNumber string) *models.ApplicationError {

	sqlStatement := `delete from devices where "serialNumber"=$1;`

	_, err := dbPool.Exec(context.Background(), sqlStatement, serialNumber)

	if err != nil {
		logger.Log.Error(err.Error())
		apiError := utils.RenderAppError(DEVICE_DAO_ERROR_BASE_CODE+4, fmt.Sprintf("Could not delete device %s.Error:%s", serialNumber, err.Error()))
		return apiError
	}

	return nil

}

func (d *devDb) UpdateDevice(deviceType, accessPointId int, serialNumber string) (appError *models.ApplicationError) {

	sqlStatement := `UPDATE devices SET "deviceType" = $1, "accessPointId" = $2 WHERE "serialNumber" = $3;`

	_, err := dbPool.Exec(context.Background(), sqlStatement, deviceType, accessPointId, serialNumber)

	if err != nil {
		logger.Log.Error(err.Error())
		appError = utils.RenderAppError(DEVICE_DAO_ERROR_BASE_CODE+5, fmt.Sprintf("Could not update device %s.Error:%s", serialNumber, err.Error()))
		return appError
	}

	return nil
}

func (d *devDb) GetConfigurationAndDeviceTypeOfDevice(serialNumber string) (configuration int, deviceType int, appError *models.ApplicationError) {

	sqlStatement := `SELECT
							ap."configuration",
							d."deviceType"
						FROM
							access_point_devices apd
						INNER JOIN
							devices d ON apd."serialNumber" = d."serialNumber"
						INNER JOIN
							access_points ap ON ap."accessPointId" = apd."accessPointId"
						WHERE
							d."serialNumber" = $1;`

	err := dbPool.QueryRow(context.Background(), sqlStatement, serialNumber).Scan(&configuration, &deviceType)

	if err != nil {
		logger.Log.Error(err.Error())
		appError = utils.RenderAppError(ACCESS_POINT_DAO_ERROR_BASE_CODE+24, fmt.Sprintf("Could not get configuration and channel no from device serial number: %s.Error:%s", serialNumber, err.Error()))
		return -1, -1, appError
	}

	return configuration, deviceType, nil
}
