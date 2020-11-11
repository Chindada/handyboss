package tasks

import (
	"handyboss/sysinit"

	"github.com/astaxie/beego/toolbox"
	"github.com/jinzhu/gorm"
)

var sqlite3db *gorm.DB

func init() {
	var err error
	sqlite3db, err = sysinit.CreateMachineListConnection()
	if err != nil {
		panic(err)
	}

	changeStable := toolbox.NewTask("changeStable", "0 0/5 * * * *", changeStable)
	toolbox.AddTask("changeStable", changeStable)

	printMemUsage := toolbox.NewTask("printMemUsage", "0 0/3 * * * *", printMemUsage)
	toolbox.AddTask("printMemUsage", printMemUsage)

	addNewSchedule := toolbox.NewTask("addNewSchedule", "0 0/1 * * * *", addNewSchedule)
	toolbox.AddTask("addNewSchedule", addNewSchedule)

	randomAbnormal := toolbox.NewTask("randomAbnormal", "0 0/15 * * * *", randomAbnormal)
	toolbox.AddTask("randomAbnormal", randomAbnormal)

	reNewToken := toolbox.NewTask("reNewToken", "0 0/10 * * * *", reNewToken)
	toolbox.AddTask("reNewToken", reNewToken)

	toolbox.StartTask()
}
