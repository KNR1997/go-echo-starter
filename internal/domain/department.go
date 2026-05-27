package domain

type UpdateDepartmentRequest struct {
	DeptID uint
	Name   string
	Desc   string
}

type DeleteDepartmentRequest struct {
	DeptID uint
}
