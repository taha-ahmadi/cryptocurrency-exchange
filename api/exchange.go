package api

import (
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/labstack/echo/v4"
	"github.com/taha-ahmadi/cryptocurrency-exchange/matchingengine"
	"math/big"
	"net/http"
	"strconv"
	"strings"
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
	Users      map[uint64]*User  // UserId to User
	orders     map[uint64]uint64 // OrderId to UserID
	privateKey *ecdsa.PrivateKey // Exchange's hot wallet privateKey for transfer coins to other user
	Client     *ethclient.Client // for future, we can use another interface and hashmap to support multiple clients
	orderbooks map[Market]*matchingengine.Orderbook
}

func NewExchange(privateKey string, client *ethclient.Client) (*Exchange, error) {
	orderbooks := make(map[Market]*matchingengine.Orderbook)
	orderbooks[MarketETH] = matchingengine.NewOrderbook()
	orderbooks[MarketBTC] = matchingengine.NewOrderbook()

	pk, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		return nil, err
	}

	return &Exchange{
		Users:      make(map[uint64]*User),
		orders:     make(map[uint64]uint64),
		privateKey: pk,
		Client:     client,
		orderbooks: orderbooks,
	}, nil
}

type User struct {
	ID         uint64
	PrivateKey *ecdsa.PrivateKey
}

func NewUser(privKey string, userId uint64) (*User, error) {
	pk, err := crypto.HexToECDSA(privKey)
	if err != nil {
		return nil, err
	}
	return &User{
		ID:         userId,
		PrivateKey: pk,
	}, nil
}

// PlaceOrderRequest is a data that somebody is gonna sent over API
type PlaceOrderRequest struct {
	UserID uint64
	Type   OrderType // Limit or Market
	IsBid  bool
	Amount float64
	Price  float64
	Market Market
}

type PlaceLimitOrderResponse struct {
	OrderID uint64
}

type Order struct {
	UserID    uint64
	ID        uint64
	Amount    float64
	IsBid     bool
	Price     float64
	Timestamp int64
}

type OrderbookResponse struct {
	TotalAsksVolume float64
	TotalBidsVolume float64

	Asks []*Order
	Bids []*Order
}

type MatchedOrder struct {
	Price        float64
	AmountFilled float64
	ID           uint64
}

func (ex *Exchange) handleMarketPlaceOrder(market Market, order *matchingengine.Order) ([]*MatchedOrder, matchingengine.Matches) {
	ob := ex.orderbooks[market]
	matches := ob.PlaceMarketOrder(order) // we use matches for transfer coins
	matchedOrders := make([]*MatchedOrder, len(matches))

	isBid := false
	if order.Bid {
		isBid = true
	}

	for i := 0; i < len(matchedOrders); i++ {
		id := matches[i].Bid.ID
		if isBid {
			id = matches[i].Ask.ID
		}

		matchedOrders[i] = &MatchedOrder{
			ID:           id,
			Price:        matches[i].Price,
			AmountFilled: matches[i].AmountFilled,
		}
	}

	return matchedOrders, matches
}
func (ex *Exchange) handleLimitPlaceOrder(market Market, price float64, order *matchingengine.Order) error {
	ob := ex.orderbooks[market]
	ob.PlaceLimitOrder(price, order) // A limit order is a bucket that holds different orders that are setting
	// at the same price level but with different amount from different people

	return nil
}
func (ex *Exchange) HandlePlaceOrder(c echo.Context) error {
	var placeOrderData PlaceOrderRequest
	if err := json.NewDecoder(c.Request().Body).Decode(&placeOrderData); err != nil {
		return err
	}

	market := placeOrderData.Market
	order := matchingengine.NewOrder(placeOrderData.IsBid, placeOrderData.Amount, placeOrderData.UserID)

	// handle Market Order
	if strings.ToTitle(string(placeOrderData.Type)) == string(MarketOrder) {
		matchedOrders, matches := ex.handleMarketPlaceOrder(market, order)
		if err := ex.handleMatches(matches); err != nil {
			return err
		}
		return c.JSON(http.StatusCreated, matchedOrders)
	}

	// handle Limit Order
	if strings.ToTitle(string(placeOrderData.Type)) == string(LimitOrder) {
		err := ex.handleLimitPlaceOrder(market, placeOrderData.Price, order)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]any{"msg": "cannot placeOrder!"})
		}

		resp := &PlaceLimitOrderResponse{
			OrderID: order.ID,
		}

		return c.JSON(http.StatusCreated, resp)
	}

	return c.JSON(http.StatusInternalServerError, map[string]any{"msg": "Makert Not Exitst!"})
}

func (ex *Exchange) HandleGetMarket(c echo.Context) error {
	market := Market(c.Param("market"))

	ob, ok := ex.orderbooks[market]

	if !ok {
		return c.JSON(http.StatusBadRequest, map[string]any{"msg": "market not found!"})
	}

	var orderbookResponse = OrderbookResponse{
		TotalAsksVolume: ob.AskTotalVolume(),
		TotalBidsVolume: ob.BidTotalVolume(),
	}

	// make asks response
	for _, limit := range ob.Asks() {
		for _, order := range limit.Orders {
			o := Order{
				UserID:    order.UserID,
				ID:        order.ID,
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
				UserID:    order.UserID,
				ID:        order.ID,
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
	id, _ := strconv.Atoi(c.Param("id"))
	ob := ex.orderbooks[MarketETH]
	ob.CancelOrder(ob.Orders[uint64(id)])

	return c.JSON(http.StatusOK, map[string]any{"msg": "Canceled"})
}

func (ex *Exchange) handleMatches(matches matchingengine.Matches) error {
	for _, match := range matches {
		fromUser, ok := ex.Users[match.Ask.UserID]
		if !ok {
			return fmt.Errorf("user not found: %d", match.Ask.UserID)
		}

		toUser, ok := ex.Users[match.Bid.UserID]
		if !ok {
			return fmt.Errorf("user not found: %d", match.Bid.UserID)
		}
		toAddresss := crypto.PubkeyToAddress(toUser.PrivateKey.PublicKey)

		amount := big.NewInt(int64(match.AmountFilled))
		transferETH(ex.Client, fromUser.PrivateKey, toAddresss, amount)
	}

	return nil
}
