package models

import "gorm.io/gorm"

type Role struct {
	gorm.Model
	Name string `json:"name" gorm:"type:varchar(20);"`
	Desc string `json:"desc" gorm:"type:varchar(500);"`
}
