package responses

import "go-echo-starter/internal/models"

type DeptResponse struct {
	ID   uint   `json:"id" example:"1"`
	Name string `json:"name" example:"IT"`
	Desc string `json:"desc" example:"notes"`
}

func NewDeptResponse(departments []models.Department) *[]DeptResponse {
	deptResponse := make([]DeptResponse, 0)

	for i := range departments {
		deptResponse = append(deptResponse, DeptResponse{
			ID:   departments[i].ID,
			Name: departments[i].Name,
			Desc: departments[i].Desc,
		})
	}

	return &deptResponse
}
