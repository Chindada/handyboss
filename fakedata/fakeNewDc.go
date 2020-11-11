package fakedata

import (
	"encoding/json"
	"handyboss/libs/restapitools"
	"handyboss/models"
	"io/ioutil"
	"math/rand"
	"runtime"
	"strconv"
	"strings"

	"github.com/astaxie/beego"
	"github.com/tjarratt/babble"
)

// NewMachineLoop NewMachineLoop
func NewMachineLoop() (err error) {
	defer func() {
		if r := recover(); r != nil {
			if _, ok := r.(runtime.Error); ok {
				panic(r)
			}
			err = r.(error)
			beego.Error(err)
		}
	}()
	var insertnew bool
	if insertnew, err = beego.AppConfig.Bool("fakedata::addWsMachine"); err != nil {
		panic(err)
	} else if !insertnew {
		return
	}
	if err := NewWorkShop(); err != nil {
		panic(err)
	} else if ws, err := GetSystemWorkShop(); err != nil {
		panic(err)
	} else {
		for _, v := range ws {
			NewMachine(v)
		}
	}
	if err := BindMachine(); err != nil {
		panic(err)
	}
	return err
}

// NewDc NewDc
func NewDc(id int64) (err error) {
	babbler := babble.NewBabbler()
	babbler.Count = 1
	dcs := models.Dc{
		CreateTime:      0,
		DcAuthorization: "Basic cm9vdDptaXRyb290",
		ID:              0,
		IdleTime:        300,
		MacAddress:      strings.Title(strings.ToLower(babbler.Babble())),
		MachineID:       id,
		PutTimeInterval: 86400,
		WorkShop:        "100",
	}
	var api restapitools.PostArg
	api.Body = dcs
	api.IP = systemIP
	api.Token = jwt
	api.URL = systemNewDc
	headers := make(map[string]string)
	headers["token"] = "82589155"
	api.Headers = headers
	resp, err := api.Post()
	if err != nil {
		return err
	} else if resp != nil {
		defer resp.Body.Close()
	}
	return err
}

// BindMachine BindMachine
func BindMachine() (err error) {
	err = GetSystemMachine()
	for _, v := range systemMachineDetail {
		if v.DcID == 0 {
			if err := NewDc(v.ID); err != nil {
				return err
			}
		}
	}
	return err
}

// NewMachine NewMachine
func NewMachine(ws models.WorkShop) (err error) {
	babbler := babble.NewBabbler()
	babbler.Count = 1
	if err := GetSystemMachine(); err != nil {
		return err
	}
	var machineNum int
	availableName := make(map[string]bool)
	for _, v := range systemMachineDetail {
		if v.WorkShopID == ws.ID {
			machineNum++
			availableName[v.MachineNumber] = true
		}
	}
	var machineNumSetting int
	if machineNumSetting, err = beego.AppConfig.Int("fakedata::machineNum"); err != nil {
		return err
	} else if machineNum < machineNumSetting {
		for i := 0; i < machineNumSetting-machineNum; i++ {
			number := ws.Name + strconv.Itoa(rand.Intn(20)+1)
			if availableName[number] {
				i--
				continue
			}
			fakeMachine := models.SystemMachine{
				Analyze:       true,
				Brand:         strings.Title(strings.ToLower(babbler.Babble())),
				DcID:          0,
				ID:            0,
				MachineName:   strings.Title(strings.ToLower(babbler.Babble())),
				MachineNumber: number,
				Model:         strconv.Itoa((rand.Intn(50)+1)*10) + "T",
				WorkShopID:    ws.ID,
				MaterialType:  1,
			}
			var api restapitools.PostArg
			api.Body = fakeMachine
			api.IP = systemIP
			api.Token = jwt
			api.URL = systemMachine
			resp, err := api.Post()
			if err != nil {
				return err
			} else if resp != nil {
				defer resp.Body.Close()
			}
			availableName[number] = true
		}
	}
	return err
}

// GetSystemWorkShop GetSystemWorkShop
func GetSystemWorkShop() (workshop []models.WorkShop, err error) {
	var api restapitools.GetArg
	api.IP = systemIP
	api.URL = systemNewWorkShop
	api.Token = jwt
	resp, err := api.Get()
	if err != nil {
		return workshop, err
	} else if resp != nil {
		defer resp.Body.Close()
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return workshop, err
	}
	if err := json.Unmarshal(body, &workshop); err != nil {
		return workshop, err
	}
	return workshop, err
}

// NewWorkShop NewWorkShop
func NewWorkShop() (err error) {
	var wsNumSetting int
	if wsNumSetting, err = beego.AppConfig.Int("fakedata::wsNum"); err != nil {
		return err
	}
	var ws []models.WorkShop
	availableNumber := make(map[string]bool)
	if ws, err = GetSystemWorkShop(); err != nil {
		return err
	}
	for _, v := range ws {
		availableNumber[v.Number] = true
	}
	if len(ws) < wsNumSetting {
		for i := 0; i < wsNumSetting-len(ws); i++ {
			number := strconv.Itoa(rand.Intn(10) + 1)
			if availableNumber[number] {
				i--
				continue
			}
			babbler := babble.NewBabbler()
			babbler.Count = 1
			wsName := strings.Title(strings.ToLower(babbler.Babble()))
			insert := models.WorkShop{
				ID:     0,
				Name:   wsName,
				Number: number,
				Type:   1,
			}
			var api restapitools.PostArg
			api.Body = insert
			api.IP = systemIP
			api.Token = jwt
			api.URL = systemNewWorkShop
			resp, err := api.Post()
			if err != nil {
				return err
			} else if resp != nil {
				defer resp.Body.Close()
			}
			availableNumber[number] = true
		}
	}
	return err
}
