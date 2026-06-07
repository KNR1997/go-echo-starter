package responses

import "go-echo-starter/internal/models"

type ApiResponse struct {
	ID      uint   `json:"id" example:"1"`
	Path    string `json:"path" example:"/roles/get"`
	Method  string `json:"method" example:"GET"`
	Summary string `json:"summary" example:"summary"`
	Tags    string `json:"tags" example:"roles"`
}

func NewApiResponse(apis []models.Api) *[]ApiResponse {
	apiResponse := make([]ApiResponse, 0)

	for i := range apis {
		apiResponse = append(apiResponse, ApiResponse{
			ID:      apis[i].ID,
			Path:    apis[i].Path,
			Method:  apis[i].Method,
			Summary: apis[i].Summary,
			Tags:    apis[i].Tags,
		})
	}

	return &apiResponse
}
