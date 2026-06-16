package models

import "time"

type Menu struct {
	ID          uint      `json:"id" gorm:"primarykey"`
	Name        string    `json:"name" gorm:"type:varchar(20);not null"`
	Remark      *string   `json:"remark" gorm:"type:json"`                                    // Using *string to allow null values
	MenusType   *string   `json:"menus_type" gorm:"type:enum('catalog','menu');default:null"` // Using *string to allow null
	Icon        *string   `json:"icon" gorm:"type:varchar(100)"`
	Path        string    `json:"path" gorm:"type:varchar(100);not null"`
	OrderNumber int       `json:"order_number" gorm:"not null;default:0"`
	ParentID    int       `json:"parent_id" gorm:"not null;default:0"`
	IsHidden    bool      `json:"is_hidden" gorm:"not null;default:false"`
	Component   string    `json:"component" gorm:"type:varchar(100);not null"`
	Keepalive   bool      `json:"keepalive" gorm:"not null;default:true"`
	Redirect    *string   `json:"redirect" gorm:"type:varchar(100)"` // Using *string to allow null
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
