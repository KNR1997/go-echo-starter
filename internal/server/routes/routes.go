package routes

import (
	"go-echo-starter/internal/server/handlers"
	"net/http"

	echomiddleware "github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"

	"github.com/labstack/echo/v4"
)

type Handlers struct {
	PostHandler        *handlers.PostHandlers
	AuthHandler        *handlers.AuthHandler
	OAuthHandler       *handlers.OAuthHandler
	RegisterHandler    *handlers.RegisterHandler
	UserHandlers       *handlers.UserHandlers
	RoleHandlers       *handlers.RoleHandlers
	DepartmentHandlers *handlers.DepartmentHandlers
	ApiHandlers        *handlers.ApiHandlers
	MenuHandlers       *handlers.MenuHandlers

	AuthMiddleware            echo.MiddlewareFunc
	RequestLoggerMiddleware   echo.MiddlewareFunc
	RequestDebuggerMiddleware echo.MiddlewareFunc
}

func ConfigureRoutes(handlers Handlers) *echo.Echo {
	engine := echo.New()

	engine.Use(echomiddleware.CORSWithConfig(echomiddleware.CORSConfig{
		AllowOrigins: []string{
			"http://localhost:5173",
		},
		AllowMethods: []string{
			echo.GET,
			echo.POST,
			echo.PUT,
			echo.DELETE,
			echo.PATCH,
			echo.OPTIONS,
		},
		AllowHeaders: []string{
			echo.HeaderOrigin,
			echo.HeaderContentType,
			echo.HeaderAccept,
			echo.HeaderAuthorization,
		},
		AllowCredentials: true,
	}))

	// Technical API route initialization.
	//
	// These endpoints exist solely to keep the service running and must not include any
	// business or processing logic.
	engine.GET("/swagger/*", echoSwagger.WrapHandler)
	engine.GET("/health", func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	})

	api := engine.Group("")

	// Private API routes initialization.
	//
	// These endpoints are used primarily for authentication/authorization and may carry sensitive data.
	// Do NOT log request or response bodies; doing so could expose client information.
	privateAPI := api.Group("", handlers.RequestDebuggerMiddleware)

	privateAPI.POST("/login", handlers.AuthHandler.Login)
	privateAPI.POST("/register", handlers.RegisterHandler.Register)

	// Authorized API route initialization.
	//
	// These endpoints implement the core application logic and require authentication
	// before they can be accessed.
	authorizedAPI := api.Group(
		"",
		handlers.AuthMiddleware,
		handlers.RequestDebuggerMiddleware,
	)

	privateAPI.GET("/posts", handlers.PostHandler.GetPostPaginated)
	authorizedAPI.POST("/posts", handlers.PostHandler.CreatePost)
	authorizedAPI.PUT("/posts/:id", handlers.PostHandler.UpdatePost)
	authorizedAPI.DELETE("/posts/:id", handlers.PostHandler.DeletePost)

	privateAPI.GET("/users", handlers.UserHandlers.GetUsers)

	privateAPI.GET("/roles", handlers.RoleHandlers.GetRolePaginated)
	authorizedAPI.POST("/roles", handlers.RoleHandlers.CreateRole)
	authorizedAPI.PUT("/roles/:id", handlers.RoleHandlers.UpdateRole)
	authorizedAPI.DELETE("/roles/:id", handlers.RoleHandlers.DeleteRole)

	privateAPI.GET("/departments", handlers.DepartmentHandlers.GetDepartmentPaginated)
	authorizedAPI.POST("/departments", handlers.DepartmentHandlers.CreateDepartment)
	authorizedAPI.PUT("/departments/:id", handlers.DepartmentHandlers.UpdateDepartment)
	authorizedAPI.DELETE("/departments/:id", handlers.DepartmentHandlers.DeleteDepartment)

	privateAPI.GET("/apis", handlers.ApiHandlers.GetApiPaginated)
	authorizedAPI.POST("/apis", handlers.ApiHandlers.CreateApi)
	authorizedAPI.PUT("/apis/:id", handlers.ApiHandlers.UpdateApi)
	authorizedAPI.DELETE("/apis/:id", handlers.ApiHandlers.DeleteApi)

	privateAPI.GET("/menus", handlers.MenuHandlers.GetMenuPaginated)
	authorizedAPI.POST("/menus", handlers.MenuHandlers.CreateMenu)
	authorizedAPI.PUT("/menus/:id", handlers.MenuHandlers.UpdateMenu)
	authorizedAPI.DELETE("/menus/:id", handlers.MenuHandlers.DeleteMenu)

	return engine

}
