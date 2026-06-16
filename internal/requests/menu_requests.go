package requests

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type BasicMenu struct {
	Name      string  `json:"name" validate:"required" example:"User Management"`
	Remark    *string `json:"remark" example:"{\"key\":\"value\"}"`
	MenusType *string `json:"menus_type" validate:"omitempty,oneof=catalog menu" example:"menu"`
	Icon      *string `json:"icon" example:"UserOutlined"`
	Path      string  `json:"path" validate:"required" example:"/users"`
	Order     int     `json:"order" example:"1"`
	ParentID  int     `json:"parent_id" example:"0"`
	IsHidden  bool    `json:"is_hidden" example:"false"`
	Component string  `json:"component" validate:"required" example:"views/UserList"`
	Keepalive bool    `json:"keepalive" example:"true"`
	Redirect  *string `json:"redirect" example:"/users/list"`
}

func (bm BasicMenu) Validate() error {
	return validation.ValidateStruct(&bm,
		validation.Field(&bm.Name, validation.Required, validation.Length(1, 20)),
		validation.Field(&bm.MenusType, validation.In("catalog", "menu")),
		validation.Field(&bm.Path, validation.Required, validation.Length(1, 100)),
		validation.Field(&bm.Order, validation.Min(0)),
		validation.Field(&bm.ParentID, validation.Min(0)),
		// validation.Field(&bm.Component, validation.Required, validation.Length(1, 100)),
		validation.Field(&bm.Icon, validation.Length(0, 100)),
		validation.Field(&bm.Redirect, validation.Length(0, 100)),
	)
}

type CreateMenuRequest struct {
	BasicMenu
}

type UpdateMenuRequest struct {
	BasicMenu
}

type PatchMenuRequest struct {
	IsHidden  *bool `json:"is_hidden" example:"false"`
	Keepalive *bool `json:"keepalive" example:"true"`
}
