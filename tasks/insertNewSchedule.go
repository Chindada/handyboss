package tasks

import (
	"emuMolding/fakedata"

	"github.com/astaxie/beego"
)

func addNewSchedule() (err error) {
	if lock, ok := fakedata.AddScheduleLock.LoadOrStore("active", false); ok {
		if lock.(bool) {
			beego.Informational("Last Insert Not yet")
			return
		}
		fakedata.FakeNewSchedule()
		fakedata.CreateMultiScheduleFromSlice()
	}
	return err
}
