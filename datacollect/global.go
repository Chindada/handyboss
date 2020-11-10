package datacollect

import (
	"emuMolding/models"
	"net/http"
	"sync"
)

func init() {
	dcTr = &http.Transport{
		DisableKeepAlives: true,
		MaxIdleConns:      -1,
	}
	dcA := models.Wise{
		Token:           "Basic cm9vdDptaXRyb290",
		MacAddress:      "00D0C9E34A12",
		IdleTime:        300,
		PutTimeInterval: 1800,
		IP:              "192.168.10.119",
	}
	dcB := models.Wise{
		Token:           "Basic cm9vdDptaXRyb290",
		MacAddress:      "00D0C9E50453",
		IdleTime:        300,
		PutTimeInterval: 1800,
		IP:              "192.168.10.192",
	}
	dcC := models.Wise{
		Token:           "Basic cm9vdDptaXRyb290",
		MacAddress:      "00D0C9E349F4",
		IdleTime:        300,
		PutTimeInterval: 1800,
		IP:              "192.168.10.145",
	}
	Dcs = append(Dcs, dcA, dcB, dcC)
}

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
