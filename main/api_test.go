package main

import (
	"AutoTrading/api"
	"AutoTrading/models"
	"log"
	"testing"
)

func TestDepth(t *testing.T) {
	//log.Println(api.BianDepth("BTCUSDT", 5))
	//log.Println(api.GetMarketDepth("btcusdt", "step0"))

	log.Println(api.GetBianKLine("BTCUSDT", models.BianDay, 500, 0, 0))
	log.Println(api.GetHuobiKLine("btcusdt", "15min", 10))
}
