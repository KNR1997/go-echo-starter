package routes

import (
	"go-echo-starter/internal/server/handlers"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Handlers struct {
	PostHandler     *handlers.PostHandlers
	AuthHandler     *handlers.AuthHandler
	OAuthHandler    *handlers.OAuthHandler
	RegisterHandler *handlers.RegisterHandler

	AuthMiddleware            echo.MiddlewareFunc
	RequestLoggerMiddleware   echo.MiddlewareFunc
	RequestDebuggerMiddleware echo.MiddlewareFunc
}

func ConfigureRoutes(handlers Handlers) *echo.Echo {
	engine := echo.New()

	// Technical API route initialization.
	//
	// These endpoints exist solely to keep the service running and must not include any
	// business or processing logic.
	// engine.GET("/swagger/*", echoSwagger.WrapHandler)
	engine.GET("/health", func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	})

	api := engine.Group("", handlers.RequestLoggerMiddleware)

	// Private API routes initialization.
	//
	// These endpoints are used primarily for authentication/authorization and may carry sensitive data.
	// Do NOT log request or response bodies; doing so could expose client information.
	privateAPI := api.Group("")

	privateAPI.POST("/login", handlers.AuthHandler.Login)
	privateAPI.POST("/register", handlers.RegisterHandler.Register)

	// Authorized API route initialization.
	//
	// These endpoints implement the core application logic and require authentication
	// before they can be accessed.
	authorizedAPI := api.Group("", handlers.RequestDebuggerMiddleware, handlers.AuthMiddleware)

	authorizedAPI.POST("/posts", handlers.PostHandler.CreatePost)
	authorizedAPI.GET("/posts", handlers.PostHandler.GetPosts)

	return engine

}
