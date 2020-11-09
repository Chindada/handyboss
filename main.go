package main

import (
	"emuMolding/datacollect"

	// _ "emuMolding/routers"
	// _ "emuMolding/tasks"

	_ "emuMolding/libs/restapitools"
	_ "emuMolding/sysinit"

	"github.com/astaxie/beego"
)

func main() {
	// fakedata.NewMachineLoop()
	// fakedata.NewMold()
	// go fakedata.Loop()
	// datacollect.FetchDc()
	go datacollect.FetchLoop(datacollect.Dcs)
	beego.Run()
}
