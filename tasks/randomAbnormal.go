package tasks

import (
	"handyboss/fakedata"
	"handyboss/models"
	"math/rand"
	"runtime"

	"github.com/astaxie/beego"
)

func randomAbnormal() (err error) {
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
		probabilityA := rand.Intn(100) + 1
		if probabilityA <= 10 {
			sqlite3db.Model(&machineFakeData).Update("status", 3)
		}
	}
	return err
}
