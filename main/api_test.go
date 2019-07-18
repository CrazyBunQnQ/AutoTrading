package main

import (
	"AutoTrading/api"
	"log"
	"testing"
)

func TestDepth(t *testing.T) {
	//log.Println(api.BianDepth("BTCUSDT", 5))
	//log.Println(api.GetMarketDepth("btcusdt", "step0"))

	//log.Println(api.BianKLine("BTCUSDT", models.BianDay, 10, 0, 0))
	//log.Println(api.GetHuobiKLine("btcusdt", "15min", 10))

	log.Println(api.BianTrade("BTCUSDT", 5))
	log.Println(api.HuobiTrade("btcusdt", 5))
}
