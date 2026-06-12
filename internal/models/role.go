package models

import "time"

type Role struct {
	ID        uint   `gorm:"primarykey"`
	Name      string `json:"name" gorm:"type:varchar(20);"`
	Desc      string `json:"desc" gorm:"type:varchar(500);"`
	Menus     []Menu `json:"menus" gorm:"many2many:role_menu;"`
	Apis      []Api  `json:"apis" gorm:"many2many:role_api;"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
