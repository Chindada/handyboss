package tasks

import (
	"emuMolding/fakedata"
	"runtime"

	"github.com/astaxie/beego"
)

func addNewSchedule() (err error) {
	defer func() {
		if r := recover(); r != nil {
			if _, ok := r.(runtime.Error); ok {
				panic(r)
			}
			err = r.(error)
			beego.Error(err)
		}
	}()
	if lock, ok := fakedata.AddScheduleLock.LoadOrStore("active", false); ok {
		if lock.(bool) {
			beego.Informational("Last Insert Not yet")
			return
		}
		if err := fakedata.FakeNewSchedule(); err != nil {
			panic(err)
		} else {
			err = fakedata.CreateMultiScheduleFromSlice()
			if err != nil {
				panic(err)
			}
		}
	}
	return err
}
