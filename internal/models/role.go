package models

import "time"

type Role struct {
	ID          uint   `gorm:"primarykey"`
	Name        string `json:"name" gorm:"type:varchar(20);"`
	Description string `json:"description" gorm:"type:varchar(500);"`
	Menus       []Menu `json:"menus" gorm:"many2many:role_menus;"`
	Apis        []Api  `json:"apis" gorm:"many2many:role_apis;"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
