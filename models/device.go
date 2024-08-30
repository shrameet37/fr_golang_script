package models

var (
	ENTRY                    int = 1
	EXIT                     int = 2
	ONE_DOOR_CONTROLLER      int = 3
	MULTI_CHANNEL_CONTROLLER int = 4
)

type Device struct {
	Id           int    `json:"id"`
	DeviceType   int    `json:"deviceType"`
	SerialNumber string `json:"serialNumber"`
	OrgId        int    `json:"organisationId"`
}
type CreateDeviceRequest struct {
	SerialNumber string
	DeviceType   int
	OrgId        int
}
