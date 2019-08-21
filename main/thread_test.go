package main

import (
	"AutoTrading/api"
	"AutoTrading/utils"
	"log"
	"math"
	"strconv"
	"sync"
	"testing"
	"time"
)

var getPriceThread = sync.WaitGroup{}
var bianPrice, huobiPrice, diffPrice float64

func TestThread(t *testing.T) {
	getPriceThread.Add(1)
	go getBianLastPrice("BTCUSDT")
	getPriceThread.Add(1)
	go getHuobiLastPrice("btcusdt")

	getPriceThread.Wait()
	log.Println(getDiffPrice(bianPrice, huobiPrice))
}

func getBianLastPrice(symbol string) float64 {
	log.Println(strconv.FormatInt(utils.UnixMillis(time.Now()), 10) + " bian")
	defer getPriceThread.Done()
	bianPrice, _ = strconv.ParseFloat(api.BianLastPrice(symbol).Price, 64)
	log.Println(strconv.FormatInt(utils.UnixMillis(time.Now()), 10) + " bian : " + strconv.FormatFloat(bianPrice, 'f', -1, 64))
	return bianPrice
}

func getHuobiLastPrice(symbol string) float64 {
	log.Println(strconv.FormatInt(utils.UnixMillis(time.Now()), 10) + " huobi")
	defer getPriceThread.Done()
	huobiPrice = api.HuobiLastPrice(symbol).Tick.Data[0].Price
	log.Println(strconv.FormatInt(utils.UnixMillis(time.Now()), 10) + " huobi: " + strconv.FormatFloat(huobiPrice, 'f', -1, 64))
	return huobiPrice
}

func getDiffPrice(a, b float64) float64 {
	return math.Abs(a - b)
}
