package fakedata

import (
	"emuMolding/libs/restapitools"
	"emuMolding/models"
	"emuMolding/sysinit"
	"encoding/json"
	"errors"
	"io/ioutil"
	"math/rand"
	"runtime"
	"strconv"
	"time"

	"github.com/astaxie/beego"
)

// Loop Loop
func Loop() {
	ticker := time.NewTicker(5 * time.Second)
	for range ticker.C {
		CalibrateMachine()
		getRealStatus()
		if len(Machine) == 0 {
			beego.Informational("Machine is empty")
			continue
		}
		for _, k := range Machine {
			// checkIdle(k, false)
			if savedclock, ok := savedcFetchLockMap.LoadOrStore(k.MacAddress, false); ok {
				if !savedclock.(bool) {
					go func(k models.Machine) {
						saveDc(k)
					}(k)
				}
			}
		}
	}
}

// CalibrateMachine CalibrateMachine
func CalibrateMachine() {
	GetIP()
	LoginSystem()
	GetSystemMachine()
	GetWorkShop()
	for _, n := range workShopNumber.GetWorkShops() {
		GetMachine(n)
	}
}

// GetIP GetIP
func GetIP() (err error) {
	defer func() {
		if r := recover(); r != nil {
			if _, ok := r.(runtime.Error); ok {
				panic(r)
			}
			err = r.(error)
			beego.Error(err)
		}
	}()
	settingdb, err := sysinit.CreateSettingsConnection()
	if err != nil {
		panic(err)
	}
	defer settingdb.Close()
	var dbip models.Settings
	settingdb.Where("`key` = 'ip'").Last(&dbip)
	if dbip.Value == "" {
		dbip = models.Settings{
			Key:   "ip",
			Value: beego.AppConfig.String("fakedata::dataip") + ":8885",
		}
		settingdb.Create(&dbip)
	}
	ip = dbip.Value
	var systemip models.Settings
	settingdb.Where("`key` = 'system_ip'").Last(&systemip)
	if systemip.Value == "" {
		systemip = models.Settings{
			Key:   "system_ip",
			Value: beego.AppConfig.String("fakedata::systemip") + ":8887",
		}
		settingdb.Create(&systemip)
	}
	systemIP = systemip.Value
	return err
}

// GetWorkShop GetWorkShop
func GetWorkShop() (err error) {
	defer func() {
		if r := recover(); r != nil {
			if _, ok := r.(runtime.Error); ok {
				panic(r)
			}
			err = r.(error)
			beego.Error(err)
		}
	}()
	settingdb, err := sysinit.CreateSettingsConnection()
	if err != nil {
		panic(err)
	}
	defer settingdb.Close()

	var dbWorkShop models.Settings
	settingdb.Where("`key` = 'work_shop'").Last(&dbWorkShop)
	if dbWorkShop.Value == "" {
		dbWorkShop = models.Settings{
			Key:   "work_shop",
			Value: "100",
		}
		settingdb.Create(&dbWorkShop)
	}
	workShopNumber = dbWorkShop
	return err
}

// GetMachine GetWorkshopNumber
func GetMachine(workShopNumber string) (err error) {
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
	api.IP = ip
	api.URL = dataGetMachine
	headers := make(map[string]string)
	headers["workShopNumber"] = workShopNumber
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
	var data models.MachineData
	if err := json.Unmarshal(body, &data); err != nil {
		panic(err)
	}
	// beego.Informational("System machine lenth", len(systemMachineDetail))
	var tempMachine []models.Machine
	for _, v := range data.Info {
		tempMachine = append(tempMachine, v)
		Machine = tempMachine
		var temp models.LocalMachineList
		var isStable bool
		sqlite3db.Where("`mac_address` = ?", v.MacAddress).Last(&temp)
		if temp.MachineNumber == "" {
			probabilityA := rand.Intn(100) + 1
			if probabilityA <= 40 {
				isStable = false
			} else {
				isStable = true
			}
			var machineID int64
			for _, k := range systemMachineDetail {
				if k.MachineNumber == v.MachineNumber {
					machineID = k.ID
				}
			}
			temp = models.LocalMachineList{
				MachineID:       machineID,
				MachineNumber:   v.MachineNumber,
				MacAddress:      v.MacAddress,
				PutTimeInterval: v.PutTimeInterval,
				IdleTime:        v.IdleTime,
				DcAuthorization: v.DcAuthorization,
				ActionTime:      0,
				CycleTime:       int64(rand.Intn(20) + 20),
				Stable:          isStable,
				Idle:            true,
			}
			sqlite3db.Create(&temp)
		}
	}
	return err
}

func generateFakeActionTime(macAddress string) (err error) {
	defer func() {
		if r := recover(); r != nil {
			if _, ok := r.(runtime.Error); ok {
				panic(r)
			}
			err = r.(error)
			beego.Error(err)
		}
	}()
	var temp models.LocalMachineList
	sqlite3db.Where("`mac_address` = ?", macAddress).Last(&temp)
	if temp.Status == 1 || temp.Status == 2 {
		if temp.ActionTime != 0 {
			if temp.CycleTime == 0 {
				temp.CycleTime = 60
			}
			temp.ActionTime += temp.CycleTime / 2
			if time.Now().Unix()-temp.ActionTime > temp.IdleTime {
				temp.ActionTime = time.Now().Unix() - 10
			}
			sqlite3db.Model(&temp).Update("action_time", temp.ActionTime)
		} else {
			sqlite3db.Model(&temp).Update("action_time", time.Now().Unix()-10)
		}
	} else {
		if int64(time.Now().Unix())-temp.ActionTime > 0 {
			sqlite3db.Model(&temp).Update("status", 1)
		}
	}
	return err
}

func saveDc(param models.Machine) (err error) {
	defer func() {
		if r := recover(); r != nil {
			if _, ok := r.(runtime.Error); ok {
				panic(r)
			}
			err = r.(error)
			beego.Error(err)
		}
	}()
	savedcFetchLockMap.Store(param.MacAddress, true)
	defer savedcFetchLockMap.Store(param.MacAddress, false)
	if isInit, ok := firstInitMap.Load(param.MachineNumber); ok {
		if isInit.(bool) {
			return
		}
	}
	if start, end, err := getFetchTime(param.MachineNumber); err != nil {
		panic(err)
	} else {
		// beego.Informational(start, end)
		fake := generateArr(start, end, param.MacAddress)
		var api restapitools.PostArg
		api.IP = ip
		api.URL = dataDcData
		api.Body = fake
		headers := make(map[string]string)
		headers["machineNumber"] = param.MachineNumber
		headers["lastTime"] = strconv.FormatInt(end, 10)
		headers["idleTime"] = strconv.FormatInt(param.IdleTime, 10)
		api.Headers = headers
		resp, err := api.Post()
		if err != nil {
			panic(err)
		} else if resp != nil {
			defer resp.Body.Close()
		}
	}
	return err
}

func generateArr(startTime, endTime int64, macAddress string) (fake []models.Di) {
	var machineFakeData models.LocalMachineList
	var cycleTime int64
	sqlite3db.Where("`mac_address` = ?", macAddress).Last(&machineFakeData)
	if machineFakeData.ActionTime > endTime {
		return nil
	}
	if len(realSchedules) == 0 {
		return nil
	}
	for _, v := range realSchedules {
		if v.MachineID == machineFakeData.MachineID {
			if startTime*1000 < v.StartTime && endTime*1000 > v.StartTime {
				cycleTime = v.MoldCycleTime
				startTime = v.StartTime / 1000
				break
			} else if startTime*1000 > v.StartTime && endTime*1000 < v.EndTime {
				cycleTime = v.MoldCycleTime
				break
			} else if startTime*1000 < v.EndTime && endTime*1000 > v.EndTime {
				cycleTime = v.MoldCycleTime
				endTime = v.EndTime / 1000
				break
			}
		}
	}
	if cycleTime == 0 {
		return nil
	}
	var shots, start int64
	action := machineFakeData.Status
	if machineFakeData.ActionTime == 0 {
		machineFakeData.ActionTime = startTime
	}
	if endTime-startTime > cycleTime {
		shots = (endTime - machineFakeData.ActionTime) / cycleTime
		start = machineFakeData.ActionTime * 1000
	}
	beego.Informational(machineFakeData.MachineNumber, start, shots)
	var i int64
	for i = 0; i <= shots; i++ {
		diTime := start
		probability := rand.Intn(100) + 1
		if probability <= 5 {
			diTime -= (cycleTime*1000/2/5 + int64(rand.Intn(9))*100)
		} else if probability <= 10 {
			diTime += (cycleTime*1000/2/5 + int64(rand.Intn(9))*100)
		} else {
			diTime += int64(rand.Intn(9)) * 100
		}
		if action == 2 {
			temp := models.Di{
				Timestamp: diTime,
				Di0:       1,
				Di1:       0,
				Di2:       1,
			}
			fake = append(fake, temp)
			action = 1
		} else if action == 1 {
			temp := models.Di{
				Timestamp: diTime,
				Di0:       1,
				Di1:       1,
				Di2:       0,
			}
			fake = append(fake, temp)
			action = 2
		}
		start += (cycleTime * 1000) / 2
	}
	sqlite3db.Model(&machineFakeData).Update("action_time", start/1000)
	sqlite3db.Model(&machineFakeData).Update("status", action)
	sqlite3db.Model(&machineFakeData).Update("cycle_time", cycleTime)
	fakeDataMode := beego.AppConfig.String("fakedata::fakeDataMode")
	if fakeDataMode == "init" {

		temp1 := models.Di{
			Timestamp: (firstDayTimeStamp - 20) * 1000,
			Di0:       1,
			Di1:       0,
			Di2:       1,
		}
		temp2 := models.Di{
			Timestamp: (firstDayTimeStamp - 10) * 1000,
			Di0:       1,
			Di1:       1,
			Di2:       0,
		}
		temp3 := models.Di{
			Timestamp: firstDayTimeStamp * 1000,
			Di0:       1,
			Di1:       0,
			Di2:       1,
		}
		fake = nil
		fake = append(fake, temp1)
		fake = append(fake, temp2)
		fake = append(fake, temp3)
		firstInitMap.Store(machineFakeData.MachineNumber, true)
	}
	return fake
}

func getFetchTime(machineNumber string) (startTime, endTime int64, err error) {
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
	api.IP = ip
	api.URL = dataFetchInterval
	headers := make(map[string]string)
	headers["machineNumber"] = machineNumber
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
	var data models.FetchTime
	if err := json.Unmarshal(body, &data); err != nil {
		panic(err)
	}
	startTime = data.Min
	endTime = data.Max
	if endTime-time.Now().Unix() >= 0 {
		endTime = time.Now().Unix() - 2
	}
	return startTime, endTime, err
}

func checkIdle(param models.Machine, init bool) (err error) {
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
	api.IP = ip
	api.URL = dataCheckIdle
	headers := make(map[string]string)
	headers["machineNumber"] = param.MachineNumber
	if init {
		headers["idleTime"] = "3"
	} else {
		headers["idleTime"] = strconv.FormatInt(param.IdleTime, 10)
	}
	api.Headers = headers
	resp, err := api.Get()
	if err != nil {
		panic(err)
	} else if resp != nil {
		defer resp.Body.Close()
	}
	return err
}

func getRealStatus() (err error) {
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
	api.IP = ip
	api.URL = dataMachineStatus
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
	var data []models.DataRealStatus
	if err := json.Unmarshal(body, &data); err != nil {
		panic(err)
	}
	for _, v := range data {
		realStatusMap[v.MachineNumber] = v.Status
		scheduledMap[v.MachineNumber] = v.Scheduled
		PlanCycleTimeMap[v.MachineNumber] = v.PlanCycleTime
	}
	return err
}

// GetSettings GetSettings
func GetSettings() (data []models.Settings, err error) {
	defer func() {
		if r := recover(); r != nil {
			if _, ok := r.(runtime.Error); ok {
				panic(r)
			}
			err = r.(error)
			beego.Error(err)
		}
	}()
	settingdb, err := sysinit.CreateSettingsConnection()
	if err != nil {
		panic(err)
	}
	defer settingdb.Close()
	settingdb.Find(&data)
	return data, err
}

// GetMachineDetail GetMachineDetail
func GetMachineDetail(machineNumber string) (data models.LocalMachineList, err error) {
	result := sqlite3db.Where("`machine_number` = ?", machineNumber).Last(&data)
	err = result.Error
	return data, err
}

// TaskTime TaskTime
func TaskTime(machineAction models.MachineAction) (err error) {
	defer func() {
		if r := recover(); r != nil {
			if _, ok := r.(runtime.Error); ok {
				panic(r)
			}
			err = r.(error)
			beego.Error(err)
		}
	}()
	var temp models.LocalMachineList
	sqlite3db.Where("`machine_number` = ?", machineAction.MachineNumber).Last(&temp)
	result := sqlite3db.Model(&temp).Updates(models.LocalMachineList{
		Status:     int(machineAction.Status),
		ActionTime: int64(time.Now().Unix()) + machineAction.Interval,
		CycleTime:  machineAction.CycleTime,
	})
	if result.Error != nil {
		err = result.Error
		panic(err)
	}
	var tempMachine models.Machine
	tempMachine.MachineNumber = machineAction.MachineNumber
	tempMachine.IdleTime = machineAction.IdleTime
	if machineAction.Status != 3 {
		checkIdle(tempMachine, true)
	}
	return err
}

// UpdateWorkShop UpdateWorkShop
func UpdateWorkShop(workshop string) (err error) {
	defer func() {
		if r := recover(); r != nil {
			if _, ok := r.(runtime.Error); ok {
				panic(r)
			}
			err = r.(error)
			beego.Error(err)
		}
	}()
	if workshop == "" {
		err = errors.New("WorkShop is empty")
		panic(err)
	}
	settingdb, err := sysinit.CreateSettingsConnection()
	if err != nil {
		panic(err)
	}
	defer settingdb.Close()
	var dbWorkShop models.Settings
	settingdb.Where("`key` = 'work_shop'").Last(&dbWorkShop)
	result := settingdb.Model(&dbWorkShop).Update("value", workshop)
	if result.Error != nil {
		err = result.Error
		panic(err)
	}
	return err
}

// UpdateIP UpdateIP
func UpdateIP(ip string) (err error) {
	defer func() {
		if r := recover(); r != nil {
			if _, ok := r.(runtime.Error); ok {
				panic(r)
			}
			err = r.(error)
			beego.Error(err)
		}
	}()
	if ip == "" {
		err = errors.New("IP is empty")
		panic(err)
	}
	settingdb, err := sysinit.CreateSettingsConnection()
	if err != nil {
		panic(err)
	}
	defer settingdb.Close()
	var dbip models.Settings
	settingdb.Where("`key` = 'ip'").Last(&dbip)
	result := settingdb.Model(&dbip).Update("value", ip)
	if result.Error != nil {
		err = result.Error
		panic(err)
	}
	return err
}
