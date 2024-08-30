package models

type FaceEventInternalAccess struct {
	Version            int    `json:"version"`            //
	EventType          string `json:"eventType"`          // fingerprint_access
	EventTime          uint32 `json:"eventTime"`          //
	DeviceSerialNumber string `json:"deviceSerialNumber"` //
	AccessorID         int    `json:"accessorId"`         //
}
