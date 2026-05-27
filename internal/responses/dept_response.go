package responses

import "go-echo-starter/internal/models"

type DeptResponse struct {
	Name string `json:"name" example:"John Doe"`
	Desc string `json:"phone" example:"0113123888"`
}

func NewDeptResponse(roles []models.Department) *[]DeptResponse {
	roleResponse := make([]DeptResponse, 0)

	for i := range roles {
		roleResponse = append(roleResponse, DeptResponse{
			Name: roles[i].Name,
			Desc: roles[i].Desc,
		})
	}

	return &roleResponse
}
