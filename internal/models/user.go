package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username    string `json:"username" gorm:"type:varchar(20);"`
	Email       string `json:"email" gorm:"type:varchar(200);"`
	Name        string `json:"name" gorm:"type:varchar(200);"`
	Phone       string `json:"phone" gorm:"type:varchar(20);"`
	Password    string `json:"password" gorm:"type:varchar(200);"`
	Post        []Post
	IsActive    bool        `json:"is_active" gorm:"type:boolean;"`
	DeptId      *int        `json:"dept_id" gorm:"type:varchar(20);"`
	IsSuperUser bool        `json:"is_superuser" gorm:"column:is_superuser;type:boolean;default:false"`
	Roles       []Role      `json:"roles" gorm:"many2many:user_role;"`
	Department  *Department `json:"department" gorm:"foreignKey:DeptId;references:ID"`
}
