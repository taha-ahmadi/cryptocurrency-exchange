package api

import (
	"encoding/json"
	"github.com/labstack/echo/v4"
	"github.com/taha-ahmadi/cryptocurrency-exchange/matchingengine"
	"net/http"
	"time"
)

type Market string
type OrderType string

const (
	MarketETH Market = "ETH"
	MarketBTC Market = "BTC"
)
const (
	MarketOrder OrderType = "MARKET"
	LimitOrder  OrderType = "LIMIT" // A limit order is a way to provide liquidity to exchange
)

type Exchange struct {
	orderbooks map[Market]*matchingengine.Orderbook
}

func NewExchange() *Exchange {
	orderbooks := make(map[Market]*matchingengine.Orderbook)
	orderbooks[MarketETH] = matchingengine.NewOrderbook()
	orderbooks[MarketBTC] = matchingengine.NewOrderbook()

	return &Exchange{
		orderbooks: orderbooks,
	}
}

// PlaceOrderRequest is a data that somebody is gonna sent over API
type PlaceOrderRequest struct {
	Type   OrderType // Limit or Market
	IsBid  bool
	Amount float64
	Price  float64
	Market Market
}

type Order struct {
	Amount    float64
	IsBid     bool
	Price     float64
	Timestamp int64
}

type OrderbookResponse struct {
	Asks []*Order
	Bids []*Order
}

func (ex *Exchange) HandlePlaceOrder(c echo.Context) error {
	var placeOrderData PlaceOrderRequest
	if err := json.NewDecoder(c.Request().Body).Decode(&placeOrderData); err != nil {
		return err
	}

	market := placeOrderData.Market
	order := matchingengine.NewOrder(placeOrderData.IsBid, placeOrderData.Amount)
	ob := ex.orderbooks[market]

	if string(placeOrderData.Type) == string(MarketOrder) {
		ob.PlaceMarketOrder(order) // we use matches to find who should get its coin
		return c.JSON(http.StatusCreated, map[string]any{"msg": "order submitted successfully"})
	}

	if string(placeOrderData.Type) == string(LimitOrder) {
		ob.PlaceLimitOrder(placeOrderData.Price, order) // A limit order is a bucket that holds different orders that are setting
		// at the same price level but with different amount from different people
		return c.JSON(http.StatusCreated, map[string]any{"msg": "order submitted successfully"})
	}

	return c.JSON(http.StatusInternalServerError, map[string]any{"msg": "Internal Error!"})
}

func (ex *Exchange) HandleGetMarket(c echo.Context) error {
	market := Market(c.Param("market"))

	ob, ok := ex.orderbooks[market]

	if !ok {
		return c.JSON(http.StatusBadRequest, map[string]any{"msg": "market not found!"})
	}

	var orderbookResponse OrderbookResponse

	// make asks response
	for _, limit := range ob.Asks() {
		for _, order := range limit.Orders {
			o := Order{
				Price:     limit.Price,
				Amount:    order.Amount,
				IsBid:     order.Bid,
				Timestamp: time.Now().UnixNano(),
			}
			orderbookResponse.Asks = append(orderbookResponse.Asks, &o)
		}
	}

	// make bids response
	for _, limit := range ob.Bids() {
		for _, order := range limit.Orders {
			o := Order{
				Price:     limit.Price,
				Amount:    order.Amount,
				IsBid:     order.Bid,
				Timestamp: time.Now().UnixNano(),
			}
			orderbookResponse.Bids = append(orderbookResponse.Bids, &o)
		}
	}

	return c.JSON(200, orderbookResponse)
}

// CancelOrder is the most important API because of market making bot that need to cancel orders
// And manage our liquidises quickly
func (ex *Exchange) CancelOrder(c echo.Context) error {
	// TODO: Implement!
	return nil
}
