package main

import (
	//please leave this space so that utils package gets imported first. This is needed so that env variables get loaded first!

	"face_management/database"
	"face_management/services"
)

func main() {

	defer func() {
		database.CloseDatabasePool()
	}()

	if err := database.InitializeDatabasePool(); err != nil {
		panic(err)
	}
	services.AssignPermsInExcelR(9467)
}
