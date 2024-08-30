package models

type AddDataToDataRepoRequest struct {
	DataType       string `json:"dataType"`
	DataId         int    `json:"dataId"`
	Data           string `json:"data"`
	NoOfKeysNeeded int    `json:"noOfKeysNeeded"`
}

type AddDataToDataRepoResponse struct {
	Type    string  `json:"type"`
	Message KeyList `json:"message"`
}
type KeyList struct {
	Keys []Keys `json:"keyList"`
}

type Keys struct {
	Id  int    `json:"id"`
	Key string `json:"Key"`
}
