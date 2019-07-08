package main

import (
	"AutoTrading/api"
	"log"
)

func main() {
	log.Println(api.OtcbtcTrades("btcusdt", "1").Get(0).ToString())
	log.Println(api.BianTrades("BTCUSDT", "1").Get(0).ToString())
}
