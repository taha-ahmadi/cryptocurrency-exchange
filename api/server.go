package api

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
	"log"
)

var (
	exchangePrivateKey string
	ETHHost            string
	ETHClient          *ethclient.Client
)

func ServerEngine() {
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	exchangePrivateKey = viper.GetString("ExchangePrivateKey")
	ETHHost = viper.GetString("ETHHost")

	e := echo.New()

	// ETH client setup
	ETHClient, err := ethclient.Dial(ETHHost)
	if err != nil {
		log.Fatalln(err)
	}

	ex, err := NewExchange(exchangePrivateKey, ETHClient)
	if err != nil {
		log.Fatalln("cannot make Exchange instance", err)
	}

	user1, _ := NewUser("0c4678963e0aa2cf580300be0536f69e0b77f7dea52ba9de5f18a739e4c26d3c", 1)
	user2, _ := NewUser("2ab26e4ce65e1d4f36aad9b48f39dde64b263b6134b8998015caa981124bb867", 2)

	ex.Users[user1.ID] = user1
	ex.Users[user2.ID] = user2

	user1Add := "0xbEA21206bEFd190A1198eaf283866Ea831a5704d"
	user1Balance, err := ETHClient.BalanceAt(context.Background(), common.HexToAddress(user1Add), nil)
	if err != nil {
		log.Fatalln(err)
	}

	user2Add := "0x46445675277F7153D8F4658fe13F2C2cEf1625Ad"
	user2Balance, err := ETHClient.BalanceAt(context.Background(), common.HexToAddress(user2Add), nil)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("user1", user1Balance)
	fmt.Println("user2", user2Balance)

	e.GET("/books/:market", ex.HandleGetMarket)
	e.GET("/books/:market/best/bid", ex.HandleGetBestBidLimit)
	e.GET("/books/:market/best/ask", ex.HandleGetBestAskLimit)
	e.GET("/orders/:userID", ex.HandleGetOrder)
	e.POST("/orders", ex.HandlePlaceOrder)
	e.GET("/trades/:market", ex.HandleGetTrades)
	e.DELETE("/orders/:id", ex.CancelOrder)

	e.Start(":3000")
}
