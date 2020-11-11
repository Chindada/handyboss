package datacollect

import (
	"encoding/json"
	"handyboss/libs/restapitools"
	"handyboss/models"
	"io/ioutil"
)

func checkGetDcLog(dc models.Wise) (tst, tend int64, err error) {
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
		return tst, tend, err
	} else if resp != nil {
		defer resp.Body.Close()
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return tst, tend, err
	}

	var data models.DcGetLogBody
	if err := json.Unmarshal(body, &data); err != nil {
		return tst, tend, err
	}
	tst = data.TSt
	tend = data.TEnd
	return tst, tend, err
}

func putDc(dc models.Wise) (err error) {
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
		return err
	} else if resp != nil {
		defer resp.Body.Close()
	}
	return err
}

func patchDc(dc models.Wise, timeRange models.DcPatchBody) (err error) {
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
		return err
	} else if resp != nil {
		defer resp.Body.Close()
	}
	return err
}

func getDcLog(dc models.Wise) (data models.DcGetLogBody, err error) {
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
		return data, err
	}
	if err := json.Unmarshal(body, &data); err != nil {
		return data, err
	}
	return data, err
}

func getDcData(dc models.Wise) (dis []models.Di, err error) {
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
		return dis, err
	}
	var data models.WiseData
	if err := json.Unmarshal(body, &data); err != nil {
		return dis, err
	}
	dis, err = data.ConvertToDi()
	if err != nil {
		return dis, err
	}
	return dis, err
}
