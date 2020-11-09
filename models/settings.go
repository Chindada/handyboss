package models

import (
	"strings"

	"github.com/jinzhu/gorm"
)

// Settings Settings
type Settings struct {
	gorm.Model `json:"-"`
	Key        string `gorm:"column:key;" json:"key"`
	Value      string `gorm:"column:value;" json:"value"`
}

// GetWorkShops GetWorkShops
func (c Settings) GetWorkShops() []string {
	return strings.Split(c.Value, ",")
}
