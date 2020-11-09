package datacollect

import (
	dcstatusdao "emuMolding/daos/dcStatusDao"
	"emuMolding/daos/didao"
	"emuMolding/models"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"runtime"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

func init() {
	dcTr = &http.Transport{
		DisableKeepAlives: true,
		MaxIdleConns:      -1,
	}
}

// FetchLoop FetchLoop
func FetchLoop(dcs []models.Wise) {
	random := time.Duration(rand.Intn(6) + 5)
	ticker := time.NewTicker(random * time.Second)
	for range ticker.C {
		for _, dc := range dcs {
			FetchDc(dc)
		}
	}
}

// FetchDc FetchDc
func FetchDc(wise models.Wise) (err error) {
	defer func() {
		if r := recover(); r != nil {
			if _, ok := r.(runtime.Error); ok {
				panic(r)
			}
			err = r.(error)
			beego.Error(err)
		}
	}()
	if lock, ok := fetchDcLockMap.LoadOrStore(wise.MacAddress, false); ok {
		if lock.(bool) {
			return
		}
	}
	fetchDcLockMap.Store(wise.MacAddress, true)
	defer fetchDcLockMap.Store(wise.MacAddress, false)
	dcSetting, err := getDcLog(wise)
	if err != nil {
		panic(err)
	} else if dcSetting.UID != 1 || dcSetting.MAC != 0 || dcSetting.TmF != 0 || dcSetting.Fltr != 1 {
		putDc(wise)
		return
	}
	fetchTime, err := getFetchTime(wise, dcSetting.TFst)
	if err != nil {
		panic(err)
	}
	if err := patchDc(wise, fetchTime); err != nil {
		panic(err)
	}
	if tst, tend, err := checkGetDcLog(wise); err != nil {
		panic(err)
	} else if tst != fetchTime.TSt || tend != fetchTime.TEnd {
		beego.Informational(tst, fetchTime.TSt, tend, fetchTime.TEnd)
		panic(errors.New("Fetch time is not match"))
	}
	dis, err := getDcData(wise)
	if err != nil {
		panic(err)
	}
	wiseFetchTimeMap.Store(wise.MacAddress, fetchTime)
	if _, err := didao.CreateMultiDi(dis); err != nil {
		panic(err)
	} else {
		if err := convertDiToStatus(dis, wise); err != nil {
			panic(err)
		}
	}
	return err
}

func getFetchTime(wise models.Wise, wiseMax int64) (timeRange models.DcPatchBody, err error) {
	if last, ok := wiseFetchTimeMap.Load(wise.MacAddress); ok {
		timeRange.TSt = last.(models.DcPatchBody).TEnd + 1
	} else {
		if lastdi, err := didao.ReadOneDi(wise.MacAddress); err != nil {
			if err == orm.ErrNoRows {
				timeRange.TSt = 0
			} else {
				return timeRange, err
			}
		} else {
			timeRange.TSt = lastdi.Timestamp + 1
		}
	}
	if timeRange.TSt == 0 {
		timeRange.TSt = time.Now().Unix() - wise.PutTimeInterval - 2
	}
	timeRange.TEnd = wiseMax - 2
	if timeRange.TSt > timeRange.TEnd {
		return
	}
	return timeRange, err
}

func convertDiToStatus(dis []models.Di, wise models.Wise) (err error) {
	diChan := make(chan models.Di, len(dis))
	statusChan := make(chan models.DcStatus)
	go func() {
		for _, di := range dis {
			diChan <- di
		}
		close(diChan)
	}()
	go func() (err error) {
		var statusArr []models.DcStatus
		for status := range statusChan {
			statusArr = append(statusArr, status)
			if len(statusArr)%100 == 0 {
				if _, err := dcstatusdao.CreateMultiStatus(statusArr); err != nil {
					return err
				}
				statusArr = nil
			}
			if _, err := dcstatusdao.CreateMultiStatus(statusArr); err != nil {
				return err
			}
			statusArr = nil
		}
		return err
	}()
	for v := range diChan {
		temp := cycleStep{}
		tmpStatus := models.DcStatus{
			MacAddress: v.MacAddress,
			Timestamp:  v.Timestamp,
		}
		if lastCycleStatus, ok := cycleMap.Load(wise.MacAddress); ok {
			temp = lastCycleStatus.(cycleStep)
		}
		if v.Di0 == 0 {
			temp = cycleStep{}
			tmpStatus.Status = 2
			statusChan <- tmpStatus
		} else if v.Di3 == 1 {
			temp = cycleStep{}
			tmpStatus.Status = 4
			statusChan <- tmpStatus
		} else {
			if temp.step1.timeStamp == 0 && temp.step2.timeStamp == 0 && temp.step3.timeStamp == 0 {
				if v.Di2 == 1 {
					temp.step1.timeStamp = v.Timestamp
					temp.step1.sysTk = v.SysTk
				}
			} else if temp.step1.timeStamp != 0 && temp.step2.timeStamp == 0 && temp.step3.timeStamp == 0 {
				if v.Di1 == 1 {
					if float64(v.SysTk-temp.step1.sysTk)/1000 > float64(wise.IdleTime) {
						temp = cycleStep{}
						tmpStatus.Status = 3
						statusChan <- tmpStatus
					} else {
						temp.step2.timeStamp = v.Timestamp
						temp.step2.sysTk = v.SysTk
					}
				}
			} else if temp.step1.timeStamp != 0 && temp.step2.timeStamp != 0 && temp.step3.timeStamp == 0 {
				if v.Di2 == 1 {
					if float64(v.SysTk-temp.step1.sysTk)/1000 > float64(wise.IdleTime) {
						tmpStatus.Status = 3
						statusChan <- tmpStatus
					} else {
						cycleTime := float64((v.SysTk - temp.step1.sysTk)) / 1000
						// ct2 := float64((v.Timestamp*1000 + v.SysTk%1000 - temp.step1.timeStamp*1000 - temp.step1.sysTk%1000)) / 1000
						result := fmt.Sprintf("%s, CycleTime: %.3f", wise.MacAddress, cycleTime)
						beego.Informational(result)
						tmpStatus.CycleTime = cycleTime
						tmpStatus.Status = 5
						statusChan <- tmpStatus
					}
					temp.step1.timeStamp = v.Timestamp
					temp.step1.sysTk = v.SysTk
					temp.step2.timeStamp = 0
					temp.step2.sysTk = 0
				}
			}
		}
		cycleMap.Store(wise.MacAddress, temp)
	}
	close(statusChan)
	return err
}
