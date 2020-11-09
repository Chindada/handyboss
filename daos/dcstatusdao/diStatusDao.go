package dcstatusdao

import (
	"emuMolding/models"

	"github.com/astaxie/beego/orm"
)

// CreateOneStatus CreateOneStatus
func CreateOneStatus(status models.DcStatus) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(&status)
	return id, err
}

// CreateMultiStatus CreateMultiStatus
func CreateMultiStatus(statusArr []models.DcStatus) (ids int64, err error) {
	o := orm.NewOrm()
	ids, err = o.InsertMulti(100, statusArr)
	return ids, err
}

// ReadOneStatus ReadOneStatus
func ReadOneStatus(macAddress string) (status models.DcStatus, err error) {
	o := orm.NewOrm()
	err = o.QueryTable(status.TableName()).Filter("mac_address", macAddress).OrderBy("-timestamp").Limit(1).One(&status)
	return status, err
}
