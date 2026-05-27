package responses

import "go-echo-starter/internal/models"

type UserResponse struct {
	Email    string `json:"email" example:"admin@demo.com"`
	Name     string `json:"name" example:"John Doe"`
	Phone    string `json:"phone" example:"0113123888"`
	IsActive bool   `json:"is_active" example:"true"`
	ID       uint   `json:"id" example:"1"`
}

func NewUserResponse(users []models.User) *[]UserResponse {
	userResponse := make([]UserResponse, 0)

	for i := range users {
		userResponse = append(userResponse, UserResponse{
			Email:    users[i].Email,
			Name:     users[i].Name,
			Phone:    users[i].Phone,
			IsActive: users[i].IsActive,
			ID:       users[i].ID,
		})
	}

	return &userResponse
}
