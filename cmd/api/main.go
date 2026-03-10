package main

import (
	"context"
	"errors"
	"fmt"
	"go-echo-starter/internal/config"
	"go-echo-starter/internal/db"
	"go-echo-starter/internal/repositories"
	"go-echo-starter/internal/server"
	"go-echo-starter/internal/server/handlers"
	"go-echo-starter/internal/server/middleware"
	"go-echo-starter/internal/server/routes"
	"go-echo-starter/internal/services/auth"
	"go-echo-starter/internal/services/oauth"
	"go-echo-starter/internal/services/post"
	"go-echo-starter/internal/services/token"
	"go-echo-starter/internal/services/user"
	"go-echo-starter/internal/slogx"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/swaggo/swag/example/override/docs"
)

const shutdownTimeout = 20 * time.Second

func main() {
	if err := run(); err != nil {
		slog.Error("Service run error", "err", err.Error())
		os.Exit(1)
	}
}

func run() error {
	if err := godotenv.Load(); err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("load env file: %w", err)
	}

	var cfg config.Config
	if err := env.Parse(&cfg); err != nil {
		return fmt.Errorf("parse env: %w", err)
	}

	docs.SwaggerInfo.Host = fmt.Sprintf("%s:%s", cfg.HTTP.Host, cfg.HTTP.Port)

	if err := slogx.Init(cfg.Logger); err != nil {
		return fmt.Errorf("init logger: %w", err)
	}

	gormDB, err := db.NewGormDB(cfg.DB)
	if err != nil {
		return fmt.Errorf("new db connection: %w", err)
	}

	userRepository := repositories.NewUserRepository(gormDB)
	userService := user.NewService(userRepository)

	postRepository := repositories.NewPostRepository(gormDB)
	postService := post.NewService(postRepository)

	provider, err := oidc.NewProvider(context.Background(), "https://accounts.google.com")
	if err != nil {
		return fmt.Errorf("oidc.NewProvider: %w", err)
	}

	verifier := provider.Verifier(&oidc.Config{ClientID: cfg.OAuth.ClientID})

	tokenService := token.NewService(
		time.Now,
		cfg.Auth.AccessTokenDuration,
		cfg.Auth.RefreshTokenDuration,
		[]byte(cfg.Auth.AccessSecret),
		[]byte(cfg.Auth.RefreshSecret),
	)

	authService := auth.NewService(userService, tokenService)
	oAuthService := oauth.NewService(verifier, tokenService, userService)

	postHandler := handlers.NewPostHandlers(postService)
	authHandler := handlers.NewAuthHandler(authService)
	oAuthHandler := handlers.NewOAuthHandler(oAuthService)
	registerHandler := handlers.NewRegisterHandler(userService)

	authMiddleware := middleware.NewAuthMiddleware(cfg.Auth.AccessSecret)
	reguestLoggerMiddleware := middleware.NewRequestLogger(slogx.NewTraceStarter(uuid.NewV7))
	requestDebuggerMiddleware := middleware.NewRequestDebugger()

	engine := routes.ConfigureRoutes(routes.Handlers{
		PostHandler:               postHandler,
		AuthHandler:               authHandler,
		OAuthHandler:              oAuthHandler,
		RegisterHandler:           registerHandler,
		AuthMiddleware:            authMiddleware,
		RequestLoggerMiddleware:   reguestLoggerMiddleware,
		RequestDebuggerMiddleware: requestDebuggerMiddleware,
	})
	if err != nil {
		return fmt.Errorf("configure routes: %w", err)
	}

	app := server.NewServer(engine)
	go func() {
		if err = app.Start(cfg.HTTP.Port); err != nil {
			slog.Error("Server error", "err", err.Error())
		}
	}()

	shutdownChannel := make(chan os.Signal, 1)
	signal.Notify(shutdownChannel, os.Interrupt, syscall.SIGHUP, syscall.SIGTERM)
	<-shutdownChannel

	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := app.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("http server shutdown: %w", err)
	}

	dbConnection, err := gormDB.DB()
	if err != nil {
		return fmt.Errorf("get db connection: %w", err)
	}

	if err := dbConnection.Close(); err != nil {
		return fmt.Errorf("close db connection: %w", err)
	}

	return nil
}

// func main() {
// 	// Load configuration
// 	cfg := config.LoadConfig()

// 	// Initialize database
// 	db := database.NewMySQLConnection(cfg)
// 	defer db.Close()

// 	// Initialize repositories
// 	userRepo := repositories.NewUserRepository(db)

// 	// Initialize handlers
// 	userHandler := handlers.NewUserHandler(userRepo)

// 	// Initialize Echo
// 	e := echo.New()

// 	// Middleware
// 	e.Use(middleware.Logger())
// 	e.Use(middleware.Recover())
// 	e.Use(middleware.CORS())

// 	// Custom validator
// 	e.Validator = &CustomValidator{validator: validator.New()}

// 	// Routes
// 	api := e.Group("/api/v1")

// 	// User routes
// 	api.POST("/users", userHandler.CreateUser)
// 	api.GET("/users", userHandler.GetAllUsers)
// 	api.GET("/users/:id", userHandler.GetUser)
// 	api.PUT("/users/:id", userHandler.UpdateUser)
// 	api.DELETE("/users/:id", userHandler.DeleteUser)

// 	// Health check
// 	e.GET("/health", func(c echo.Context) error {
// 		return c.JSON(200, map[string]string{"status": "ok"})
// 	})

// 	// Start server
// 	log.Printf("Server starting on port %s", cfg.ServerPort)
// 	e.Logger.Fatal(e.Start(":" + cfg.ServerPort))
// }
