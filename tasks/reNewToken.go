package tasks

import (
	"emuMolding/fakedata"
)

func reNewToken() (err error) {
	err = fakedata.LoginSystem()
	return err
}
