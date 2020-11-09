package models

// RealStatus RealStatus
type RealStatus struct {
	ID            int64   `json:"id"`
	MachineName   string  `json:"machineName"`
	MachineNumber string  `json:"machineNumber"`
	MachineModel  string  `json:"machineModel"`
	Brand         string  `json:"brand"`
	DcID          int64   `json:"dcId"`
	WorkShopID    int64   `json:"workShopId"`
	Status        int64   `json:"status"`
	CycleTime     float64 `json:"cycleTime"`
	Sct           float64 `json:"sct"`
	RedTimes      int64   `json:"redTimes"`
	YellowTimes   int64   `json:"yellowTimes"`
	GreenTimes    int64   `json:"greenTimes"`
	ContinuedTime int64   `json:"continuedTime"`
}

// DataRealStatus DataRealStatus
type DataRealStatus struct {
	MachineNumber string `json:"machineNumber"`
	Scheduled     bool   `json:"scheduled"`
	Tasked        bool   `json:"tasked"`
	Status        int64  `json:"status"`
	PlanCycleTime int64  `json:"planCycleTime"`
}

// StatusReturn StatusReturn
type StatusReturn struct {
	Data     []DataRealStatus `json:"data"`
	Response string           `json:"response"`
}
