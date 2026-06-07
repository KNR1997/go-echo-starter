package responses

import "go-echo-starter/internal/models"

type MenuResponse struct {
	ID   uint   `json:"id" example:"1"`
	Name string `json:"name" example:"/roles/get"`
}

func NewMenuResponse(menus []models.Menu) *[]MenuResponse {
	menuResponse := make([]MenuResponse, 0)

	for i := range menus {
		menuResponse = append(menuResponse, MenuResponse{
			ID:   menus[i].ID,
			Name: menus[i].Name,
		})
	}

	return &menuResponse
}
