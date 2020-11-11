package machinecontroller

import (
	"encoding/json"
	"handyboss/controllers"
	"handyboss/fakedata"
	"handyboss/models"
	"net/http"

	"github.com/astaxie/beego"
)

// MachineController MachineController
type MachineController struct {
	controllers.BaseController
}

// ChangeAction ChangeAction
func (c *MachineController) ChangeAction() {
	var machineAction models.MachineAction
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &machineAction); err != nil {
		beego.Informational(machineAction)
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Data["json"] = err.Error()
		c.ServeJSON()
	} else if err := fakedata.TaskTime(machineAction); err != nil {
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Data["json"] = err.Error()
		c.ServeJSON()
	} else {
		c.Ctx.Output.SetStatus(http.StatusOK)
	}
}

// GetMachineDetail GetMachineDetail
func (c *MachineController) GetMachineDetail() {
	machineNumber := c.Ctx.Input.Header("machineNumber")
	data, err := fakedata.GetMachineDetail(machineNumber)
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Data["json"] = err.Error()
		c.ServeJSON()
	} else {
		c.Ctx.Output.SetStatus(http.StatusOK)
		c.Data["json"] = data
		c.ServeJSON()
	}
}

// UpdateMachine UpdateMachine
func (c *MachineController) UpdateMachine() {
	ip := c.Ctx.Input.Header("ip")
	workshop := c.Ctx.Input.Header("workShop")

	if err := fakedata.UpdateIP(ip); err != nil {
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Data["json"] = err.Error()
		c.ServeJSON()
		return
	} else if err := fakedata.UpdateWorkShop(workshop); err != nil {
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Data["json"] = err.Error()
		c.ServeJSON()
	} else {
		c.Ctx.Output.SetStatus(http.StatusOK)
		// fakedata.GetIP()
		// fakedata.GetWorkShop()
	}
}

// GetSettings GetSettings
func (c *MachineController) GetSettings() {
	data, err := fakedata.GetSettings()
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Data["json"] = err.Error()
		c.ServeJSON()
	} else {
		c.Ctx.Output.SetStatus(http.StatusOK)
		c.Data["json"] = data
		c.ServeJSON()
	}
}
