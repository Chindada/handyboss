package models

// Mold Mold
type Mold struct {
	ID           int64   `json:"id"`
	CavityNumber int64   `json:"cavityNumber"`
	CycleTime    int64   `json:"cycleTime"`
	UpTime       int64   `json:"upTime"`
	DownTime     int64   `json:"downTime"`
	Name         string  `json:"name"`
	ProductModel string  `json:"productModel"`
	Number       string  `json:"number"`
	GreenRange   float64 `json:"greenRange"`
	YellowRange  float64 `json:"yellowRange"`
}

// MoldMachineRel MoldMachineRel
type MoldMachineRel struct {
	MachineID int64 `json:"machineId"`
	MoldID    int64 `json:"moldId"`
}

// SystemMachine SystemMachine
type SystemMachine struct {
	ID                         int64   `json:"id"`
	Brand                      string  `json:"brand"`
	MachineName                string  `json:"machineName"`
	MachineNumber              string  `json:"machineNumber"`
	Model                      string  `json:"model"`
	PurchasDate                int64   `json:"purchasDate"`
	ManufactureDate            int64   `json:"manufactureDate"`
	ScrewDiameter              int64   `json:"screwDiameter"`
	ScrewRatio                 int64   `json:"screwRatio"`
	TheoreticalInjectionVolume int64   `json:"theoreticalInjectionVolume"`
	InjectionVolume            int64   `json:"injectionVolume"`
	MaxInjectionPressure       int64   `json:"maxInjectionPressure"`
	InjectionRate              int64   `json:"injectionRate"`
	InjectionSpeed             int64   `json:"injectionSpeed"`
	ShotStroke                 int64   `json:"shotStroke"`
	ShootingStroke             int64   `json:"shootingStroke"`
	NozzleClosedPower          int64   `json:"nozzleClosedPower"`
	HeatingSegmentNumber       int64   `json:"heatingSegmentNumber"`
	TubeTotalHeat              int64   `json:"tubeTotalHeat"`
	ClampingPower              int64   `json:"clampingPower"`
	MaxOpenStroke              int64   `json:"maxOpenStroke"`
	MinMoldThickness           int64   `json:"minMoldThickness"`
	MaxMoldThickness           int64   `json:"maxMoldThickness"`
	MoldWidth                  int64   `json:"moldWidth"`
	MoldHeight                 int64   `json:"moldHeight"`
	MaxPlasticAmount           int64   `json:"maxPlasticAmount"`
	LastUpdateTime             int64   `json:"lastUpdateTime"`
	DcID                       int64   `json:"dcId"`
	WorkShopID                 int64   `json:"workShopId"`
	Analyze                    bool    `json:"analyze"`
	WorkShopName               string  `json:"workShopName"`
	WorkShopNumber             string  `json:"workShopNumber"`
	Status                     int64   `json:"status"`
	MaterialType               int64   `json:"materialType"`
	Remark                     string  `json:"remark"`
	Width                      float64 `json:"width"`
	Height                     float64 `json:"height"`
}
