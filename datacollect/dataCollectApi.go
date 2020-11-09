package datacollect

import (
	"emuMolding/libs/restapitools"
	"emuMolding/models"
	"encoding/json"
	"io/ioutil"
	"runtime"

	"github.com/astaxie/beego"
)

func checkGetDcLog(dc models.Wise) (tst, tend int64, err error) {
	defer func() {
		if r := recover(); r != nil {
			if _, ok := r.(runtime.Error); ok {
				panic(r)
			}
			err = r.(error)
			beego.Error(err)
		}
	}()
	headers := make(map[string]string)
	headers["Authorization"] = dc.Token
	api := restapitools.GetArg{
		URL:     "/log_output",
		IP:      dc.IP,
		Headers: headers,
		Tr:      dcTr,
	}
	resp, err := api.Get()
	if err != nil {
		panic(err)
	} else if resp != nil {
		defer resp.Body.Close()
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	var data models.DcGetLogBody
	if err := json.Unmarshal(body, &data); err != nil {
		panic(err)
	}
	tst = data.TSt
	tend = data.TEnd
	return tst, tend, err
}

func putDc(dc models.Wise) (err error) {
	defer func() {
		if r := recover(); r != nil {
			if _, ok := r.(runtime.Error); ok {
				panic(r)
			}
			err = r.(error)
			beego.Error(err)
		}
	}()
	puBody := models.DcPutBody{
		UID:  1,
		MAC:  0,
		TmF:  0,
		Fltr: 1,
	}
	headers := make(map[string]string)
	headers["Authorization"] = dc.Token
	api := restapitools.PutArg{
		URL:     "/log_output",
		IP:      dc.IP,
		Headers: headers,
		Body:    puBody,
		Tr:      dcTr,
	}
	resp, err := api.Put()
	if err != nil {
		panic(err)
	} else if resp != nil {
		defer resp.Body.Close()
	}
	return err
}

func patchDc(dc models.Wise, timeRange models.DcPatchBody) (err error) {
	defer func() {
		if r := recover(); r != nil {
			if _, ok := r.(runtime.Error); ok {
				panic(r)
			}
			err = r.(error)
			beego.Error(err)
		}
	}()
	headers := make(map[string]string)
	headers["Authorization"] = dc.Token
	api := restapitools.PatchArg{
		URL:     "/log_output",
		IP:      dc.IP,
		Headers: headers,
		Body:    timeRange,
		Tr:      dcTr,
	}
	resp, err := api.Patch()
	if err != nil {
		panic(err)
	} else if resp != nil {
		defer resp.Body.Close()
	}
	return err
}

func getDcLog(dc models.Wise) (data models.DcGetLogBody, err error) {
	defer func() {
		if r := recover(); r != nil {
			if _, ok := r.(runtime.Error); ok {
				panic(r)
			}
			err = r.(error)
			beego.Error(err)
		}
	}()
	headers := make(map[string]string)
	headers["Authorization"] = dc.Token
	api := restapitools.GetArg{
		URL:     "/log_output",
		IP:      dc.IP,
		Headers: headers,
		Tr:      dcTr,
	}
	resp, err := api.Get()
	if err != nil {
		panic(err)
	} else if resp != nil {
		defer resp.Body.Close()
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	if err := json.Unmarshal(body, &data); err != nil {
		panic(err)
	}
	return data, err
}

func getDcData(dc models.Wise) (dis []models.Di, err error) {
	defer func() {
		if r := recover(); r != nil {
			if _, ok := r.(runtime.Error); ok {
				panic(r)
			}
			err = r.(error)
			beego.Error(err)
		}
	}()
	headers := make(map[string]string)
	headers["Authorization"] = dc.Token
	api := restapitools.GetArg{
		URL:     "/log_message",
		IP:      dc.IP,
		Headers: headers,
		Tr:      dcTr,
	}
	resp, err := api.Get()
	if err != nil {
		panic(err)
	} else if resp != nil {
		defer resp.Body.Close()
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	var data models.WiseData
	if err := json.Unmarshal(body, &data); err != nil {
		panic(err)
	}
	dis, err = data.ConvertToDi()
	if err != nil {
		panic(err)
	}
	return dis, err
}
