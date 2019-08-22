package main

import (
	"AutoTrading/api"
	"log"
	"testing"
)

func TestDepth(t *testing.T) {
	log.Println(api.BianLastPrice("BTCUSDT"))
	//log.Println(api.BianLastAllPrice())
	huobiPrice := api.HuobiLastPrice("btcusdt")
	log.Println(huobiPrice)

	//log.Println(api.BianBestTicker("XRPUSDT"))
	//log.Println(api.BianAllBestTicker())

	//log.Println(api.BianDepth("BTCUSDT", 5))
	//log.Println(api.HuobiDepth("btcusdt", "step0"))

	//log.Println(api.BianKLine("BTCUSDT", models.BianDay, 10, 0, 0))
	//log.Println(api.GetHuobiKLine("btcusdt", "15min", 10))

	//curBianPrice := api.BianTrade("BTCUSDT", 1)[0].Price
	//log.Println("当前币安的 BTC 价格为: $" + curBianPrice)
	//log.Println(api.BianAggTrade("BTCUSDT", 5, 0, 0, 0))
	//curHuobiPrice := api.HuobiTrade("btcusdt", 1).Data[0].Data[0].Price
	//log.Println("当前火币的 BTC 价格为: $" + strconv.FormatFloat(curHuobiPrice, 'f', -1, 64))

	//log.Println(api.BianAvgPrice("BTCUSDT"))

	//log.Println(api.BianTicker24("BTCUSDT"))
	//log.Println(api.BianTicker24All())

	// ************************* Account Test ***********************

	// orderid: 207779114
	//log.Println(api.BianOrderByLimit("XRPUSDT", models.SideSell, models.GTC, 50, 0.5, 0))
	//log.Println(api.BianOpenOrder("XRPUSDT"))
	//log.Println(api.BianOrderDelete("XRPUSDT", "", 223811828))
	//log.Println(api.BianOpenOrder("XRPUSDT"))
	//log.Println(api.BianOrderQuery("XRPUSDT", "", 207779114))
	//queryStrategyLBHS("XRPUSDT", "binance")
	//updateQuantity()
	//updateAccount()
	//updateStrategyLBHS("XRP", "binance")
	//RunLBHS()
}
