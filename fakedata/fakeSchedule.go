package fakedata

import (
	"emuMolding/libs/restapitools"
	"emuMolding/models"
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"runtime"
	"strconv"
	"sync"
	"time"

	"github.com/astaxie/beego"
)

// GetSchedule GetSchedule
func GetSchedule() (schedule []models.NewSchedule, err error) {
	var api restapitools.GetArg
	api.IP = systemIP
	api.URL = systemScheduleURL
	api.Token = jwt
	resp, err := api.Get()
	if err != nil {
		return schedule, err
	} else if resp != nil {
		defer resp.Body.Close()
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return schedule, err
	}
	if err := json.Unmarshal(body, &schedule); err != nil {
		return schedule, err
	}
	var temp []models.NewSchedule
	for _, machineID := range systemMachineID {
		for _, s := range schedule {
			if s.MachineID == machineID {
				temp = append(temp, s)
			}
		}
		machineRealSchedules.Store(machineID, temp)
		temp = nil
	}
	realSchedules = schedule
	return schedule, err
}

// GetScheduleFilter GetScheduleFilter
func GetScheduleFilter(timeRange string) (schedule []models.NewSchedule, err error) {
	var api restapitools.GetArg
	api.IP = systemIP
	api.URL = systemScheduleURL
	api.Token = jwt
	headers := make(map[string]string)
	headers["timeRange"] = timeRange
	api.Headers = headers
	resp, err := api.Get()
	if err != nil {
		return schedule, err
	} else if resp != nil {
		defer resp.Body.Close()
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return schedule, err
	}
	if err := json.Unmarshal(body, &schedule); err != nil {
		return schedule, err
	}
	return schedule, err
}

// CreateMultiScheduleFromSlice CreateMultiScheduleFromSlice
func CreateMultiScheduleFromSlice() (err error) {
	defer func() {
		if r := recover(); r != nil {
			if _, ok := r.(runtime.Error); ok {
				panic(r)
			}
			err = r.(error)
			beego.Error(err)
		}
	}()
	defer func() {
		AddScheduleLock.Store("active", false)
		beego.Informational("End fake sc.")
	}()
	var fakeSchedules []models.NewSchedule
	beego.Informational("There are", len(key), "sc")
	total := len(key)
	for _, v := range key {
		fakeSchedules = append(fakeSchedules, v)
		if len(fakeSchedules)%50 == 0 {
			total -= 50
			var api restapitools.PostArg
			api.IP = systemIP
			api.URL = systemNewSchedule
			api.Token = jwt
			api.Body = fakeSchedules
			resp, err := api.Post()
			if err != nil {
				beego.Error(err)
			} else if resp != nil {
				defer resp.Body.Close()
			}
			beego.Informational("Remain", total)
			fakeSchedules = nil
		}
	}
	if len(fakeSchedules) != 0 {
		var api restapitools.PostArg
		api.IP = systemIP
		api.URL = systemNewSchedule
		api.Token = jwt
		api.Body = fakeSchedules
		resp, err := api.Post()
		if err != nil {
			beego.Error(err)
		} else if resp != nil {
			defer resp.Body.Close()
		}
		beego.Informational("Final Insert", len(fakeSchedules))
		fakeSchedules = nil
	}
	key = nil
	return err
}

// FakeNewSchedule FakeNewSchedule
func FakeNewSchedule() (err error) {
	defer func() {
		if r := recover(); r != nil {
			if _, ok := r.(runtime.Error); ok {
				panic(r)
			}
			err = r.(error)
			beego.Error(err)
		}
	}()
	if err := GetSystemMold(); err != nil {
		panic(err)
	}
	if err := GetSystemMachine(); err != nil {
		panic(err)
	}
	if len(allMolds) == 0 || len(systemMachineID) == 0 {
		return
	}
	AddScheduleLock.Store("active", true)
	beego.Informational("Start fake sc.")
	num, err := beego.AppConfig.Int("fakedata::autoCreateScheduleNum")
	if err != nil {
		panic(err)
	}
	var scheduleMolds []models.Mold
	for i := 0; i < num; i++ {
		randomIndex := rand.Intn(len(allMolds))
		randomMold := allMolds[randomIndex]
		scheduleMolds = append(scheduleMolds, randomMold)
	}
	if _, err := GetSchedule(); err != nil {
		panic(err)
	}
	for _, k := range scheduleMolds {
		for _, v := range realSchedules {
			temp := tempTime{
				startTime: v.StartTime,
				endTime:   v.EndTime,
				machineID: v.MachineID,
			}
			if k.ID == v.MoldID {
				if moldLast, ok := moldLastUsedTime.LoadOrStore(v.MoldID, temp); ok {
					if v.EndTime > moldLast.(tempTime).endTime {
						moldLastUsedTime.Store(v.MoldID, temp)
					}
				}
			}
		}
	}
	for _, machineID := range systemMachineID {
		for _, v := range realSchedules {
			if machineID == v.MachineID {
				if machineLast, ok := machineLastUsedTime.LoadOrStore(v.MachineID, v.EndTime); ok {
					if v.EndTime > machineLast.(int64) {
						machineLastUsedTime.Store(v.MachineID, v.EndTime)
					}
				}
			}
		}
	}
	firstDayTimeStamp, err := beego.AppConfig.Int64("fakedata::firstDayTimeStamp")
	if err != nil {
		panic(err)
	}
	now := time.Now().Unix() * 1000
	for _, m := range scheduleMolds {
		if _, ok := moldLastUsedTime.Load(m.ID); !ok {
			temp := tempTime{
				startTime: firstDayTimeStamp*1000 - 86400*7*1000,
				endTime:   firstDayTimeStamp*1000 - 86400*6*1000,
				machineID: 0,
			}
			moldLastUsedTime.Store(m.ID, temp)
		}
	}
	for _, m := range systemMachineID {
		if _, ok := machineLastUsedTime.Load(m); !ok {
			machineLastUsedTime.Store(m, firstDayTimeStamp*1000-86400*6*1000)
		}
	}
	moldWithTimeChan := make(chan moldTime, len(scheduleMolds))
	go func() {
		for _, v := range scheduleMolds {
			temp := moldTime{
				moldID:       v.ID,
				totalTime:    rand.Float64() + float64(rand.Intn(5)+6),
				cycleTime:    v.CycleTime,
				cavityNumber: v.CavityNumber,
			}
			moldWithTimeChan <- temp
		}
		close(moldWithTimeChan)
	}()

	var scTime sync.Map
	for mold := range moldWithTimeChan {
		moldEndTimeChan := make(chan tempTime, len(systemMachineID))
		for _, v := range systemMachineID {
			generateEndTime(v, mold, moldEndTimeChan)
		}
		close(moldEndTimeChan)
		for k := range moldEndTimeChan {
			if fast, ok := scTime.LoadOrStore(mold, k); ok {
				if fast.(tempTime).endTime > k.endTime {
					scTime.Store(mold, k)
				}
			}
		}
		var temp tempTime
		var scStatus int64
		if fastest, ok := scTime.Load(mold); ok {
			temp = fastest.(tempTime)
		} else {
			continue
		}
		start := temp.startTime
		end := temp.endTime
		machineID := temp.machineID
		if now > end {
			scStatus = 1
		} else {
			scStatus = 2
		}
		if start == 0 {
			continue
		}
		if end > now+86400*3*1000 {
			continue
		}
		schedule := models.NewSchedule{
			ScheduleSerial:  "FC#" + strconv.Itoa(rand.Intn(10000)),
			Status:          scStatus,
			StartTime:       start,
			EndTime:         end,
			MachineID:       machineID,
			MoldID:          mold.moldID,
			MoldCycleTime:   mold.cycleTime,
			MoldGreenRange:  0.1,
			MoldYellowRange: 0.2,
			Qty:             int64(mold.totalTime*3600) / int64(mold.cycleTime),
		}
		moldLastUsedTime.Store(mold.moldID, temp)
		machineLastUsedTime.Store(machineID, end)

		startTm := time.Unix(start/1000, 0)
		startWeekDay := startTm.Weekday()
		if startWeekDay != time.Saturday && startWeekDay != time.Sunday {
			probabilityA := rand.Intn(100) + 1
			if probabilityA <= 80 {
				key = append(key, schedule)
			}
		} else {
			probabilityB := rand.Intn(100) + 1
			if probabilityB <= 40 {
				key = append(key, schedule)
			}
		}
	}
	return err
}

func generateEndTime(machineID int64, x moldTime, moldEndTimeChan chan tempTime) {
	var start, end int64
	if machineLast, ok := machineLastUsedTime.Load(machineID); ok {
		start = machineLast.(int64) + 3600*1000*(rand.Int63n(5)+1)
		end = start + int64(x.totalTime*3600*1000)
	}
	if moldLast, ok := moldLastUsedTime.Load(x.moldID); ok {
		if start < moldLast.(tempTime).startTime && end > moldLast.(tempTime).startTime {
			return
		} else if start > moldLast.(tempTime).startTime && end < moldLast.(tempTime).endTime {
			return
		} else if start < moldLast.(tempTime).endTime && end > moldLast.(tempTime).endTime {
			return
		} else if start < moldLast.(tempTime).startTime && end > moldLast.(tempTime).endTime {
			return
		}
	}
	temp := tempTime{
		machineID: machineID,
		startTime: start,
		endTime:   end,
	}
	moldEndTimeChan <- temp
}
