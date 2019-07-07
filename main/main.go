package main

import (
	"AutoTrading/api"
	"log"
)

func main() {
	log.Println(api.OtcbtcDepth("btcusdt", "1"))
	log.Println(api.BianDepth("BTCUSDT", "5"))
}
