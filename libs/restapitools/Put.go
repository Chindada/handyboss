package restapitools

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

// PutArg PutArg
type PutArg struct {
	Token   *http.Cookie
	URL     string
	IP      string
	Headers map[string]string
	Body    interface{}
	Tr      *http.Transport
}

// Put Put
func (c *PutArg) Put() (resp *http.Response, err error) {
	url := "http://" + c.IP + c.URL
	b, err := json.Marshal(c.Body)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}
	if c.Token != nil {
		req.AddCookie(c.Token)
	}
	req.Header.Set("Content-Type", "application/json")
	if len(c.Headers) != 0 {
		for k, v := range c.Headers {
			req.Header.Set(k, v)
		}
	}
	if c.Tr == nil {
		c.Tr = tr
	}
	client := &http.Client{Transport: c.Tr}
	resp, err = client.Do(req)
	if err != nil {
		return nil, err
	} else if resp.StatusCode != 200 {
		defer resp.Body.Close()
		responseData, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		err = errors.New(string(responseData))
		return nil, err
	}
	return resp, err
}
