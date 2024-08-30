package models

type ExclerDevicePermission struct {
	SerialNumber string
	AccessorId   int
	FaceDataId   int
	FaceData     []byte
	Username     string
}
