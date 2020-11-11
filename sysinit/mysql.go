package sysinit

import (
	"handyboss/models"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"

	// orm mysql driver
	_ "github.com/go-sql-driver/mysql"
)

func init() {
	name := "default"
	dbUser := beego.AppConfig.String("dbUser")
	dbPasswd := beego.AppConfig.String("dbPasswd")
	dbHost := beego.AppConfig.String("dbHost")
	dbPort := beego.AppConfig.String("dbPort")
	dbTimeZone := beego.AppConfig.String("dbTimeZone")
	orm.RegisterDriver("mysql", orm.DRMySQL)
	orm.RegisterDataBase(name, "mysql", dbUser+":"+dbPasswd+"@tcp("+dbHost+":"+dbPort+")/fakedata?charset=utf8&loc="+dbTimeZone)

	orm.RegisterModel(new(models.Di))
	orm.RegisterModel(new(models.DcStatus))
	orm.RegisterModel(new(models.Wise))

	err := orm.RunSyncdb(name, false, true)
	if err != nil {
		panic(err)
	}
}
