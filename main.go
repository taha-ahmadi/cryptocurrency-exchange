package main

import (
	"github.com/labstack/echo/v4"
	"github.com/taha-ahmadi/cryptocurrency-exchange/api"
)

func main() {
	e := echo.New()
	ex := api.NewExchange()

	e.GET("/books/:market", ex.HandleGetMarket)
	e.POST("/orders", ex.HandlePlaceOrder)
	e.DELETE("/orders/:id", ex.CancelOrder)

	e.Start(":3000")
}
