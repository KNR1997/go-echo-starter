package models

import "gorm.io/gorm"

type Api struct {
	gorm.Model
	Path    string `json:"path" gorm:"type:varchar(100);"`
	Method  string `json:"method" gorm:"type:varchar(500);"`
	Summary string `json:"summary" gorm:"type:varchar(500);"`
	Tags    string `json:"tags" gorm:"type:varchar(500);"`
}
