package companymodel

import "github.com/jinzhu/gorm"

// Company Company
type Company struct {
	gorm.Model `json:"-"`
	Number     string `gorm:"column:number;" json:"number"`
	Name       string `gorm:"column:name;" json:"name"`
}
