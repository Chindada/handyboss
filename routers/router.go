package routers

import (
	"handyboss/controllers/machinecontroller"

	"github.com/astaxie/beego"
)

// CookieJWTage CookieJWTage
// const CookieJWTage int64 = 3600

const (
	actionMethod   = "post:ChangeAction;get:GetMachineDetail"
	settingsMethod = "patch:UpdateMachine;get:GetSettings"
)

func init() {
	// Filter for jwt
	// beego.InsertFilter("/api/*", beego.BeforeRouter, func(ctx *context.Context) {
	// 	cookie, err := ctx.Request.Cookie("jwt")
	// 	if err != nil {
	// 		utils.LogError(err)
	// 	} else if username, ok := utils.CheckToken(cookie.Value); ok {
	// 		token := utils.GenToken(username, CookieJWTage)
	// 		ctx.SetCookie("jwt", token, CookieJWTage)
	// 	} else {
	// 		ctx.Redirect(500, "/")
	// 		beego.Informational("Not OK")
	// 	}
	// })
	beego.Router("/iomfake",
		&machinecontroller.MachineController{},
		actionMethod)
	beego.Router("/iomfake/settings",
		&machinecontroller.MachineController{},
		settingsMethod)
}
