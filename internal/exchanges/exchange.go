package exchanges

import (
	"crypto/ecdsa"
	"errors"
	"fmt"
	"log"
	"math/big"
	"strings"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/taha-ahmadi/cryptocurrency-exchange/internal/eth"
	"github.com/taha-ahmadi/cryptocurrency-exchange/internal/matchingengine"
	"github.com/taha-ahmadi/cryptocurrency-exchange/internal/models"
)

// Exchange represents the main exchange functionality
type Exchange struct {
	Users      map[uint64]*models.User
	Orders     map[uint64][]*matchingengine.Order
	PrivateKey *ecdsa.PrivateKey
	ETHClient  *eth.Client
	Orderbooks map[Market]*matchingengine.Orderbook
	mu         sync.RWMutex
}

// NewExchange creates a new exchange instance
func NewExchange(privateKey string, ethClient *eth.Client) (*Exchange, error) {
	orderbooks := make(map[Market]*matchingengine.Orderbook)
	orderbooks[MarketETH] = matchingengine.NewOrderbook()
	orderbooks[MarketBTC] = matchingengine.NewOrderbook()

	pk, err := cryptoHexToECDSA(privateKey)
	if err != nil {
		return nil, err
	}

	return &Exchange{
		Users:      make(map[uint64]*models.User),
		Orders:     make(map[uint64][]*matchingengine.Order),
		PrivateKey: pk,
		ETHClient:  ethClient,
		Orderbooks: orderbooks,
	}, nil
}

// AddUser adds a new user to the exchange
func (ex *Exchange) AddUser(user *models.User) {
	ex.Users[user.ID] = user
}

// HandleMarketOrder handles a market order
func (ex *Exchange) HandleMarketOrder(market Market, order *matchingengine.Order) ([]*MatchedOrder, error) {
	ob, exists := ex.Orderbooks[market]
	if !exists {
		return nil, fmt.Errorf("market %s does not exist", market)
	}

	matches := ob.PlaceMarketOrder(order)
	matchedOrders := make([]*MatchedOrder, len(matches))

	isBid := order.Bid

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

	// Update orders map by removing filled orders
	ex.UpdateOrdersAfterMatch()

	// Process the actual transfers
	if err := ex.ProcessMatches(matches); err != nil {
		return matchedOrders, err
	}

	return matchedOrders, nil
}

// HandleLimitOrder handles a limit order
func (ex *Exchange) HandleLimitOrder(market Market, price float64, order *matchingengine.Order) error {
	ob, exists := ex.Orderbooks[market]
	if !exists {
		return fmt.Errorf("market %s does not exist", market)
	}

	ob.PlaceLimitOrder(price, order)

	ex.mu.Lock()
	defer ex.mu.Unlock()
	ex.Orders[order.UserID] = append(ex.Orders[order.UserID], order)

	log.Printf("New LIMIT order => type: [%t] | price [%.2f] | size [%.2f]", order.Bid, order.Limit.Price, order.Amount)

	return nil
}

// UpdateOrdersAfterMatch updates the orders map after matches
func (ex *Exchange) UpdateOrdersAfterMatch() {
	// Create a new map to hold unfilled orders
	newOrderMap := make(map[uint64][]*matchingengine.Order)

	ex.mu.Lock()
	defer ex.mu.Unlock()

	for userID, orders := range ex.Orders {
		for i := 0; i < len(orders); i++ {
			if !orders[i].IsFilled() {
				newOrderMap[userID] = append(newOrderMap[userID], orders[i])
			}
		}
	}

	ex.Orders = newOrderMap
}

// PlaceOrder places a new order
func (ex *Exchange) PlaceOrder(req *PlaceOrderRequest) (interface{}, error) {
	market := req.Market
	order := matchingengine.NewOrder(req.IsBid, req.Amount, req.UserID)

	// Handle market order
	if strings.ToUpper(string(req.Type)) == string(MarketOrder) {
		matchedOrders, err := ex.HandleMarketOrder(market, order)
		if err != nil {
			return nil, err
		}
		return matchedOrders, nil
	}

	// Handle limit order
	if strings.ToUpper(string(req.Type)) == string(LimitOrder) {
		err := ex.HandleLimitOrder(market, req.Price, order)
		if err != nil {
			return nil, err
		}
		return &PlaceOrderResponse{OrderID: order.ID}, nil
	}

	return nil, errors.New("unknown order type")
}

// GetOrderbook gets the orderbook for a market
func (ex *Exchange) GetOrderbook(market Market) (*OrderbookResponse, error) {
	ob, exists := ex.Orderbooks[market]
	if !exists {
		return nil, fmt.Errorf("market %s does not exist", market)
	}

	var orderbookResponse = OrderbookResponse{
		TotalAsksVolume: ob.AskTotalVolume(),
		TotalBidsVolume: ob.BidTotalVolume(),
	}

	// Add asks to response
	for _, limit := range ob.Asks() {
		for _, order := range limit.Orders {
			o := &Order{
				UserID:    order.UserID,
				ID:        order.ID,
				Price:     limit.Price,
				Amount:    order.Amount,
				IsBid:     order.Bid,
				Timestamp: order.Timestamp,
			}
			orderbookResponse.Asks = append(orderbookResponse.Asks, o)
		}
	}

	// Add bids to response
	for _, limit := range ob.Bids() {
		for _, order := range limit.Orders {
			o := &Order{
				UserID:    order.UserID,
				ID:        order.ID,
				Price:     limit.Price,
				Amount:    order.Amount,
				IsBid:     order.Bid,
				Timestamp: order.Timestamp,
			}
			orderbookResponse.Bids = append(orderbookResponse.Bids, o)
		}
	}

	return &orderbookResponse, nil
}

// CancelOrder cancels an order
func (ex *Exchange) CancelOrder(orderID uint64) error {
	found := false

	for _, ob := range ex.Orderbooks {
		order, exists := ob.Orders[orderID]
		if exists {
			ob.CancelOrder(order)
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("order %d not found", orderID)
	}

	return nil
}

// GetBestBidPrice gets the best bid price for a market
func (ex *Exchange) GetBestBidPrice(market Market) (float64, error) {
	ob, exists := ex.Orderbooks[market]
	if !exists {
		return 0, fmt.Errorf("market %s does not exist", market)
	}

	if len(ob.Bids()) == 0 {
		return 0, errors.New("no bids available")
	}

	return ob.Bids()[0].Price, nil
}

// GetBestAskPrice gets the best ask price for a market
func (ex *Exchange) GetBestAskPrice(market Market) (float64, error) {
	ob, exists := ex.Orderbooks[market]
	if !exists {
		return 0, fmt.Errorf("market %s does not exist", market)
	}

	if len(ob.Asks()) == 0 {
		return 0, errors.New("no asks available")
	}

	return ob.Asks()[0].Price, nil
}

// GetTrades gets all trades for a market
func (ex *Exchange) GetTrades(market Market) ([]*matchingengine.Trade, error) {
	ob, exists := ex.Orderbooks[market]
	if !exists {
		return nil, fmt.Errorf("market %s does not exist", market)
	}

	return ob.Trades, nil
}

// GetUserOrders gets all orders for a user
func (ex *Exchange) GetUserOrders(userID uint64) (*GetOrdersResponse, error) {
	ex.mu.RLock()
	defer ex.mu.RUnlock()

	orderbookOrders, exists := ex.Orders[userID]
	if !exists {
		return &GetOrdersResponse{}, nil
	}

	ordersResp := &GetOrdersResponse{
		Asks: []Order{},
		Bids: []Order{},
	}

	for i := 0; i < len(orderbookOrders); i++ {
		if orderbookOrders[i].Limit == nil {
			continue
		}

		order := Order{
			ID:        orderbookOrders[i].ID,
			UserID:    orderbookOrders[i].UserID,
			Price:     orderbookOrders[i].Limit.Price,
			Amount:    orderbookOrders[i].Amount,
			Timestamp: orderbookOrders[i].Timestamp,
			IsBid:     orderbookOrders[i].Bid,
		}

		if order.IsBid {
			ordersResp.Bids = append(ordersResp.Bids, order)
		} else {
			ordersResp.Asks = append(ordersResp.Asks, order)
		}
	}

	return ordersResp, nil
}

// ProcessMatches processes matches by transferring ETH
func (ex *Exchange) ProcessMatches(matches matchingengine.Matches) error {
	for _, match := range matches {
		fromUser, ok := ex.Users[match.Ask.UserID]
		if !ok {
			return fmt.Errorf("user not found: %d", match.Ask.UserID)
		}

		toUser, ok := ex.Users[match.Bid.UserID]
		if !ok {
			return fmt.Errorf("user not found: %d", match.Bid.UserID)
		}

		toAddress := common.HexToAddress(toUser.Address)
		amount := big.NewInt(int64(match.AmountFilled))

		err := ex.ETHClient.TransferETH(fromUser.PrivateKey, toAddress, amount)
		if err != nil {
			return fmt.Errorf("failed to transfer ETH: %w", err)
		}
	}

	return nil
}

// Helper functions
func cryptoHexToECDSA(hexKey string) (*ecdsa.PrivateKey, error) {
	return ecdsa.GenerateKey(nil, nil) // This is not safe for production, just a stub for testing
}
