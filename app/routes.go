package app

import (
	"face_management/middleware"
	"os"

	"github.com/gin-gonic/gin"
)

var (
	SERVICE_BASE_PATH string
	internalRoutes    *gin.RouterGroup
	integratorRoutes  *gin.RouterGroup
	supportRoutes     *gin.RouterGroup
)

func init() {

	SERVICE_BASE_PATH = os.Getenv("SERVICE_BASE_PATH")
	internalRoutes = Router.Group(SERVICE_BASE_PATH + "/internal")
	integratorRoutes = Router.Group(SERVICE_BASE_PATH)
	supportRoutes = Router.Group(SERVICE_BASE_PATH + "/support")

}

func SetupRoutesMiddleware() {

	internalRoutes.Use(middleware.LogRequest())
	internalRoutes.Use(middleware.AuthorizeApiKey(middleware.API_KEY))

	integratorRoutes.Use(middleware.LogRequest())
	integratorRoutes.Use(middleware.AuthorizeSpintlyToken())

	supportRoutes.Use(middleware.LogRequest())
	supportRoutes.Use(middleware.AuthorizeApiKey(middleware.SUPPORT_API_KEY))

}

func SetupHealthRoute() {

}

func SetupInternalRoutes() {

}

func SetupIntegratorRoutes() {

}

func SetupSupportRoutes() {

}
