package main

type Configuration struct {
	ApiUser        string
	ApiSecret      string
	ApiPasswd      string
	ApiAccessToken string
	ServerRoot     string
	InfluxServer   string
	InfluxDB       string
	InfluxUser     string
	InfluxPasswd   string
}

type accessData struct {
	Access_token  string
	Expires_in    string
	Refresh_token string
}

type zoneData struct {
	Id                            float64
	Href                          string
	ZoneOwnerId                   float64
	Name                          string
	MasterNameServers             interface{}
	ServiceTypeId                 float64
	CurrentProvisioningStatusCira string
	CurrentProvisioningStatus     string
	CurrentEventId                float64
	CustomerId                    float64
	CustomerName                  string
	Status                        string
	Conflict                      bool
}

type zoneList struct {
	HasNext     bool
	HasPrevious bool
	Page        float64
	Items       []zoneData
}
