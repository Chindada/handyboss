package models

import "github.com/jinzhu/gorm"

// LocalMachineList LocalMachineList
type LocalMachineList struct {
	gorm.Model      `json:"-"`
	MachineNumber   string `gorm:"column:machine_number;" json:"machineNumber"`
	MachineID       int64  `gorm:"column:machine_id;" json:"machineID"`
	MacAddress      string `gorm:"column:mac_address;" json:"macAddress"`
	PutTimeInterval int64  `gorm:"column:put_time_interval;default:86400;" json:"putTimeInterval"`
	IdleTime        int64  `gorm:"column:idle_time;default:300;" json:"idleTime"`
	DcAuthorization string `gorm:"column:dc_authorization;" json:"dcAuthorization"`
	Status          int    `gorm:"column:status;default:1;" json:"status"`
	ActionTime      int64  `gorm:"column:action_time;" json:"actionTime"`
	CycleTime       int64  `gorm:"column:cycle_time;" json:"cycleTime"`
	Stable          bool   `gorm:"column:stable;" json:"stable"`
	Idle            bool   `gorm:"column:idle;default:1;" json:"idle"`
}

// Machine Machine
type Machine struct {
	MachineNumber   string `json:"machineNumber"`
	MacAddress      string `json:"macAddress"`
	PutTimeInterval int64  `json:"putTimeInterval"`
	IdleTime        int64  `json:"idleTime"`
	DcAuthorization string `json:"dcAuthorization"`
}

// MachineData MachineData
type MachineData struct {
	Info     []Machine `json:"info"`
	Response string    `json:"response"`
}

// FetchTime FetchTime
type FetchTime struct {
	Min      int64  `json:"min"`
	Max      int64  `json:"max"`
	Response string `json:"response"`
}

// MachineAction MachineAction for API
type MachineAction struct {
	MachineNumber string `json:"machineNumber"`
	MacAddress    string `json:"macAddress"`
	IdleTime      int64  `json:"idleTime"`
	Status        int64  `json:"status"`
	Interval      int64  `json:"interval"`
	CycleTime     int64  `json:"cycleTime"`
}
