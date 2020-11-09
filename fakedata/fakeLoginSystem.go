package fakedata

import (
	"emuMolding/libs/restapitools"
	"runtime"

	"github.com/astaxie/beego"
)

// LoginSystem LoginSystem
func LoginSystem() (err error) {
	defer func() {
		if r := recover(); r != nil {
			if _, ok := r.(runtime.Error); ok {
				panic(r)
			}
			err = r.(error)
			beego.Error(err)
		}
	}()
	var api restapitools.GetArg
	api.IP = systemIP
	api.URL = systemLoginURL
	headers := make(map[string]string)
	headers["account"] = beego.AppConfig.String("fakedata::account")
	headers["password"] = beego.AppConfig.String("fakedata::password")
	api.Headers = headers
	resp, err := api.Get()
	if err != nil {
		panic(err)
	}
	if resp != nil {
		defer resp.Body.Close()
		cookies := resp.Cookies()
		for _, k := range cookies {
			// beego.Informational(k.Value)
			jwt = k
		}
	}
	return err
}
