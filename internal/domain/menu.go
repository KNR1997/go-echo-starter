package domain

type UpdateMenuRequest struct {
	MenuID    uint
	Name      string
	Remark    *string
	MenusType *string
	Icon      *string
	Path      string
	Order     int
	ParentID  int
	IsHidden  bool
	Component string
	Keepalive bool
	Redirect  *string
}

type PatchMenuRequest struct {
	MenuID    uint
	IsHidden  *bool
	Keepalive *bool
}

type DeleteMenuRequest struct {
	MenuID uint
}
