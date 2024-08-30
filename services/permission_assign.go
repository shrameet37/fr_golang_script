package services

import (
	"encoding/csv"
	"encoding/json"
	"face_management/database"
	"face_management/interactors"
	"face_management/models"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/google/uuid"
)

func AssignPermsInExcelR(accessPointId int) {

	filename := fmt.Sprintf("SPINTLY_AP_%d.csv", accessPointId)
	file, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// Initialize the CSV writer
	writer := csv.NewWriter(file)
	defer writer.Flush()

	dpDetailsAll, appError := database.DevPermDb.ExcelrGetDevicePermissionFromAccessPointId(accessPointId)
	if appError != nil {
		log.Panic("CP1")
	}

	for _, dpDetails := range dpDetailsAll {

		log.Println("AccessorId = ", dpDetails.AccessorId)
		log.Println("FaceDataId = ", dpDetails.FaceDataId)
		log.Println("SerialNumber = ", dpDetails.SerialNumber)
		log.Println("Username = ", dpDetails.Username)

		DataRepoRequest := models.AddDataToDataRepoRequest{
			DataType:       "FR",
			DataId:         dpDetails.FaceDataId,
			Data:           string(dpDetails.FaceData),
			NoOfKeysNeeded: 1,
		}

		dataRepoRespone, appError := interactors.AddDataToDataRepo(&DataRepoRequest)
		if appError != nil {
			errMsg := fmt.Sprintf("processDevicePermission->GetFaceDataIdDetails: Couldn't GetFaceDataIdDetails fro face Id %d ", dpDetails.FaceDataId)
			log.Panic(errMsg)
		}

		for _, dataRepoKey := range dataRepoRespone.Message.Keys {

			iotEngineCommand := models.IotEngineCommandMsg{
				Version:        2,
				SrcAppID:       0x81,
				DestAppID:      models.GATEWAY_FACE_MANAGER_APP_ID,
				IsLiveMsg:      false,
				Target:         "device",
				TargetSerialNo: dpDetails.SerialNumber,
				MsgType:        "add_user_face",
				MsgTypeVer:     1,
				UserId:         uint32(dpDetails.AccessorId),
				UserName:       dpDetails.Username,
				KeyId:          uint32(dataRepoKey.Id),
				ValidationKey:  dataRepoKey.Key,
				MessageId:      uuid.New(),
			}

			msgIot, err := json.Marshal(iotEngineCommand)
			if err != nil {
				log.Panic("CP2")
			}
			log.Println("Msg to send", iotEngineCommand) //todo
			log.Println("Msg to send", msgIot)           //todo
			log.Println("Msg to send", string(msgIot))   //todo

			appError = interactors.SendMsg(iotEngineCommand)

			if appError != nil {
				log.Panic(appError)
			}

			record := []string{strconv.Itoa(int(uint32(dataRepoKey.Id)))}

			// Write the record to the CSV
			if err := writer.Write(record); err != nil {
				panic(err)
			}

			// Make sure all data is written to the file
			writer.Flush()

			// Check if there were any errors during the writing process
			if err := writer.Error(); err != nil {
				panic(err)
			}

		}
	}

}
