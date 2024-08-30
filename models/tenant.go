package models

type SharedAccessPointRequest struct {
	LandlordOrgId     int                 `json:"landlordOrgId"`
	LinkedTenantSites []LinkedTenantSites `json:"linkedTenantSites"`
}

type LinkedTenantSites struct {
	TenantOrgId              int                        `json:"tenantOrgId"`
	TenantSiteId             int                        `json:"tenantSiteId"`
	SharedAccessPoints       []int                      `json:"sharedAccessPoints"`
	SharedAccessPointDevices []SharedAccessPointDevices `json:"-"`
}

type SharedAccessPointDevices struct {
	Id      int      `json:"id"`
	Devices []Device `json:"-"`
}
