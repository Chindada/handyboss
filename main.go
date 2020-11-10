package main

import (
	"emuMolding/fakedata"
	_ "emuMolding/libs/restapitools"
	_ "emuMolding/routers"
	_ "emuMolding/sysinit"
	_ "emuMolding/tasks"

	"github.com/astaxie/beego"
)

func main() {
	fakedata.NewMachineLoop()
	fakedata.NewMold()
	go fakedata.Loop()
	// go datacollect.FetchLoop(datacollect.Dcs)
	beego.Run()
}
