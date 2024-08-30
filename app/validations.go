package app

import (
	"github.com/go-playground/validator/v10"

	"face_management/utils"
)

var IsValidCardEnrollmentMode validator.Func = func(fl validator.FieldLevel) bool {

	cardMode, ok := fl.Field().Interface().(string)

	if !ok {
		// The field is not a string, so it is invalid.
		return false
	}

	switch cardMode {
	case "bulk_card_enrollment", "card_assignment":
		return true
	default:
		return false
	}
}

var isDeviceSerialNumber validator.Func = func(fl validator.FieldLevel) bool {

	value, ok := fl.Field().Interface().(string)
	if ok {

		matched, _ := utils.IsSerialNumber(value)
		return matched
	}

	return false
}

var IsValidCardCommandMsgType validator.Func = func(fl validator.FieldLevel) bool {

	msgType, ok := fl.Field().Interface().(string)

	if !ok {
		// The field is not a string, so it is invalid.
		return false
	}

	switch msgType {
	case "can_change_card_type", "change_card_type":
		return true
	default:
		return false
	}
}

var IsValidSupportCommandMsgType validator.Func = func(fl validator.FieldLevel) bool {

	msgType, ok := fl.Field().Interface().(string)

	if !ok {
		// The field is not a string, so it is invalid.
		return false
	}

	switch msgType {
	case "get_card_info", "get_accessor_permissions", "get_access_point_permissions", "get_card_assigned_to_accessor", "get_pending_permission_on_device", "get_pending_permission_of_credential", "get_device_permission_status":
		return true
	default:
		return false
	}
}
