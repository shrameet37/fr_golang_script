package database

import (
	"context"
	"face_management/logger"
	"face_management/models"
	"face_management/utils"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/lib/pq"
)

type apDb struct{}

type apDbInterface interface {
	CreateAccessPoint(accessPoint models.AccessPoint) *models.ApplicationError
	CreateAccessPointAndDevice(apd models.AccessPointDevice) *models.ApplicationError
	GetOrganisationOfAccessPoint(accessPointId int) (int, *models.ApplicationError)
	GetAccessPointsUnderOrganisationCount(orgId int) (int, *models.ApplicationError)
	DeleteAccessPoint(accessPointId, orgId int) *models.ApplicationError
	DeleteAccessPointDevices(accessPointId int) *models.ApplicationError
	GetDevicesUnderAccessPoint(accessPointId int) (devices []string, apiError *models.ApplicationError)
	GetControllerCountUnderAccessPoint(serialNumber string, orgId int) (controllerCount int, appError *models.ApplicationError)
	GetConfigurationAndChannelNoOfAccessPoint(accessPointId int) (int, int, *models.ApplicationError)
	GetDevicesOfAccessPoint(accessPointId int) ([]models.Device, *models.ApplicationError)
	GetControllerOfAccessPoint(accessPointId int) ([]models.Device, *models.ApplicationError)
	GetExitDeviceOfAccessPoint(accessPointId int) ([]models.Device, *models.ApplicationError)
	GetAccessPointsUnderOrganisation(orgId int) ([]models.AccessPoint, *models.ApplicationError)
	GetAccessPointIdsUnderOrganisation(orgId int) ([]models.AccessPoint, *models.ApplicationError)
	GetExistingAccessPoints(accessPointIds []int, orgId int) ([]int, *models.ApplicationError)
	GetAllDevicesWithControllerOfAccessPoint(accessPointId, orgId int) ([]models.Device, *models.ApplicationError)
	GetAccessPointIdFromDeviceSerialNumber(serialNumber string) (apId int, appError *models.ApplicationError)
	GetAccessPointIdsFromDeviceSerialNumber(serialNumber string) (apIds []int, appError *models.ApplicationError)
	UpdateAccessPointDevice(apd models.UpdateAccessPointRequest) (appError *models.ApplicationError)
}

var ApDb apDbInterface

const (
	ACCESS_POINT_DAO_BYTE3_VALUE = 0x00310000 //3 = dao, 1 = access point
)

var (
	ACCESS_POINT_DAO_ERROR_BASE_CODE int
)

func init() {
	ApDb = &apDb{}
	ACCESS_POINT_DAO_ERROR_BASE_CODE = utils.ServiceErrorBaseCode + ACCESS_POINT_DAO_BYTE3_VALUE
}

func (d *apDb) GetOrganisationOfAccessPoint(accessPointId int) (int, *models.ApplicationError) {

	var orgId int

	sqlStatement := `SELECT "organisationId" FROM access_points WHERE "accessPointId"=$1;`

	err := dbPool.QueryRow(context.Background(), sqlStatement, accessPointId).Scan(&orgId)

	if err != nil {

		logger.Log.Error(err.Error())
		appError := utils.RenderAppError(ACCESS_POINT_DAO_ERROR_BASE_CODE+1, fmt.Sprintf("Could not get orgId of access point.Error:%s", err.Error()))
		return 0, appError
	}

	return orgId, nil
}

func (d *apDb) CreateAccessPoint(accessPoint models.AccessPoint) *models.ApplicationError {

	sqlStatement := `INSERT INTO access_points ("accessPointId", "organisationId", "siteId", "configuration", "channelNo")
							VALUES ($1, $2, $3, $4, $5)`

	_, err := dbPool.Exec(context.Background(), sqlStatement, accessPoint.AccessPointId, accessPoint.OrgId, accessPoint.SiteId, accessPoint.Configuration, accessPoint.ChannelNo)

	if err != nil {
		logger.Log.Error(err.Error())
		appError := utils.RenderAppError(ACCESS_POINT_DAO_ERROR_BASE_CODE+2, fmt.Sprintf("Could not create access point.Error:%s", err.Error()))
		return appError
	}

	return nil

}

func (d *apDb) CreateAccessPointAndDevice(apd models.AccessPointDevice) *models.ApplicationError {

	sqlStatement := `INSERT INTO access_point_devices ("accessPointId", "serialNumber")
							VALUES ($1, $2)`
	_, err := dbPool.Exec(context.Background(), sqlStatement, apd.AccessPointId, apd.SerialNumber)

	if err != nil {
		logger.Log.Error(err.Error())
		appError := utils.RenderAppError(ACCESS_POINT_DAO_ERROR_BASE_CODE+3, fmt.Sprintf("Could not create access point and device in access_points_devices table.Error:%s", err.Error()))
		return appError
	}

	return nil

}

func (d *apDb) DeleteAccessPoint(accessPointId, orgId int) *models.ApplicationError {

	sqlStatement := `DELETE FROM access_points WHERE "accessPointId"=$1 AND "organisationId"=$2`

	_, err := dbPool.Exec(context.Background(), sqlStatement, accessPointId, orgId)

	if err != nil {
		logger.Log.Error(err.Error())
		appError := utils.RenderAppError(ACCESS_POINT_DAO_ERROR_BASE_CODE+4, fmt.Sprintf("Could not delete access_point %d.Error:%s", accessPointId, err.Error()))
		return appError
	}

	return nil

}

func (d *apDb) DeleteAccessPointDevices(accessPointId int) *models.ApplicationError {

	sqlStatement := `DELETE from access_point_devices where "accessPointId"=$1`

	_, err := dbPool.Exec(context.Background(), sqlStatement, accessPointId)

	if err != nil {
		logger.Log.Error(err.Error())
		appError := utils.RenderAppError(ACCESS_POINT_DAO_ERROR_BASE_CODE+5, fmt.Sprintf("Could not delete access_point %d.Error:%s", accessPointId, err.Error()))
		return appError
	}

	return nil

}

func (d *apDb) GetConfigurationAndChannelNoOfAccessPoint(accessPointId int) (configuration int, channelNo int, appError *models.ApplicationError) {

	sqlStatement := `SELECT "configuration", "channelNo" FROM access_points WHERE "accessPointId"=$1;`

	err := dbPool.QueryRow(context.Background(), sqlStatement, accessPointId).Scan(&configuration, &channelNo)

	if err != nil {
		errMsg := fmt.Sprintf("GetConfigurationAndChannelNoOfAccessPoint:Could not get configuration and channel no of access point %d.Error:%s!", accessPointId, err.Error())
		logger.Log.Error(errMsg)
		appError = utils.RenderAppError(ACCESS_POINT_DAO_ERROR_BASE_CODE+6, errMsg)
		return 0, 0, appError
	}

	return configuration, channelNo, nil
}

func (d *apDb) GetAccessPointsUnderOrganisation(orgId int) ([]models.AccessPoint, *models.ApplicationError) {

	sqlStatement := `SELECT "accessPointId", "configuration", "channelNo" FROM access_points WHERE "organisationId" = $1`

	rows, err := dbPool.Query(context.Background(), sqlStatement, orgId)

	if err != nil {
		logger.Log.Error(err.Error())
		apiError := utils.RenderAppError(ACCESS_POINT_DAO_ERROR_BASE_CODE+7, fmt.Sprintf("Could not get access points under organisation %d.Error:%s", orgId, err.Error()))
		return nil, apiError
	}

	defer rows.Close()

	accessPoints := []models.AccessPoint{}

	for rows.Next() {

		ap := models.AccessPoint{}
		err := rows.Scan(&ap.AccessPointId, &ap.Configuration, &ap.ChannelNo)
		if err != nil {
			logger.Log.Error(err.Error())
			apiError := utils.RenderAppError(ACCESS_POINT_DAO_ERROR_BASE_CODE+8, fmt.Sprintf("Error:%s", err.Error()))
			return nil, apiError
		}

		accessPoints = append(accessPoints, ap)

	}

	return accessPoints, nil

}

func (d *apDb) GetAccessPointIdsUnderOrganisation(orgId int) ([]models.AccessPoint, *models.ApplicationError) {

	sqlStatement := `SELECT "accessPointId" FROM access_points WHERE "organisationId" = $1`

	rows, err := dbPool.Query(context.Background(), sqlStatement, orgId)

	if err != nil {
		logger.Log.Error(err.Error())
		apiError := utils.RenderAppError(ACCESS_POINT_DAO_ERROR_BASE_CODE+9, fmt.Sprintf("Could not get access points under organisation %d.Error:%s", orgId, err.Error()))
		return nil, apiError
	}

	defer rows.Close()

	accessPoints := []models.AccessPoint{}

	for rows.Next() {

		ap := models.AccessPoint{}
		err := rows.Scan(&ap.AccessPointId)
		if err != nil {
			logger.Log.Error(err.Error())
			apiError := utils.RenderAppError(ACCESS_POINT_DAO_ERROR_BASE_CODE+10, fmt.Sprintf("Error:%s", err.Error()))
			return nil, apiError
		}

		accessPoints = append(accessPoints, ap)

	}

	return accessPoints, nil

}

func (d *apDb) GetDevicesOfAccessPoint(accessPointId int) ([]models.Device, *models.ApplicationError) {

	sqlStatement := `SELECT d."id", d."deviceType", d."serialNumber", d."organisationId" FROM devices d
						INNER JOIN access_point_devices apd ON apd."serialNumber" = d."serialNumber"
						WHERE apd."accessPointId"=$1 AND d."deviceType" IN (1, 2)`

	rows, err := dbPool.Query(context.Background(), sqlStatement, accessPointId)

	if err != nil {
		errMsg := fmt.Sprintf("GetDevicesOfAccessPoint:Could not get devices from access_points_devices for access point %d.Error:%s!", accessPointId, err.Error())
		logger.Log.Error(errMsg)
		appError := utils.RenderAppError(ACCESS_POINT_DAO_ERROR_BASE_CODE+11, errMsg)
		return nil, appError
	}

	defer rows.Close()

	devices := []models.Device{}

	for rows.Next() {
		device := models.Device{}
		err = rows.Scan(&device.Id, &device.DeviceType, &device.SerialNumber, &device.OrgId)
		if err != nil {
			errMsg := fmt.Sprintf("GetDevicesOfAccessPoint:Could not get build array of devices for access point %d.Error:%s!", accessPointId, err.Error())
			logger.Log.Error(errMsg)
			appError := utils.RenderAppError(ACCESS_POINT_DAO_ERROR_BASE_CODE+12, errMsg)
			return nil, appError
		}
		devices = append(devices, device)
	}

	return devices, nil

}

func (d *apDb) GetAllDevicesWithControllerOfAccessPoint(accessPointId, orgId int) ([]models.Device, *models.ApplicationError) {

	sqlStatement := `SELECT d."id", d."deviceType", d."serialNumber", d."organisationId"
							FROM devices d
							INNER JOIN access_point_devices apd ON apd."serialNumber" = d."serialNumber"
							WHERE apd."accessPointId"=$1 AND d."organisationId"=$2;
							`

	rows, err := dbPool.Query(context.Background(), sqlStatement, accessPointId, orgId)

	if err != nil {
		logger.Log.Error(err.Error())
		appError := utils.RenderAppError(ACCESS_POINT_DAO_ERROR_BASE_CODE+13, fmt.Sprintf("Could not get devices from access_points_devices for access point %d.Error:%s", accessPointId, err.Error()))
		return nil, appError
	}

	defer rows.Close()

	devices := []models.Device{}

	for rows.Next() {
		device := models.Device{}
		err = rows.Scan(&device.Id, &device.DeviceType, &device.SerialNumber, &device.OrgId)
		if err != nil {
			logger.Log.Error(err.Error())
			appError := utils.RenderAppError(ACCESS_POINT_DAO_ERROR_BASE_CODE+14, fmt.Sprintf("Error:%s", err.Error()))
			return nil, appError
		}
		devices = append(devices, device)
	}

	return devices, nil

}

func (d *apDb) GetControllerOfAccessPoint(accessPointId int) ([]models.Device, *models.ApplicationError) {

	sqlStatement := `SELECT d."id", d."deviceType", d."serialNumber", d."organisationId" FROM devices d
						INNER JOIN access_point_devices apd ON apd."serialNumber" = d."serialNumber"
						WHERE apd."accessPointId"=$1 AND d."deviceType" IN (3, 4)`

	rows, err := dbPool.Query(context.Background(), sqlStatement, accessPointId)

	if err != nil {
		errMsg := fmt.Sprintf("GetControllerOfAccessPoint:Could not get controller from access_points_devices for access point %d.Error:%s!", accessPointId, err.Error())
		logger.Log.Error(errMsg)
		appError := utils.RenderAppError(ACCESS_POINT_DAO_ERROR_BASE_CODE+15, errMsg)
		return nil, appError
	}

	defer rows.Close()

	devices := []models.Device{}

	for rows.Next() {
		device := models.Device{}
		err = rows.Scan(&device.Id, &device.DeviceType, &device.SerialNumber, &device.OrgId)
		if err != nil {
			errMsg := fmt.Sprintf("GetControllerOfAccessPoint:Could not get build array of controllers for access point %d.Error:%s!", accessPointId, err.Error())
			logger.Log.Error(errMsg)
			appError := utils.RenderAppError(ACCESS_POINT_DAO_ERROR_BASE_CODE+16, errMsg)
			return nil, appError
		}
		devices = append(devices, device)
	}

	return devices, nil

}

func (d *apDb) GetExitDeviceOfAccessPoint(accessPointId int) ([]models.Device, *models.ApplicationError) {

	sqlStatement := `SELECT d."id", d."deviceType", d."serialNumber", d."organisationId" FROM devices d
						INNER JOIN access_point_devices apd ON apd."serialNumber" = d."serialNumber"
						WHERE apd."accessPointId"=$1 AND d."deviceType" = 2`

	rows, err := dbPool.Query(context.Background(), sqlStatement, accessPointId)

	if err != nil {
		errMsg := fmt.Sprintf("GetExitDeviceOfAccessPoint:Could not get exit device from access_points_devices for access point %d.Error:%s!", accessPointId, err.Error())
		logger.Log.Error(errMsg)
		appError := utils.RenderAppError(ACCESS_POINT_DAO_ERROR_BASE_CODE+17, errMsg)
		return nil, appError
	}

	defer rows.Close()

	devices := []models.Device{}

	for rows.Next() {
		device := models.Device{}
		err = rows.Scan(&device.Id, &device.DeviceType, &device.SerialNumber, &device.OrgId)
		if err != nil {
			errMsg := fmt.Sprintf("GetExitDeviceOfAccessPoint:Could not get build array of exit devices for access point %d.Error:%s!", accessPointId, err.Error())
			logger.Log.Error(errMsg)
			appError := utils.RenderAppError(ACCESS_POINT_DAO_ERROR_BASE_CODE+18, errMsg)
			return nil, appError
		}
		devices = append(devices, device)
	}

	return devices, nil

}

func (d *apDb) GetAccessPointsUnderOrganisationCount(orgId int) (int, *models.ApplicationError) {

	var count int

	sqlStatement := `SELECT COUNT(*) FROM access_points WHERE "organisationId" =$1`

	err := dbPool.QueryRow(context.Background(), sqlStatement, orgId).Scan(&count)

	if err != nil {
		logger.Log.Error(err.Error())
		apiError := utils.RenderAppError(ACCESS_POINT_DAO_ERROR_BASE_CODE+19, fmt.Sprintf("Could not get access_point count for org %d.Error:%s", orgId, err.Error()))
		return -1, apiError
	}

	return count, nil

}

func (d *apDb) GetDevicesUnderAccessPoint(accessPointId int) (serialNumber []string, appError *models.ApplicationError) {

	sqlStatement := `SELECT apd."serialNumber"
							FROM access_point_devices apd
							INNER JOIN devices d ON d."serialNumber" = apd."serialNumber"
							WHERE apd."accessPointId" = $1
							GROUP BY apd."serialNumber";`

	rows, err := dbPool.Query(context.Background(), sqlStatement, accessPointId)

	if err != nil {

		appError := utils.RenderAppError(ACCESS_POINT_DAO_ERROR_BASE_CODE+20, fmt.Sprintf("Could not get devices and details for access points with id %d.Error:%s", accessPointId, err.Error()))
		return nil, appError
	}

	defer rows.Close()

	deviceCount := 0

	var deviceSerialNumber string

	for rows.Next() {
		err := rows.Scan(&deviceSerialNumber)
		if err != nil {
			logger.Log.Error(err.Error())
			continue
		}
		serialNumber = append(serialNumber, deviceSerialNumber)
		deviceCount++
	}

	return serialNumber, nil
}

func (d *apDb) GetControllerCountUnderAccessPoint(serialNumber string, orgId int) (controllerCount int, appError *models.ApplicationError) {

	sqlStatement := `select COUNT(*) as "controllerCount"
								from devices d 
								inner join access_point_devices apd 
								on apd."serialNumber" = d."serialNumber" 
								where d."deviceType" = 4 
								and d."serialNumber" = $1 and d."organisationId" = $2;`

	row := dbPool.QueryRow(context.Background(), sqlStatement, serialNumber, orgId)

	err := row.Scan(&controllerCount)

	if err != nil {

		if err == pgx.ErrNoRows {

			controllerCount = 0
			return controllerCount, nil
		}

		appError := utils.RenderAppError(ACCESS_POINT_DAO_ERROR_BASE_CODE+21, fmt.Sprintf("Could not get devices and details for access points with serialId %s.Error:%s", serialNumber, err.Error()))

		return controllerCount, appError
	}

	return controllerCount, nil

}

func (d *apDb) GetExistingAccessPoints(accessPointIds []int, orgId int) ([]int, *models.ApplicationError) {

	if len(accessPointIds) == 0 {
		return nil, nil
	}

	sqlStatement := `SELECT "accessPointId" FROM access_points WHERE "accessPointId" = ANY($1) AND "organisationId" = $2;`

	rows, err := dbPool.Query(context.Background(), sqlStatement, pq.Array(accessPointIds), orgId)
	if err != nil {
		logger.Log.Error(err.Error())
		appError := utils.RenderAppError(ACCESS_POINT_DAO_ERROR_BASE_CODE+22, fmt.Sprintf("Could not verify if access points exist. Error: %s", err.Error()))
		return nil, appError
	}

	defer rows.Close()

	var apIds []int
	for rows.Next() {
		var apId int
		if err := rows.Scan(&apId); err != nil {
			logger.Log.Error(err.Error())
			appError := utils.RenderAppError(ACCESS_POINT_DAO_ERROR_BASE_CODE+23, fmt.Sprintf("Error: %s", err.Error()))
			return nil, appError
		}
		apIds = append(apIds, apId)
	}

	return apIds, nil
}

func (d *apDb) GetAccessPointIdFromDeviceSerialNumber(serialNumber string) (apId int, appError *models.ApplicationError) {

	sqlStatement := `SELECT ap."accessPointId" FROM access_points ap INNER JOIN access_point_devices apd ON ap."accessPointId" = apd."accessPointId" WHERE apd."serialNumber"=$1`

	row := dbPool.QueryRow(context.Background(), sqlStatement, serialNumber)

	err := row.Scan(&apId)

	if err != nil {
		logger.Log.Error(err.Error())
		appError = utils.RenderAppError(ACCESS_POINT_DAO_ERROR_BASE_CODE+24, fmt.Sprintf("Could not get access point Id from device serial number: %s.Error:%s", serialNumber, err.Error()))
		return -1, appError
	}

	return apId, nil
}

func (d *apDb) GetAccessPointIdsFromDeviceSerialNumber(serialNumber string) (apIds []int, appError *models.ApplicationError) {

	sqlStatement := `SELECT ap."accessPointId" FROM access_points ap INNER JOIN access_point_devices apd ON ap."accessPointId" = apd."accessPointId" WHERE apd."serialNumber"=$1`

	rows, err := dbPool.Query(context.Background(), sqlStatement, serialNumber)

	if err != nil {

		if err == pgx.ErrNoRows {
			return apIds, nil
		}

		appError = utils.RenderAppError(ACCESS_POINT_DAO_ERROR_BASE_CODE+25, fmt.Sprintf("Could not get access point Ids from device serial number: %s.Error:%s", serialNumber, err.Error()))
		return apIds, appError
	}

	defer rows.Close()

	for rows.Next() {

		var apId int
		err = rows.Scan(&apId)

		if err != nil {
			appError = utils.RenderAppError(ACCESS_POINT_DAO_ERROR_BASE_CODE+26, fmt.Sprintf("Could not get access point Ids from device serial number: %s.Error:%s", serialNumber, err.Error()))
			return apIds, appError
		}

		apIds = append(apIds, apId)

	}

	return apIds, nil
}

func (d *apDb) UpdateAccessPointDevice(apd models.UpdateAccessPointRequest) (appError *models.ApplicationError) {

	sqlStatement := `UPDATE access_point_devices 
							SET "serialNumber" = $1
							WHERE "serialNumber" = $2 AND "accessPointId" = $3`

	_, err := dbPool.Exec(context.Background(), sqlStatement, apd.NewSerialNumber, apd.OldSerialNumber, apd.AccessPointId)
	if err != nil {
		logger.Log.Error(err.Error())
		appError := utils.RenderAppError(ACCESS_POINT_DAO_ERROR_BASE_CODE+27, fmt.Sprintf("Could not update access point Id for new device serial number: %s with old serial number: %s.Error:%s", apd.NewSerialNumber, apd.OldSerialNumber, err.Error()))
		return appError
	}

	return nil
}
