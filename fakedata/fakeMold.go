package fakedata

import (
	"emuMolding/libs/restapitools"
	"emuMolding/models"
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"runtime"
	"strconv"
	"strings"

	"github.com/astaxie/beego"
	"github.com/tjarratt/babble"
)

// NewMold NewMold
func NewMold() (err error) {
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
	var num int
	if insertnew, err = beego.AppConfig.Bool("fakedata::insertNewMold"); !insertnew {
		return
	}
	GetSystemMold()
	num, err = beego.AppConfig.Int("fakedata::moldNum")
	if err != nil {
		panic(err)
	}
	if len(allMolds) > num {
		beego.Informational("Mold Num is enough.")
		return
	}
	var mold []models.Mold
	babbler := babble.NewBabbler()
	babbler.Count = 1

	for i := 1; i <= num; i++ {
		in := []int64{1, 2, 4, 8, 16, 32}
		randomIndex := rand.Intn(len(in))
		tempMold := models.Mold{
			CavityNumber: in[randomIndex],
			CycleTime:    rand.Int63n(40) + 10,
			UpTime:       60,
			DownTime:     60,
			Name:         strings.Title(strings.ToLower(babbler.Babble())),
			ProductModel: strings.Title(strings.ToLower(babbler.Babble())),
			Number:       "#" + strconv.Itoa(10000+i),
			GreenRange:   0.1,
			YellowRange:  0.2,
		}
		mold = append(mold, tempMold)
		if len(mold)%100 == 0 {
			var api restapitools.PostArg
			api.IP = systemIP
			api.URL = systemMold
			api.Token = jwt
			headers := make(map[string]string)
			headers["multi"] = "true"
			api.Headers = headers
			api.Body = mold
			resp, err := api.Post()
			if err != nil {
				panic(err)
			} else if resp != nil {
				defer resp.Body.Close()
			}
			mold = nil
		}
	}
	var api restapitools.PostArg
	api.IP = systemIP
	api.URL = systemMold
	api.Token = jwt
	headers := make(map[string]string)
	headers["multi"] = "true"
	api.Headers = headers
	api.Body = mold
	resp, err := api.Post()
	if err != nil {
		panic(err)
	} else if resp != nil {
		defer resp.Body.Close()
	}
	mold = nil
	if err := CreateMoldMachineRel(); err != nil {
		panic(err)
	}
	return err
}

// CreateMoldMachineRel CreateMoldMachineRel
func CreateMoldMachineRel() (err error) {
	defer func() {
		if r := recover(); r != nil {
			if _, ok := r.(runtime.Error); ok {
				panic(r)
			}
			err = r.(error)
			beego.Error(err)
		}
	}()
	if err := GetSystemMachine(); err != nil {
		panic(err)
	}
	if err := GetSystemMold(); err != nil {
		panic(err)
	}
	var moldMachineRel []models.MoldMachineRel
	for _, i := range systemMachineID {
		for _, j := range systemMoldID {
			temp := models.MoldMachineRel{
				MachineID: i,
				MoldID:    j,
			}
			moldMachineRel = append(moldMachineRel, temp)
			if len(moldMachineRel)%500 == 0 {
				var api restapitools.PutArg
				api.IP = systemIP
				api.URL = systemMoldMachineRel
				api.Token = jwt
				api.Body = moldMachineRel
				resp, err := api.Put()
				if err != nil {
					panic(err)
				} else if resp != nil {
					defer resp.Body.Close()
				}
				moldMachineRel = nil
			}
		}
		var api restapitools.PutArg
		api.IP = systemIP
		api.URL = systemMoldMachineRel
		api.Token = jwt
		api.Body = moldMachineRel
		resp, err := api.Put()
		if err != nil {
			panic(err)
		} else if resp != nil {
			defer resp.Body.Close()
		}
		moldMachineRel = nil
	}
	return err
}

// GetSystemMachine GetSystemMachine
func GetSystemMachine() (err error) {
	var api restapitools.GetArg
	api.IP = systemIP
	api.URL = systemMachine
	api.Token = jwt
	resp, err := api.Get()
	if err != nil {
		return err
	} else if resp != nil {
		defer resp.Body.Close()
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var machine []models.SystemMachine
	if err := json.Unmarshal(body, &machine); err != nil {
		return err
	}
	systemMachineDetail = machine
	systemMachineID = nil
	for _, v := range machine {
		systemMachineID = append(systemMachineID, v.ID)
	}
	return err
}

// GetSystemMold GetSystemMold
func GetSystemMold() (err error) {
	var api restapitools.GetArg
	api.IP = systemIP
	api.URL = systemMold
	api.Token = jwt
	resp, err := api.Get()
	if err != nil {
		return err
	} else if resp != nil {
		defer resp.Body.Close()
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(body, &allMolds); err != nil {
		return err
	}
	for _, v := range allMolds {
		systemMoldID = append(systemMoldID, v.ID)
	}
	return err
}
