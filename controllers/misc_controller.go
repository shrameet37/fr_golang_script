package controllers

import (
	"face_management/utils"
)

const (
	MISC_CONTROLLER_DAO_BYTE3_VALUE       = 0x00130000 //1 = controller, 3 = Misc
	CONTROLLER_DOESNT_NOT_HAVE_PERMISSION = 0x00000600
)

var (
	MISC_CONTROLLER_ERROR_BASE_CODE int
)

func init() {
	MISC_CONTROLLER_ERROR_BASE_CODE = utils.ServiceErrorBaseCode + MISC_CONTROLLER_DAO_BYTE3_VALUE
}
