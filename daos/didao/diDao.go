package didao

import (
	"emuMolding/models"

	"github.com/astaxie/beego/orm"
)

// CreateOneDi CreateOneDi
func CreateOneDi(di models.Di) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(&di)
	return id, err
}

// CreateMultiDi CreateMultiDi
func CreateMultiDi(dis []models.Di) (ids int64, err error) {
	o := orm.NewOrm()
	ids, err = o.InsertMulti(100, dis)
	return ids, err
}

// ReadOneDi ReadOneDi
func ReadOneDi(macAddress string) (di models.Di, err error) {
	o := orm.NewOrm()
	err = o.QueryTable(di.TableName()).Filter("mac_address", macAddress).OrderBy("-timestamp").Limit(1).One(&di)
	return di, err
}
