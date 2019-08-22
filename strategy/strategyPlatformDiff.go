package strategy

import (
	"AutoTrading/api"
	"AutoTrading/config"
	"AutoTrading/utils"
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"
	"sync"
	"time"
)

var getPriceThread = sync.WaitGroup{}
var bianPrice, huobiPrice float64

func RunPlatformDiffStrategy() {
	symbol := "BTCUSDT"
	symbolUpper := strings.ToUpper(symbol)
	symbolLowwer := strings.ToLower(symbol)
	getPriceThread.Add(1)
	go getHuobiLastPrice(symbolLowwer)
	getPriceThread.Add(1)
	go getBianLastPrice(symbolUpper)

	getPriceThread.Wait()
	diffPrice, gtA := getDiffPrice(huobiPrice, bianPrice)
	if diffPrice >= config.PlatformDiff.BTC {
		log.Println(fmt.Sprintf("diff price is %.10f, the Huobi Price is greater than the Bian Price: %t", diffPrice, gtA))
		if gtA {

		} else {

		}
	}
}

// TODO Platforms with low balances sell more
// TODO Send alerts when there is a serious imbalance

func getBianLastPrice(symbol string) float64 {
	startTime := utils.UnixMillis(time.Now())
	defer getPriceThread.Done()
	bianPrice, _ = strconv.ParseFloat(api.BianLastPrice(symbol).Price, 64)
	costTime := utils.UnixMillis(time.Now()) - startTime
	log.Println(fmt.Sprintf("The last %s price in Bian : %.10f, start at %d, takes %dms", symbol, bianPrice, startTime, costTime))
	return bianPrice
}

func getHuobiLastPrice(symbol string) float64 {
	startTime := utils.UnixMillis(time.Now())
	defer getPriceThread.Done()
	huobiPrice = api.HuobiLastPrice(symbol).Tick.Data[0].Price
	costTime := utils.UnixMillis(time.Now()) - startTime
	log.Println(fmt.Sprintf("The last %s price in Huobi: %.10f, start at %d, takes %dms", symbol, huobiPrice, startTime, costTime))
	return huobiPrice
}

func getDiffPrice(a, b float64) (float64, bool) {
	return math.Abs(a - b), a > b
}