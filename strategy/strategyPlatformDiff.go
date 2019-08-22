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
		//log.Println(fmt.Sprintf("\nOrder at Binance: %s %.10f %s at the price of %.10f\nresult: %s", side, quantity, symbol, price, order.Err))
		log.Println(fmt.Sprintf("diff price is %.10f, the Huobi Price is greater than the Bian Price: %t", diffPrice, gtA))
		if gtA {

		} else {

		}
	}
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

func getDiffPrice(a, b float64) (float64, bool) {
	return math.Abs(a - b), a > b
}
