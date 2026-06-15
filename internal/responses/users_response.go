// responses/user_response.go
package responses

import (
	"go-echo-starter/internal/models"
	"time"
)

// DepartmentResponse represents the department data in the response
type DepartmentResponse struct {
	ID   uint   `json:"id" example:"1"`
	Name string `json:"name" example:"Engineering"`
	Code string `json:"code,omitempty" example:"ENG"`
}

// UserResponse represents the user data in the response
type UserResponse struct {
	ID          uint                `json:"id" example:"1"`
	Username    string              `json:"username" example:"john"`
	Email       string              `json:"email" example:"admin@demo.com"`
	Name        string              `json:"name" example:"John Doe"`
	Phone       string              `json:"phone" example:"0113123888"`
	IsActive    bool                `json:"is_active" example:"true"`
	IsSuperUser bool                `json:"is_superuser" example:"false"`
	Roles       []RoleResponse      `json:"roles"`
	Department  *DepartmentResponse `json:"department,omitempty"` // Use omitempty to handle nil
	LastLogin   *time.Time          `json:"last_login"`
}

func NewUserResponse(users []models.User) *[]UserResponse {
	userResponse := make([]UserResponse, 0, len(users))

	for i := range users {
		// Convert roles
		roles := make([]RoleResponse, 0, len(users[i].Roles))
		for _, role := range users[i].Roles {
			roles = append(roles, RoleResponse{
				ID:   role.ID,
				Name: role.Name,
			})
		}

		// Convert department (if exists)
		var department *DepartmentResponse
		if users[i].Department != nil {
			department = &DepartmentResponse{
				ID:   users[i].Department.ID,
				Name: users[i].Department.Name,
			}
		}

		userResponse = append(userResponse, UserResponse{
			ID:          users[i].ID,
			Username:    users[i].Username,
			Email:       users[i].Email,
			Name:        users[i].Name,
			Phone:       users[i].Phone,
			IsActive:    users[i].IsActive,
			IsSuperUser: users[i].IsSuperUser,
			Roles:       roles,
			Department:  department,
			LastLogin:   users[i].LastLogin,
		})
	}

	return &userResponse
}

// For single user response
func NewSingleUserResponse(user *models.User) *UserResponse {
	// Convert roles
	// roles := make([]RoleResponse, 0, len(user.Roles))
	// for _, role := range user.Roles {
	// 	roles = append(roles, RoleResponse{
	// 		ID:   role.ID,
	// 		Name: role.Name,
	// 	})
	// }

	// Convert department
	// var department *DepartmentResponse
	// if user.Department != nil {
	// 	department = &DepartmentResponse{
	// 		ID:   user.Department.ID,
	// 		Name: user.Department.Name,
	// 	}
	// }

	return &UserResponse{
		ID:          user.ID,
		Username:    user.Username,
		Email:       user.Email,
		Name:        user.Name,
		Phone:       user.Phone,
		IsActive:    user.IsActive,
		IsSuperUser: user.IsSuperUser,
		// Roles:       roles,
		// Department:  department,
	}
}
