package controllers

import (
	"net/http"

	"github.com/astaxie/beego"
)

// BaseController BaseController
type BaseController struct {
	beego.Controller
	ControllerName string
	ActionName     string
}

// Prepare Prepare
func (c *BaseController) Prepare() {
	c.ControllerName, c.ActionName = c.GetControllerAndAction()
	if !c.AdapterUserInfo() {
		if c.ActionName == "Login" || c.ActionName == "NewUser" {
			return
		}
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.ServeJSON()
		c.StopRun()
	}
}

// AdapterUserInfo AdapterUserInfo
func (c *BaseController) AdapterUserInfo() bool {
	// k := c.GetSession("usersession")
	// if k != nil {
	// 	sessionmap := k.(map[string]bool)
	// 	if sessionmap[c.ActionName] || sessionmap["isSuper"] {
	// 		return true
	// 	}
	// }
	return true
}
