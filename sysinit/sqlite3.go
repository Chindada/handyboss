package sysinit

import (
	"emuMolding/models"

	"github.com/astaxie/beego"
	"github.com/jinzhu/gorm"

	// sqlite3 driver
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

// func init() {
// 	InitMachineSettings()
// 	InitSettings()
// }

// InitMachineSettings InitMachineSettings
func InitMachineSettings() {
	db, err := CreateMachineListConnection()
	if err != nil {
		beego.Alert(err)
	}
	defer db.Close()
	// db.LogMode(false)
	var resetIP bool
	if resetIP, err = beego.AppConfig.Bool("fakedata::resetIP"); err != nil {
		panic(err)
	} else if resetIP {
		beego.Informational("Reset Machine to Config")
		db.DropTable(&models.LocalMachineList{})
	}
	db.AutoMigrate(&models.LocalMachineList{})
}

// InitSettings InitSettings
func InitSettings() {
	db, err := CreateSettingsConnection()
	if err != nil {
		beego.Alert(err)
	}
	defer db.Close()
	// db.LogMode(false)
	var resetIP bool
	if resetIP, err = beego.AppConfig.Bool("fakedata::resetIP"); err != nil {
		panic(err)
	} else if resetIP {
		beego.Informational("Reset IP to Config")
		db.DropTable(&models.Settings{})
	}
	db.AutoMigrate(&models.Settings{})
}

// CreateMachineListConnection CreateMachineListConnection
func CreateMachineListConnection() (db *gorm.DB, err error) {
	db, err = gorm.Open("sqlite3", "./machine_list.db")
	if err != nil {
		beego.Alert(err)
	}
	return db, err
}

// CreateSettingsConnection CreateSettingsConnection
func CreateSettingsConnection() (db *gorm.DB, err error) {
	db, err = gorm.Open("sqlite3", "./settings.db")
	if err != nil {
		beego.Alert(err)
	}
	return db, err
}
