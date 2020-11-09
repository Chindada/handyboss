package models

// Dc Dc
type Dc struct {
	ID              int64  `json:"id"`
	DcAuthorization string `json:"dcAuthorization"`
	PutTimeInterval int64  `json:"putTimeInterval"`
	IdleTime        int64  `json:"idleTime"`
	CreateTime      int64  `json:"createTime"`
	MacAddress      string `json:"macAddress"`
	MachineID       int64  `json:"machineId"`
	WorkShop        string `json:"workShop"`
	MachineName     string `json:"machineName"`
	MachineNumber   string `json:"machineNumber"`
}

// WorkShop WorkShop
type WorkShop struct {
	ID              int64  `json:"id"`
	Number          string `json:"number"`
	Name            string `json:"name"`
	Type            int64  `json:"type"`
	Remark          string `json:"remark"`
	Count           int64  `json:"count"`
	FirstStatusTime int64  `json:"firstStatusTime"`
}
