package tasks

import (
	"handyboss/fakedata"
)

func reNewToken() (err error) {
	err = fakedata.LoginSystem()
	return err
}
