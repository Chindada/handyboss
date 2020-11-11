package tasks

import (
	"handyboss/fakedata"
	"handyboss/models"
	"runtime"

	"github.com/astaxie/beego"
)

func changeCycle() (err error) {
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
		planCycle := fakedata.PlanCycleTimeMap[k.MachineNumber]
		var machineFakeData models.LocalMachineList
		sqlite3db.Where("`mac_address` = ?", k.MacAddress).Last(&machineFakeData)
		sqlite3db.Model(&machineFakeData).Update("cycle_time", planCycle)
	}
	return err
}
