package domain

type UpdateRoleRequest struct {
	RoleID uint
	Name   string
	Desc   string
}

type DeleteRoleRequest struct {
	RoleID uint
}
