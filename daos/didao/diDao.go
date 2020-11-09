package didao

import (
	"emuMolding/models"

	"github.com/astaxie/beego/orm"
)

// CreateOneDi CreateOneDi
func CreateOneDi(di models.Di) (id int64, err error) {
	o := orm.NewOrm()
	err = o.Begin()
	id, err = o.Insert(&di)
	if err != nil {
		err = o.Rollback()
	} else {
		err = o.Commit()
	}
	return id, err
}

// CreateMultiDi CreateMultiDi
func CreateMultiDi(dis []models.Di) (ids int64, err error) {
	o := orm.NewOrm()
	err = o.Begin()
	ids, err = o.InsertMulti(100, dis)
	if err != nil {
		err = o.Rollback()
	} else {
		err = o.Commit()
	}
	return ids, err
}

// ReadOneDi ReadOneDi
func ReadOneDi(macAddress string) (di models.Di, err error) {
	o := orm.NewOrm()
	err = o.Begin()
	err = o.QueryTable(di.TableName()).Filter("mac_address", macAddress).OrderBy("-timestamp").Limit(1).One(&di)
	if err != nil {
		err = o.Rollback()
	} else {
		err = o.Commit()
	}
	return di, err
}
