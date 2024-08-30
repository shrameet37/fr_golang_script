package utils

import (
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"face_management/logger"
	"fmt"
	"os"
	"reflect"
	"regexp"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func init() {

	if _, err := os.Stat(".env"); err == nil {
		err = godotenv.Load(".env")
		if err != nil {
			panic(err)
		}
	}
}

func ShowServiceInfo() {

	serviceName := os.Getenv("SERVICE_NAME")
	fmt.Printf("Starting service %s!\n", serviceName)

	envName := os.Getenv("ENV")
	fmt.Printf("Environment is %s!\n", envName)

}

func IsSerialNumber(serialNumber string) (bool, error) {

	matched, err := regexp.MatchString("^[a-fA-F0-9]{14}$", serialNumber)
	return matched, err

}

func HasDuplicates(values []string) bool {

	valueCounts := make(map[string]int)

	for _, value := range values {
		if _, exists := valueCounts[value]; exists {
			return true
		}
		valueCounts[value] = 1
	}

	return false
}

func IsPresent(inputList []string, searchInput string) bool {

	isPresent := false

	for _, eachInput := range inputList {
		if eachInput == searchInput {
			isPresent = true
			break
		}
	}

	return isPresent
}

func Contains(slice interface{}, element interface{}) bool {

	sliceValue := reflect.ValueOf(slice)
	if sliceValue.Kind() != reflect.Slice {
		panic("Contains must be called with a slice")
	}

	for i := 0; i < sliceValue.Len(); i++ {
		currentElement := sliceValue.Index(i).Interface()
		if reflect.DeepEqual(currentElement, element) {
			return true
		}
	}

	return false
}

func ArrayToIntArray(arrayInterface []interface{}) ([]int, error) {
	intArray := make([]int, len(arrayInterface))
	for i, v := range arrayInterface {
		intValue, ok := v.(float64)
		if !ok {
			return nil, fmt.Errorf("ArrayToIntArray: element at index %d is not a valid float64", i)
		}
		intArray[i] = int(intValue)
	}
	return intArray, nil
}

func GenerateMeshPayload(msgType, msgTypeVersion byte, channelNo int, credentialId int32) string {

	var meshPayload []byte

	if channelNo != 0 {

		msgPayloadSize := 7
		meshPayload = make([]byte, msgPayloadSize)

		meshPayload[0] = msgType
		meshPayload[1] = msgTypeVersion
		meshPayload[2] = uint8(channelNo)

		binary.LittleEndian.PutUint32(meshPayload[3:7], uint32(credentialId))

	} else {

		msgPayloadSize := 6
		meshPayload = make([]byte, msgPayloadSize)

		meshPayload[0] = msgType
		meshPayload[1] = msgTypeVersion

		binary.LittleEndian.PutUint32(meshPayload[2:6], uint32(credentialId))

	}

	payloadToBeReturned := hex.EncodeToString(meshPayload)

	logger.Log.Debug("Generated Mesh Payload", zap.String("payload", payloadToBeReturned))

	return payloadToBeReturned
}

func ConvertEndianUint32(no uint32) (convNo uint32) {

	return uint32(uint8(no))<<24 | uint32(uint8(no>>8))<<16 | uint32(uint8(no>>16))<<8 | uint32(uint8(no>>24))
}

func ConvertStructToString(structInput interface{}) string {

	structInputByteArray, err := json.Marshal(structInput)
	if err != nil {
		errorMsg := fmt.Sprintf("Could not convert struct input to byte array.Struct value:%v.Error message:%s!", structInput, err.Error())
		logger.Log.Error(errorMsg)
		return ""
	}

	return string(structInputByteArray)

}
