package domain

type UpdateApiRequest struct {
	ApiID   uint
	Path    string
	Method  string
	Summary string
	Tags    string
}

type DeleteApiRequest struct {
	ApiID uint
}
