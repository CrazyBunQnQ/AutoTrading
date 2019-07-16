package main

import (
	"AutoTrading/api"
	"log"
	"testing"
)

func TestDepth(t *testing.T) {
	log.Println(api.BianDepth("BTCUSDT", 5))
	log.Println(api.GetMarketDepth("btcusdt", "step0"))
}
