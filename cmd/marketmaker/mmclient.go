package marketmaker

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/taha-ahmadi/cryptocurrency-exchange/internal/exchanges"
	"github.com/taha-ahmadi/cryptocurrency-exchange/internal/matchingengine"
)

const Endpoint = "http://localhost:3000"

// MMClient (Market Making Client) For interact with market making algorithm we need a client API to our own exchange to place limit/market orders,
// retrieve and cancel them
type MMClient struct {
	*http.Client
}

func NewClinet() *MMClient {
	return &MMClient{
		Client: http.DefaultClient,
	}
}

func (c *MMClient) GetTrades(market string) ([]*matchingengine.Trade, error) {
	e := fmt.Sprintf("%s/trades/%s", Endpoint, market)
	req, err := http.NewRequest(http.MethodGet, e, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	trades := []*matchingengine.Trade{}
	if err := json.NewDecoder(resp.Body).Decode(&trades); err != nil {
		return nil, err
	}

	return trades, nil
}

func (c *MMClient) GetOrders(userID int64) (*exchanges.GetOrdersResponse, error) {
	e := fmt.Sprintf("%s/orders/%d", Endpoint, userID)

	req, err := http.NewRequest(http.MethodGet, e, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	orders := exchanges.GetOrdersResponse{}
	if err := json.NewDecoder(resp.Body).Decode(&orders); err != nil {
		return nil, err
	}

	return &orders, nil
}

type PlaceOrderParams struct {
	UserID uint64
	Bid    bool
	Price  float64
	Amount float64
}

func (c *MMClient) PlaceLimitOrder(p *PlaceOrderParams) (*exchanges.PlaceOrderResponse, error) {
	params := &exchanges.PlaceOrderRequest{
		UserID: p.UserID,
		Type:   exchanges.LimitOrder,
		IsBid:  p.Bid,
		Amount: p.Amount,
		Price:  p.Price,
		Market: exchanges.MarketETH,
	}

	body, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	e := Endpoint + "/orders"

	req, err := http.NewRequest(http.MethodPost, e, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	placeOrderResponse := &exchanges.PlaceOrderResponse{}
	if err := json.NewDecoder(resp.Body).Decode(placeOrderResponse); err != nil {
		return nil, err
	}

	return placeOrderResponse, nil
}
func (c *MMClient) PlaceMarketOrder(params *PlaceOrderParams) (*exchanges.PlaceOrderResponse, error) {
	endpoint := Endpoint + "/orders"
	data := exchanges.PlaceOrderRequest{
		UserID: params.UserID,
		Type:   exchanges.MarketOrder,
		IsBid:  params.Bid,
		Amount: params.Amount,
		Market: exchanges.MarketETH, // Hard coded ETH because we just support ETH for now
	}
	body, err := json.Marshal(data)
	req, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	var result = &exchanges.PlaceOrderResponse{}
	if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
		return nil, err
	}
	return result, nil
}

func (c *MMClient) CancelOrder(orderID uint64) error {
	endpoint := fmt.Sprintf("%s/orders/%d", Endpoint, orderID)
	req, err := http.NewRequest(http.MethodDelete, endpoint, nil)
	if err != nil {
		return err
	}

	resp, err := c.Do(req)
	if err != nil {
		return err
	}

	_ = resp

	return nil
}

func (c *MMClient) GetBestBid() (float64, error) {
	endpoint := fmt.Sprintf("%s/books/ETH/best/bid", Endpoint)
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return 0.0, err
	}

	resp, err := c.Do(req)
	if err != nil {
		return 0.0, err
	}

	priceResp := &exchanges.PriceResponse{}
	if err := json.NewDecoder(resp.Body).Decode(priceResp); err != nil {
		return 0.0, err
	}

	return priceResp.Price, nil
}

func (c *MMClient) GetBestAsk() (float64, error) {
	e := fmt.Sprintf("%s/books/ETH/best/ask", Endpoint)
	req, err := http.NewRequest(http.MethodGet, e, nil)
	if err != nil {
		return 0, err
	}

	resp, err := c.Do(req)
	if err != nil {
		return 0, err
	}

	priceResp := &exchanges.PriceResponse{}
	if err := json.NewDecoder(resp.Body).Decode(priceResp); err != nil {
		return 0, err
	}

	return priceResp.Price, err
}
