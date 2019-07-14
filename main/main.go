package main

import (
	"AutoTrading/api"
	"log"
	"strconv"
)

func main() {
	otcbtcPrice := api.OtcbtcTrades("btcusdt", "1").Get(0, "price").ToFloat64()
	log.Println("otcbtc cur btc price: " + strconv.FormatFloat(otcbtcPrice, 'f', -1, 64))
	bianPrice := api.BianTrades("BTCUSDT", "1").Get(0, "price").ToFloat64()
	log.Println("bian cur btc price: " + strconv.FormatFloat(bianPrice, 'f', -1, 64))
}
