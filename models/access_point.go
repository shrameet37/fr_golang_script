package models

// var (
// 	// EE represents an Entry and Exit Reader configuration = 1.
// 	EE int = 1
// 	// ER represents an Entry Reader & REX configuration = 2.
// 	ER int = 2
// 	// EDC represents an Entry Reader with Door Controller configuration = 3.
// 	EDC int = 3
// 	// CR representts a Controller with REX configuration = 4. No support for card.
// 	CR int = 4
// 	// CLK represents a Clock In Device configuration = 5.
// 	CLK int = 5
// 	// DL represents a Door Lock configuration = 6.
// 	DL int = 6
// 	// EEaC represents an Entry, Exit as a Controller configuration.
// 	EEaC int = 7
// )

type AccessPoint struct {
	Id            int  `json:"id"`
	AccessPointId int  `json:"accessPointId"`
	OrgId         int  `json:"organisationId"`
	SiteId        int  `json:"siteId"`
	Configuration int  `json:"configuration"`
	ChannelNo     int  `json:"channelNo"`
	Shared        bool `json:"shared"`
}

type AccessPointDevice struct {
	Id            int    `json:"id"`
	SerialNumber  string `json:"serialNumber"`
	AccessPointId int    `json:"accessPointId"`
}

type UpdateAccessPointRequest struct {
	AccessPointId   int    `json:"accessPointId"`
	OldSerialNumber string `json:"oldSerialNumber"`
	NewSerialNumber string `json:"newSerialNumber"`
}

type CreateAccessPointRequest struct {
	OrgId              int
	AccessPointId      int
	SerialNumberMap    []interface{}
	Configuration      int
	ChannelNo          int
	SiteId             int
	InstallationMethod int
}
