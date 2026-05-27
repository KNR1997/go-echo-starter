package responses

import "go-echo-starter/internal/models"

type RoleResponse struct {
	Name string `json:"name" example:"John Doe"`
	Desc string `json:"phone" example:"0113123888"`
}

func NewRoleResponse(roles []models.Role) *[]RoleResponse {
	roleResponse := make([]RoleResponse, 0)

	for i := range roles {
		roleResponse = append(roleResponse, RoleResponse{
			Name: roles[i].Name,
			Desc: roles[i].Desc,
		})
	}

	return &roleResponse
}
