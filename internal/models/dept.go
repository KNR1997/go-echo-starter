package models

import "time"

type Department struct {
	ID        uint   `gorm:"primarykey"`
	Name      string `json:"name" gorm:"type:varchar(20);"`
	Desc      string `json:"desc" gorm:"type:varchar(500);"`
	Order     int    `json:"order" gorm:"type:int;"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
