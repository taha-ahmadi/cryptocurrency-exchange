package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/taha-ahmadi/cryptocurrency-exchange/internal/exchanges"
)

// HandleGetMarket handles the GET /books/:market endpoint
func (h *Handler) HandleGetMarket(c echo.Context) error {
	market := exchanges.Market(c.Param("market"))

	orderbook, err := h.Exchange.GetOrderbook(market)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, orderbook)
}

// HandleGetBestBidLimit handles the GET /books/:market/best/bid endpoint
func (h *Handler) HandleGetBestBidLimit(c echo.Context) error {
	market := exchanges.Market(c.Param("market"))

	price, err := h.Exchange.GetBestBidPrice(market)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, exchanges.PriceResponse{Price: price})
}

// HandleGetBestAskLimit handles the GET /books/:market/best/ask endpoint
func (h *Handler) HandleGetBestAskLimit(c echo.Context) error {
	market := exchanges.Market(c.Param("market"))

	price, err := h.Exchange.GetBestAskPrice(market)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, exchanges.PriceResponse{Price: price})
}

// HandleGetTrades handles the GET /trades/:market endpoint
func (h *Handler) HandleGetTrades(c echo.Context) error {
	market := exchanges.Market(c.Param("market"))

	trades, err := h.Exchange.GetTrades(market)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, trades)
}

// HandleGetOrder handles the GET /orders/:userID endpoint
func (h *Handler) HandleGetOrder(c echo.Context) error {
	userIDStr := c.Param("userID")
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "invalid user ID"})
	}

	orders, err := h.Exchange.GetUserOrders(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, orders)
}

// HandlePlaceOrder handles the POST /orders endpoint
func (h *Handler) HandlePlaceOrder(c echo.Context) error {
	var placeOrderData exchanges.PlaceOrderRequest
	if err := json.NewDecoder(c.Request().Body).Decode(&placeOrderData); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "invalid request body"})
	}

	result, err := h.Exchange.PlaceOrder(&placeOrderData)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, result)
}

// HandleCancelOrder handles the DELETE /orders/:id endpoint
func (h *Handler) HandleCancelOrder(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": "invalid order ID"})
	}

	err = h.Exchange.CancelOrder(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"message": "order cancelled successfully"})
}
