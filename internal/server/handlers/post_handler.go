package handlers

import (
	"context"
	"go-echo-starter/internal/domain"
	"go-echo-starter/internal/models"
	"go-echo-starter/internal/requests"
	"go-echo-starter/internal/responses"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type postService interface {
	Create(ctx context.Context, post *models.Post) error
	GetPosts(ctx context.Context) ([]models.Post, error)
	GetPost(ctx context.Context, id uint) (models.Post, error)
	UpdateByUser(ctx context.Context, request domain.UpdatePostRequest) (*models.Post, error)
	DeleteByUser(ctx context.Context, request domain.DeletePostRequest) error
}

type PostHandlers struct {
	postService postService
}

func NewPostHandlers(postService postService) *PostHandlers {
	return &PostHandlers{postService: postService}
}

// CreatePost godoc
//
//	@Summary		Create a new post
//	@Description	Create a post for authenticated user
//	@ID				create-post
//	@Tags			Posts
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			request	body		requests.CreatePostRequest	true	"Post data"
//	@Success		201		{object}	responses.Data
//	@Failure		400		{object}	responses.Error
//	@Failure		401		{object}	responses.Error
//	@Router			/posts [post]
func (p *PostHandlers) CreatePost(c echo.Context) error {
	authClaims, err := getAuthClaims(c)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
	}

	var createPostRequest requests.CreatePostRequest
	if err := c.Bind(&createPostRequest); err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "Failed to bind request: "+err.Error())
	}

	if err := createPostRequest.Validate(); err != nil {
		return responses.ValidationErrorResponse(
			c,
			http.StatusBadRequest,
			"Validation failed",
			responses.ParseValidationErrors(err),
		)
	}

	post := &models.Post{
		Title:   createPostRequest.Title,
		Content: createPostRequest.Content,
		UserID:  authClaims.ID,
	}

	if err := p.postService.Create(c.Request().Context(), post); err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "Failed to create post: "+err.Error())
	}

	return responses.MessageResponse(c, http.StatusCreated, "Post successfully created")
}

// GetPosts godoc
//
//	@Summary		Get all posts
//	@Description	Retrieve all posts from database
//	@ID				get-posts
//	@Tags			Posts
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200		{array}		responses.Data
//	@Failure		404		{object}	responses.Error
//	@Router			/posts [get]
func (p *PostHandlers) GetPosts(c echo.Context) error {
	posts, err := p.postService.GetPosts(c.Request().Context())
	if err != nil {
		return responses.ErrorResponse(c, http.StatusNotFound, "Failed to get all posts: "+err.Error())
	}

	response := responses.NewPostResponse(posts)
	return responses.Response(c, http.StatusOK, response)
}

// UpdatePost godoc
//
//	@Summary		Update a post
//	@Description	Update a post owned by authenticated user
//	@ID				update-post
//	@Tags			Posts
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//
//	@Param			id		path		int						true	"Post ID"
//	@Param			request	body		requests.UpdatePostRequest	true	"Update post payload"
//
//	@Success		200		{object}	responses.Data
//	@Failure		400		{object}	responses.Error
//	@Failure		401		{object}	responses.Error
//	@Failure		500		{object}	responses.Error
//	@Router			/posts/{id} [put]
func (p *PostHandlers) UpdatePost(c echo.Context) error {
	authClaims, err := getAuthClaims(c)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
	}

	idParam := c.Param("id")
	postID, err := strconv.Atoi(idParam)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "Invalid post ID")
	}

	var updatePostRequest requests.UpdatePostRequest
	if err := c.Bind(&updatePostRequest); err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "Failed to bind request: "+err.Error())
	}

	if err := updatePostRequest.Validate(); err != nil {
		return responses.ValidationErrorResponse(
			c,
			http.StatusBadRequest,
			"Validation failed",
			responses.ParseValidationErrors(err),
		)
	}

	data := domain.UpdatePostRequest{
		Title:   updatePostRequest.Title,
		Content: updatePostRequest.Content,
		UserID:  authClaims.ID,
		PostID:  uint(postID),
	}

	if _, err := p.postService.UpdateByUser(c.Request().Context(), data); err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "Failed to update post: "+err.Error())
	}

	return responses.MessageResponse(c, http.StatusCreated, "Post successfully updated")
}

// DeletePost godoc
//
//	@Summary		Delete a post
//	@Description	Delete a post owned by authenticated user
//	@ID				delete-post
//	@Tags			Posts
//	@Produce		json
//	@Security		BearerAuth
//
//	@Param			id	path		int	true	"Post ID"
//
//	@Success		200	{object}	responses.Data
//	@Failure		400	{object}	responses.Error
//	@Failure		401	{object}	responses.Error
//	@Failure		500	{object}	responses.Error
//	@Router			/posts/{id} [delete]
func (p *PostHandlers) DeletePost(c echo.Context) error {
	authClaims, err := getAuthClaims(c)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
	}

	idParam := c.Param("id")
	postID, err := strconv.Atoi(idParam)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "Invalid post ID")
	}

	data := domain.DeletePostRequest{
		UserID: authClaims.ID,
		PostID: uint(postID),
	}

	if err := p.postService.DeleteByUser(c.Request().Context(), data); err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "Failed to delete post: "+err.Error())
	}

	return responses.MessageResponse(c, http.StatusCreated, "Post successfully deleted")
}
