package main

import (
	"blockchain-scrap/handler"
	"blockchain-scrap/infra"
	"blockchain-scrap/service"

	"github.com/gin-gonic/gin"
)

func main() {

	db := infra.GetDBInstance()
	if err := infra.AutoMigrate(db); err != nil {
		return
	}

	r := gin.Default()

	coinService := service.NewCoinService()
	coinHandler := handler.NewCoinHandler(coinService)

	r.GET("/coins/:id/:contract-address", coinHandler.GetCoinDetail)
	r.GET("/coins/all", coinHandler.GetAllCoins)
	r.GET("/coins/all/stream", coinHandler.StreamCoins)

	r.Run()
}
