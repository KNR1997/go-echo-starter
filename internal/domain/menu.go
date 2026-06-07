package domain

type UpdateMenuRequest struct {
	MenuID uint
	Name   string
}

type DeleteMenuRequest struct {
	MenuID uint
}
