package fakedata

import (
	"emuMolding/models"
	"emuMolding/sysinit"
	"net/http"
	"sync"

	"github.com/astaxie/beego"
	"github.com/jinzhu/gorm"
)

func init() {
	var err error
	realStatusMap = make(map[string]int64)
	scheduledMap = make(map[string]bool)
	PlanCycleTimeMap = make(map[string]int64)
	firstDayTimeStamp, err = beego.AppConfig.Int64("fakedata::firstDayTimeStamp")
	if err != nil {
		panic(err)
	}
	sqlite3db, err = sysinit.CreateMachineListConnection()
	if err != nil {
		panic(err)
	}
	CalibrateMachine()
	_, err = GetSchedule()
	if err != nil {
		panic(err)
	}
}

// PlanCycleTimeMap PlanCycleTimeMap
var PlanCycleTimeMap map[string]int64

// Machine All machine
var Machine []models.Machine

// System Token
var jwt *http.Cookie

// Savedc flow control lock
var savedcFetchLockMap sync.Map

// AddScheduleLock AddScheduleLock
var AddScheduleLock sync.Map

// Store status from data
var realStatusMap map[string]int64

// Store whether machine is scheduled from data
var scheduledMap map[string]bool

// Fake data sqlite
var sqlite3db *gorm.DB

// Data ip
var ip string

// System ip
var systemIP string

// Workshop number in sqlite
var workShopNumber models.Settings

// All machine ID
var systemMachineID []int64
var systemMachineDetail []models.SystemMachine

// All mold ID
var systemMoldID []int64

// All Mold
var allMolds []models.Mold

// For fakeschedule machine's latest time can start
var machineStartTimeMap sync.Map

// Mold last used time
var tempTimeMap sync.Map

// Store temp schedule sent to system
var fakeScheduleMap sync.Map

// Load from fakeScheduleMap
var schedules []models.NewSchedule
var realSchedules []models.NewSchedule
var fakeSchedule []models.NewSchedule

var moldLastUsedTime sync.Map
var machineLastUsedTime sync.Map
var key []models.NewSchedule

var firstDayTimeStamp int64
var firstInitMap sync.Map

// Schedule mold's detail
type moldTime struct {
	moldID       int64
	totalTime    float64
	cycleTime    int64
	cavityNumber int64
	startTime    int64
	endTime      int64
}

// Mold last used time and detail
type tempTime struct {
	machineID int64
	startTime int64
	endTime   int64
}

// URL
const (
	systemLoginURL       = "/ioms5system/auth/login"
	systemNewSchedule    = "/ioms5system/function/schedule/oee/multi"
	systemScheduleURL    = "/ioms5system/function/schedule/oee"
	systemMoldMachineRel = "/ioms5system/function/mold/importrel"
	systemMachine        = "/ioms5system/crud/basic/machine"
	systemMold           = "/ioms5system/crud/basic/mold"
	systemNewDc          = "/ioms5system/crud/system/dc"
	systemNewWorkShop    = "/ioms5system/crud/basic/workshop"

	dataGetMachine    = "/ioms5data/datacollect/machine"
	dataDcData        = "/ioms5data/datacollect/dcData"
	dataFetchInterval = "/ioms5data/datacollect/fetchInterval"
	dataCheckIdle     = "/ioms5data/datacollect/checkidle"
	dataMachineStatus = "/ioms5data/analyze/kanban/task"
)
