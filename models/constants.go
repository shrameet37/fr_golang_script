package models

const (
	ASSIGN_FACE       int = 1
	UNASSIGN_FACE     int = 2
	ADD_PERMISSION    int = 3
	REMOVE_PERMISSION int = 4
	UPDATE_FACE       int = 5
)

const (
	GATEWAY_FACE_MANAGER_APP_ID = 0x12

	CLOUD_FACE_MANAGEMENT_APP_ID = 0x92
)

type DEVICE_EVENT_ACCESS_CONTROL byte

const (
	//IMAGEMATCH
	DEV_EVT_ACCESS_CONTROL_FINGERPRINT_ACCESS DEVICE_EVENT_ACCESS_CONTROL = 0x05
	// DEV_EVT_ACCESS_CONTROL_NO_PERMISSION_KEYPAD_ACCESS DEVICE_EVENT_ACCESS_CONTROL = 0x49

	// DEV_EVT_AC_ACCESS_DENIED_DOOR_LOCKED                        DEVICE_EVENT_ACCESS_CONTROL = 0x0B
	DEV_EVT_ACCESS_CONTROL_ACCESS_DENIED_USER_SCHEDULE_DISABLED DEVICE_EVENT_ACCESS_CONTROL = 0x29
)
