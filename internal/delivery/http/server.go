package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/taha-ahmadi/cryptocurrency-exchange/internal/config"
	"github.com/taha-ahmadi/cryptocurrency-exchange/internal/delivery/http/handler"
	"github.com/taha-ahmadi/cryptocurrency-exchange/internal/delivery/http/middleware"
	"github.com/taha-ahmadi/cryptocurrency-exchange/internal/exchanges"
	"github.com/taha-ahmadi/cryptocurrency-exchange/internal/models"
	"github.com/taha-ahmadi/cryptocurrency-exchange/pkg/ethclient"
)

// Server represents the HTTP server for the exchange
type Server struct {
	echo    *echo.Echo
	handler *handler.Handler
	config  *config.Config
}

// New creates a new HTTP server
func New(cfg *config.Config) (*Server, error) {
	// Create Echo instance
	e := echo.New()

	// Set up middleware
	e.Use(echoMiddleware.Recover())
	e.Use(echoMiddleware.CORSWithConfig(echoMiddleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
	}))

	// Use custom middleware
	e.Use(middleware.RequestLogger(middleware.RequestLoggerConfig{
		LogRequestBody:  true,
		LogResponseBody: false,
	}))
	e.Use(middleware.ErrorHandler())

	// Create ETH client
	ethClient, err := ethclient.New(cfg.ETHHost)
	if err != nil {
		return nil, fmt.Errorf("failed to create ETH client: %w", err)
	}

	// Create exchange
	exchange, err := exchanges.New(cfg.ExchangePrivateKey, ethClient)
	if err != nil {
		return nil, fmt.Errorf("failed to create exchange: %w", err)
	}

	// Add test users
	user1, err := models.NewUser("0c4678963e0aa2cf580300be0536f69e0b77f7dea52ba9de5f18a739e4c26d3c", 1)
	if err != nil {
		return nil, fmt.Errorf("failed to create user1: %w", err)
	}

	user2, err := models.NewUser("2ab26e4ce65e1d4f36aad9b48f39dde64b263b6134b8998015caa981124bb867", 2)
	if err != nil {
		return nil, fmt.Errorf("failed to create user2: %w", err)
	}

	exchange.AddUser(user1)
	exchange.AddUser(user2)

	// Create handler
	handler := handler.New(exchange)

	return &Server{
		echo:    e,
		handler: handler,
		config:  cfg,
	}, nil
}

// Start starts the HTTP server
func (s *Server) Start() error {
	// Register routes
	s.handler.RegisterRoutes(s.echo)

	// Get port from config
	port := s.config.ServerPort
	if port == "" {
		port = "3000"
	}
	addr := ":" + port

	// Start server in a goroutine
	go func() {
		if err := s.echo.Start(addr); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	log.Printf("Server started on port %s", port)

	// Wait for interrupt signal to gracefully shut down the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// Shutdown with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := s.echo.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to gracefully shut down server: %w", err)
	}

	return nil
}
