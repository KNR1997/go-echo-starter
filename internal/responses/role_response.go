package responses

import "go-echo-starter/internal/models"

type RoleResponse struct {
	ID      uint      `json:"id" example:"1"`
	Name    string    `json:"name" example:"Admin"`
	Desc    string    `json:"desc" example:"some description"`
	MenuIDs []uint    `json:"menu_ids" example:"[1,2,3]"`
	Apis    []ApiInfo `json:"apis"`
}

type ApiInfo struct {
	ID     uint   `json:"id" example:"1"`
	Path   string `json:"path" example:"/api/v1/users"`
	Method string `json:"method" example:"GET"`
}

func NewRoleResponse(roles []models.Role) *[]RoleResponse {
	roleResponse := make([]RoleResponse, 0)

	for i := range roles {
		// Extract menu IDs from the role's menus
		menuIDs := make([]uint, len(roles[i].Menus))
		for j, menu := range roles[i].Menus {
			menuIDs[j] = menu.ID
		}

		// Extract menu IDs from the role's menus
		apis := make([]ApiInfo, len(roles[i].Apis))
		for j, api := range roles[i].Apis {
			apis[j] = ApiInfo{
				ID:     api.ID,
				Path:   api.Path,
				Method: api.Method,
			}
		}

		roleResponse = append(roleResponse, RoleResponse{
			ID:      roles[i].ID,
			Name:    roles[i].Name,
			Desc:    roles[i].Desc,
			MenuIDs: menuIDs,
			Apis:    apis,
		})
	}

	return &roleResponse
}

// Optional: If you need a single role response
func NewSingleRoleResponse(role models.Role) *RoleResponse {
	menuIDs := make([]uint, len(role.Menus))
	for i, menu := range role.Menus {
		menuIDs[i] = menu.ID
	}

	return &RoleResponse{
		ID:      role.ID,
		Name:    role.Name,
		Desc:    role.Desc,
		MenuIDs: menuIDs,
	}
}
