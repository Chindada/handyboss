package main

import (
	"emuMolding/datacollect"
	"emuMolding/models"

	// _ "emuMolding/routers"
	// _ "emuMolding/tasks"

	_ "emuMolding/libs/restapitools"
	_ "emuMolding/sysinit"

	"github.com/astaxie/beego"
)

func main() {
	dcA := models.Wise{
		Token:           "Basic cm9vdDptaXRyb290",
		MacAddress:      "00D0C9E34A12",
		IdleTime:        300,
		PutTimeInterval: 1800,
		IP:              "192.168.10.119",
	}
	dcB := models.Wise{
		Token:           "Basic cm9vdDptaXRyb290",
		MacAddress:      "00D0C9E50453",
		IdleTime:        300,
		PutTimeInterval: 1800,
		IP:              "192.168.10.192",
	}
	dcC := models.Wise{
		Token:           "Basic cm9vdDptaXRyb290",
		MacAddress:      "00D0C9E349F4",
		IdleTime:        300,
		PutTimeInterval: 1800,
		IP:              "192.168.10.145",
	}
	var dcs []models.Wise
	dcs = append(dcs, dcA, dcB, dcC)

	// fakedata.NewMachineLoop()
	// fakedata.NewMold()
	// go fakedata.Loop()
	// datacollect.FetchDc()
	go datacollect.FetchLoop(dcs)
	beego.Run()
}
