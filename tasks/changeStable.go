package tasks

import (
	"handyboss/fakedata"
	"handyboss/models"
	"math/rand"
	"runtime"

	"github.com/astaxie/beego"
)

func changeStable() (err error) {
	defer func() {
		if r := recover(); r != nil {
			if _, ok := r.(runtime.Error); ok {
				panic(r)
			}
			err = r.(error)
			beego.Error(err)
		}
	}()
	param := fakedata.Machine
	for _, k := range param {
		var machineFakeData models.LocalMachineList
		sqlite3db.Where("`mac_address` = ?", k.MacAddress).Last(&machineFakeData)
		var idleFirst bool
		if idleFirst, err = beego.AppConfig.Bool("fakedata::idleFirst"); err != nil {
			panic(err)
		}
		probabilityA := rand.Intn(100) + 1
		if probabilityA <= 60 {
			sqlite3db.Model(&machineFakeData).Update("stable", false)
			sqlite3db.Model(&machineFakeData).Update("idle", true)
		} else {
			sqlite3db.Model(&machineFakeData).Update("stable", true)
			if !idleFirst {
				sqlite3db.Model(&machineFakeData).Update("idle", false)
			}
		}
	}
	return err
}
