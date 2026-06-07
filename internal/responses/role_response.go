package responses

import "go-echo-starter/internal/models"

type RoleResponse struct {
	ID   uint   `json:"id" example:"1"`
	Name string `json:"name" example:"Admin"`
	Desc string `json:"desc" example:"some description"`
}

func NewRoleResponse(roles []models.Role) *[]RoleResponse {
	roleResponse := make([]RoleResponse, 0)

	for i := range roles {
		roleResponse = append(roleResponse, RoleResponse{
			ID:   roles[i].ID,
			Name: roles[i].Name,
			Desc: roles[i].Desc,
		})
	}

	return &roleResponse
}
