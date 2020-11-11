package dcstatusdao

import (
	"handyboss/models"

	"github.com/astaxie/beego/orm"
)

// CreateOneStatus CreateOneStatus
func CreateOneStatus(status models.DcStatus) (id int64, err error) {
	o := orm.NewOrm()
	err = o.Begin()
	id, err = o.Insert(&status)
	if err != nil {
		err = o.Rollback()
	} else {
		err = o.Commit()
	}
	return id, err
}

// CreateMultiStatus CreateMultiStatus
func CreateMultiStatus(statusArr []models.DcStatus) (ids int64, err error) {
	o := orm.NewOrm()
	err = o.Begin()
	ids, err = o.InsertMulti(100, statusArr)
	if err != nil {
		err = o.Rollback()
	} else {
		err = o.Commit()
	}
	return ids, err
}

// ReadOneStatus ReadOneStatus
func ReadOneStatus(macAddress string) (status models.DcStatus, err error) {
	o := orm.NewOrm()
	err = o.Begin()
	err = o.QueryTable(status.TableName()).Filter("mac_address", macAddress).OrderBy("-timestamp").Limit(1).One(&status)
	if err != nil {
		err = o.Rollback()
	} else {
		err = o.Commit()
	}
	return status, err
}
