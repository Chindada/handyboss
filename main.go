package main

import (
	"handyboss/fakedata"
	_ "handyboss/libs/restapitools"
	_ "handyboss/routers"
	_ "handyboss/sysinit"
	_ "handyboss/tasks"

	"github.com/astaxie/beego"
)

func main() {
	fakedata.NewMachineLoop()
	fakedata.NewMold()
	go fakedata.Loop()
	// go datacollect.FetchLoop(datacollect.Dcs)
	// test two remote
	beego.Run()
}
