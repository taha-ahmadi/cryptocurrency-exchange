package main

import (
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/labstack/echo/v4"
	"github.com/taha-ahmadi/cryptocurrency-exchange/api"
	"log"
)

const (
	exchangePrivateKey = "4f3edf983ac636a65a842ce7c78d9aa706d3b113bce9c46f30d7d21715b23b1d" // We must add it to ENV file
)

func main() {
	e := echo.New()

	// ETH client setup
	client, err := ethclient.Dial("http://localhost:8545")
	if err != nil {
		log.Fatalln(err)
	}

	ex, err := api.NewExchange(exchangePrivateKey, client)
	if err != nil {
		log.Fatalln("cannot make Exchange instance", err)
	}

	user1, _ := api.NewUser("6cbed15c793ce57650b9877cf6fa156fbef513c4e6134f022a85b1ffdd59b2a1", 1)
	user2, _ := api.NewUser("6370fd033278c143179d81c5526140625662b8daa446c22ee2d73db3707e620c", 2)

	ex.Users[user1.ID] = user1
	ex.Users[user2.ID] = user2

	e.GET("/books/:market", ex.HandleGetMarket)
	e.POST("/orders", ex.HandlePlaceOrder)
	e.DELETE("/orders/:id", ex.CancelOrder)

	e.Start(":3000")
}
