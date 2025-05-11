package api

import (
	"context"
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/taha-ahmadi/cryptocurrency-exchange/internal/config"
)

// Server represents the API server
type Server struct {
	echo      *echo.Echo
	exchange  *Exchange
	config    *config.Config
	ethClient *ethclient.Client
}

// NewServer creates a new API server
func NewServer(cfg *config.Config) (*Server, error) {
	e := echo.New()

	// Setup CORS
	corsConfig := middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:8080"},
		AllowHeaders: []string{echo.HeaderAuthorization},
	}
	e.Use(middleware.CORSWithConfig(corsConfig))

	// ETH client setup
	ethClient, err := ethclient.Dial(cfg.ETHHost)
	if err != nil {
		return nil, fmt.Errorf("cannot dial eth client: %w", err)
	}

	// Create exchange
	ex, err := NewExchange(cfg.ExchangePrivateKey, ethClient)
	if err != nil {
		return nil, fmt.Errorf("cannot make Exchange instance: %w", err)
	}

	// Create test users
	user1, _ := NewUser("0c4678963e0aa2cf580300be0536f69e0b77f7dea52ba9de5f18a739e4c26d3c", 1)
	user2, _ := NewUser("2ab26e4ce65e1d4f36aad9b48f39dde64b263b6134b8998015caa981124bb867", 2)

	ex.Users[user1.ID] = user1
	ex.Users[user2.ID] = user2

	user1Add := "0xbEA21206bEFd190A1198eaf283866Ea831a5704d"
	user1Balance, err := ethClient.BalanceAt(context.Background(), common.HexToAddress(user1Add), nil)
	if err != nil {
		return nil, err
	}

	user2Add := "0x46445675277F7153D8F4658fe13F2C2cEf1625Ad"
	user2Balance, err := ethClient.BalanceAt(context.Background(), common.HexToAddress(user2Add), nil)
	if err != nil {
		return nil, err
	}

	fmt.Println("user1", user1Balance)
	fmt.Println("user2", user2Balance)

	server := &Server{
		echo:      e,
		exchange:  ex,
		config:    cfg,
		ethClient: ethClient,
	}

	return server, nil
}

// SetupRoutes configures all the routes for the API server
func (s *Server) SetupRoutes() {
	s.echo.GET("/books/:market", s.exchange.HandleGetMarket)
	s.echo.GET("/books/:market/best/bid", s.exchange.HandleGetBestBidLimit)
	s.echo.GET("/books/:market/best/ask", s.exchange.HandleGetBestAskLimit)
	s.echo.GET("/orders/:userID", s.exchange.HandleGetOrder)
	s.echo.POST("/orders", s.exchange.HandlePlaceOrder)
	s.echo.GET("/trades/:market", s.exchange.HandleGetTrades)
	s.echo.DELETE("/orders/:id", s.exchange.CancelOrder)
}

// Start starts the server
func (s *Server) Start() error {
	port := s.config.ServerPort
	if port == "" {
		port = "3000"
	}

	return s.echo.Start(":" + port)
}

// ServerEngine creates and starts the API server - for backward compatibility
func ServerEngine() {
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	server, err := NewServer(cfg)
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	server.SetupRoutes()
	if err := server.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
