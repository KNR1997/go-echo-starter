package responses

import (
	"go-echo-starter/internal/models"
	"time"
)

type MenuResponse struct {
	ID          uint           `json:"id" example:"1"`
	Name        string         `json:"name" example:"/roles/get"`
	Remark      *string        `json:"remark" example:"{\"key\":\"value\"}"`
	MenusType   *string        `json:"menu" example:"menu"`
	Icon        *string        `json:"icon" example:"fa-user"`
	Path        string         `json:"path" example:"/roles"`
	OrderNumber int            `json:"order" example:"1"`
	ParentID    int            `json:"parent_id" example:"0"`
	IsHidden    bool           `json:"is_hidden" example:"false"`
	Component   string         `json:"component" example:"views/RoleList"`
	Keepalive   bool           `json:"keepalive" example:"true"`
	Redirect    *string        `json:"redirect" example:"/dashboard"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	Children    []MenuResponse `json:"children,omitempty"` // omitempty to hide empty children
}

// NewMenuResponse returns a flat list of menu responses
func NewMenuResponse(menus []models.Menu) *[]MenuResponse {
	menuResponse := make([]MenuResponse, 0)

	for i := range menus {
		menuResponse = append(menuResponse, MenuResponse{
			ID:          menus[i].ID,
			Name:        menus[i].Name,
			Remark:      menus[i].Remark,
			MenusType:   menus[i].MenusType,
			Icon:        menus[i].Icon,
			Path:        menus[i].Path,
			OrderNumber: menus[i].OrderNumber,
			ParentID:    menus[i].ParentID,
			IsHidden:    menus[i].IsHidden,
			Component:   menus[i].Component,
			Keepalive:   menus[i].Keepalive,
			Redirect:    menus[i].Redirect,
			CreatedAt:   menus[i].CreatedAt,
			UpdatedAt:   menus[i].UpdatedAt,
			Children:    []MenuResponse{}, // Initialize empty slice
		})
	}

	return &menuResponse
}

// NewMenuTreeResponse returns a tree structure of menus (parent with children)
func NewMenuTreeResponse(menus []models.Menu) *[]MenuResponse {
	// First, convert all menus to MenuResponse
	menuMap := make(map[uint]MenuResponse)
	menuResponseList := make([]MenuResponse, 0)

	// Create a map for quick lookup
	for i := range menus {
		menuResponse := MenuResponse{
			ID:          menus[i].ID,
			Name:        menus[i].Name,
			Remark:      menus[i].Remark,
			MenusType:   menus[i].MenusType,
			Icon:        menus[i].Icon,
			Path:        menus[i].Path,
			OrderNumber: menus[i].OrderNumber,
			ParentID:    menus[i].ParentID,
			IsHidden:    menus[i].IsHidden,
			Component:   menus[i].Component,
			Keepalive:   menus[i].Keepalive,
			Redirect:    menus[i].Redirect,
			CreatedAt:   menus[i].CreatedAt,
			UpdatedAt:   menus[i].UpdatedAt,
			Children:    []MenuResponse{},
		}
		menuMap[menus[i].ID] = menuResponse
		menuResponseList = append(menuResponseList, menuResponse)
	}

	// Build tree structure - identify parents first
	var rootMenus []MenuResponse

	// First pass: collect all root menus and ensure parents exist in map
	for i := range menuResponseList {
		menu := &menuResponseList[i]
		if menu.ParentID == 0 {
			rootMenus = append(rootMenus, *menu)
		}
	}

	// Second pass: add children to their parents
	for i := range menuResponseList {
		menu := &menuResponseList[i]
		if menu.ParentID != 0 {
			// Find parent and add as child
			if parent, exists := menuMap[uint(menu.ParentID)]; exists {
				// Update parent's children
				parentCopy := parent
				parentCopy.Children = append(parentCopy.Children, *menu)
				menuMap[uint(menu.ParentID)] = parentCopy

				// Update in rootMenus if parent is a root menu
				for j := range rootMenus {
					if rootMenus[j].ID == parent.ID {
						rootMenus[j] = parentCopy
						break
					}
				}

				// Also update in menuResponseList to keep data consistent
				for k := range menuResponseList {
					if menuResponseList[k].ID == parent.ID {
						menuResponseList[k] = parentCopy
						break
					}
				}
			}
		}
	}

	return &rootMenus
}

// Alternative: More efficient tree building using recursion
func NewMenuTreeResponseOptimized(menus []models.Menu) *[]MenuResponse {
	// Create a map of parent_id to list of menus
	menuMap := make(map[int][]models.Menu)
	menuResponseMap := make(map[uint]*MenuResponse)

	// Group menus by parent_id
	for i := range menus {
		menu := menus[i]
		menuMap[menu.ParentID] = append(menuMap[menu.ParentID], menu)
	}

	// Recursive function to build tree
	var buildTree func(parentID int) []MenuResponse
	buildTree = func(parentID int) []MenuResponse {
		var result []MenuResponse

		for _, menu := range menuMap[parentID] {
			menuResponse := MenuResponse{
				ID:          menu.ID,
				Name:        menu.Name,
				Remark:      menu.Remark,
				MenusType:   menu.MenusType,
				Icon:        menu.Icon,
				Path:        menu.Path,
				OrderNumber: menu.OrderNumber,
				ParentID:    menu.ParentID,
				IsHidden:    menu.IsHidden,
				Component:   menu.Component,
				Keepalive:   menu.Keepalive,
				Redirect:    menu.Redirect,
				CreatedAt:   menu.CreatedAt,
				UpdatedAt:   menu.UpdatedAt,
				Children:    buildTree(int(menu.ID)), // Recursively get children
			}
			result = append(result, menuResponse)
			menuResponseMap[menu.ID] = &result[len(result)-1]
		}

		return result
	}

	// Build tree starting from root (parent_id = 0)
	rootMenus := buildTree(0)
	return &rootMenus
}
