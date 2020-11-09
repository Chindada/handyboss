package tasks

import (
	"emuMolding/sysinit"

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

	// changeStable := toolbox.NewTask("changeStable", "0 0 0/2 * * *", changeStable)
	// toolbox.AddTask("changeStable", changeStable)

	printMemUsage := toolbox.NewTask("printMemUsage", "0 0/1 * * * *", printMemUsage)
	toolbox.AddTask("printMemUsage", printMemUsage)

	// changeCycle := toolbox.NewTask("changeCycle", "0/15 * * * * *", changeCycle)
	// toolbox.AddTask("changeCycle", changeCycle)

	addNewSchedule := toolbox.NewTask("addNewSchedule", "0/10 * * * * *", addNewSchedule)
	toolbox.AddTask("addNewSchedule", addNewSchedule)

	// randomAbnormal := toolbox.NewTask("randomAbnormal", "0 0/10 * * * *", randomAbnormal)
	// toolbox.AddTask("randomAbnormal", randomAbnormal)

	reNewToken := toolbox.NewTask("reNewToken", "0 0/10 * * * *", reNewToken)
	toolbox.AddTask("reNewToken", reNewToken)

	toolbox.StartTask()
}
