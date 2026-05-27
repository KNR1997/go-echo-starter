package models

import "gorm.io/gorm"

type Department struct {
	gorm.Model
	Name  string `json:"name" gorm:"type:varchar(20);"`
	Desc  string `json:"desc" gorm:"type:varchar(500);"`
	Order int    `json:"order" gorm:"type:int;"`
}
