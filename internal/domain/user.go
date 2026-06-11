package domain

type UpdateUserRequest struct {
	UserID      uint
	UserName    string
	Name        string
	Email       string
	IsSuperUser bool
	IsActive    bool
	RoleIds     []int
	DeptId      int
}

type DeleteUserRequest struct {
	UserID uint
}
