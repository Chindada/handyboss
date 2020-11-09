package datacollect

import (
	"emuMolding/models"
	"net/http"
	"sync"
)

var dcTr *http.Transport
var dcIP string

var wiseFetchTimeMap sync.Map
var fetchDcLockMap sync.Map
var cycleMap sync.Map

// Dcs Dcs
var Dcs []models.Wise

type cycleStep struct {
	step1 timeWithSysTk
	step2 timeWithSysTk
	step3 timeWithSysTk
}

type timeWithSysTk struct {
	timeStamp int64
	sysTk     int64
}
