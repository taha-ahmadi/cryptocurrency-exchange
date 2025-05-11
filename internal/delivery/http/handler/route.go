package handler

import (
	"github.com/labstack/echo/v4"
)

// RegisterRoutes registers all routes with the Echo server
func (h *Handler) RegisterRoutes(e *echo.Echo) {
	e.GET("/books/:market", h.HandleGetMarket)
	e.GET("/books/:market/best/bid", h.HandleGetBestBidLimit)
	e.GET("/books/:market/best/ask", h.HandleGetBestAskLimit)
	e.GET("/orders/:userID", h.HandleGetOrder)
	e.POST("/orders", h.HandlePlaceOrder)
	e.GET("/trades/:market", h.HandleGetTrades)
	e.DELETE("/orders/:id", h.HandleCancelOrder)
}
