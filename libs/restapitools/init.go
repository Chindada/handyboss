package restapitools

import "net/http"

var tr *http.Transport

func init() {
	tr = &http.Transport{
		DisableKeepAlives: false,
		MaxIdleConns:      100,
	}
}
