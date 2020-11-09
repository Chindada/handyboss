package models

// NewSchedule NewSchedule
type NewSchedule struct {
	ScheduleSerial  string  `json:"scheduleSerial"`
	Status          int64   `json:"status"`
	StartTime       int64   `json:"startTime"`
	EndTime         int64   `json:"endTime"`
	MachineID       int64   `json:"machineId"`
	MoldID          int64   `json:"moldId"`
	MoldCycleTime   int64   `json:"moldCycleTime"`
	MoldGreenRange  float64 `json:"moldGreenRange"`
	MoldYellowRange float64 `json:"moldYellowRange"`
	Qty             int64   `json:"qty"`
}
