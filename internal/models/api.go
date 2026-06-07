package models

import "time"

type Api struct {
	ID        uint   `gorm:"primarykey"`
	Path      string `json:"path" gorm:"type:varchar(100);"`
	Method    string `json:"method" gorm:"type:varchar(500);"`
	Summary   string `json:"summary" gorm:"type:varchar(500);"`
	Tags      string `json:"tags" gorm:"type:varchar(500);"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
