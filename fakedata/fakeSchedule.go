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
	defer func() {
		if r := recover(); r != nil {
			if _, ok := r.(runtime.Error); ok {
				panic(r)
			}
			err = r.(error)
			beego.Error(err)
		}
	}()
	var api restapitools.GetArg
	api.IP = systemIP
	api.URL = systemScheduleURL
	api.Token = jwt
	resp, err := api.Get()
	if err != nil {
		panic(err)
	} else if resp != nil {
		defer resp.Body.Close()
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	if err := json.Unmarshal(body, &schedule); err != nil {
		panic(err)
	}
	realSchedules = schedule
	return schedule, err
}

// GetScheduleFilter GetScheduleFilter
func GetScheduleFilter(timeRange string) (schedule []models.NewSchedule, err error) {
	defer func() {
		if r := recover(); r != nil {
			if _, ok := r.(runtime.Error); ok {
				panic(r)
			}
			err = r.(error)
			beego.Error(err)
		}
	}()
	var api restapitools.GetArg
	api.IP = systemIP
	api.URL = systemScheduleURL
	api.Token = jwt
	headers := make(map[string]string)
	headers["timeRange"] = timeRange
	api.Headers = headers
	resp, err := api.Get()
	if err != nil {
		panic(err)
	} else if resp != nil {
		defer resp.Body.Close()
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	if err := json.Unmarshal(body, &schedule); err != nil {
		panic(err)
	}
	return schedule, err
}

// CreateMultiScheduleFromSlice CreateMultiScheduleFromSlice
func CreateMultiScheduleFromSlice() {
	defer func() {
		AddScheduleLock.Store("active", false)
		beego.Informational("End fake sc.")
	}()
	var fakeSchedules []models.NewSchedule
	beego.Informational("There are", len(key), "sc")
	total := len(key)
	for _, v := range key {
		// if sc, ok := tempSchedule.LoadAndDelete(v); ok {
		// 	fakeSchedules = append(fakeSchedules, sc.(models.NewSchedule))
		// }
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
	GetSystemMold()
	GetSystemMachine()
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
	GetSchedule()
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
			if machineLast, ok := machineLastUsedTime.LoadOrStore(v.MachineID, v.EndTime); ok {
				if v.EndTime > machineLast.(int64) {
					machineLastUsedTime.Store(v.MachineID, v.EndTime)
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
	wg := sync.WaitGroup{}
	wg.Add(len(scheduleMolds))
	for _, v := range scheduleMolds {
		go func(moldWithTimeChan chan moldTime, v models.Mold, wg *sync.WaitGroup) {
			defer wg.Done()
			temp := moldTime{
				moldID:       v.ID,
				totalTime:    rand.Float64() + float64(rand.Intn(5)+6),
				cycleTime:    v.CycleTime,
				cavityNumber: v.CavityNumber,
			}
			moldWithTimeChan <- temp
		}(moldWithTimeChan, v, &wg)
	}
	go func() {
		wg.Wait()
		close(moldWithTimeChan)
	}()

	var scTime sync.Map
	for mold := range moldWithTimeChan {
		moldEndTimeChan := make(chan tempTime, len(systemMachineID))
		wg2 := sync.WaitGroup{}
		wg2.Add(len(systemMachineID))
		for _, v := range systemMachineID {
			go func(moldEndTimeChan chan tempTime, v int64, wg2 *sync.WaitGroup) {
				defer wg2.Done()
				generateEndTime(v, mold, moldEndTimeChan)
			}(moldEndTimeChan, v, &wg2)
		}
		go func() {
			wg2.Wait()
			close(moldEndTimeChan)
		}()
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
			return
		}
		if end > now+86400*3*1000 {
			return
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
			if probabilityB <= 30 {
				key = append(key, schedule)
			}
		}
	}
	return err
}

func generateEndTime(machineID int64, x moldTime, moldEndTimeChan chan tempTime) (err error) {
	defer func() {
		if r := recover(); r != nil {
			if _, ok := r.(runtime.Error); ok {
				panic(r)
			}
			err = r.(error)
			beego.Error(err)
		}
	}()
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
	return err
}
