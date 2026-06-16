package domain

type UpdateRoleRequest struct {
	RoleID      uint
	Name        string
	Description string
}

type DeleteRoleRequest struct {
	RoleID uint
}

type AuthorizeRoleRequest struct {
	RoleID  uint
	MenuIDs []int
	ApiIDs  []int
}
