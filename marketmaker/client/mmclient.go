package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/taha-ahmadi/cryptocurrency-exchange/api"
	"net/http"
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

type PlaceLimitOrderParams struct {
	UserID uint64
	Bid    bool
	price  float64
	amount float64
}

func (c *MMClient) PlaceLimitOrder(params *PlaceLimitOrderParams) (*api.PlaceLimitOrderResponse, error) {
	endpoint := Endpoint + "/order"
	data := api.PlaceOrderRequest{
		UserID: params.UserID,
		Type:   api.LimitOrder,
		IsBid:  params.Bid,
		Amount: params.amount,
		Price:  params.price,
		Market: api.MarketETH,
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

	var result = &api.PlaceLimitOrderResponse{}
	if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
		return nil, err
	}
	return result, nil
}

func (c *MMClient) CancelOrder(orderID uint64) error {
	endpoint := fmt.Sprintf("%s/order/%d", Endpoint, orderID)
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
