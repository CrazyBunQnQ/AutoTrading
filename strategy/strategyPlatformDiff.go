package strategy

import (
	"AutoTrading/api"
	"AutoTrading/config"
	"AutoTrading/models"
	"AutoTrading/utils"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"math"
	"strconv"
	"strings"
	"sync"
	"time"
)

var getPriceThread = sync.WaitGroup{}
var tradeThread = sync.WaitGroup{}

var bianPrice, huobiPrice float64
var bianAccount, huobiAccount models.Account

func RunPlatformDiffStrategy() {
	// Query the current balance of each platform account
	updateAccountBalance()

	symbol := "BTCUSDT"
	symbolUpper := strings.ToUpper(symbol)
	symbolLowwer := strings.ToLower(symbol)
	getPriceThread.Add(1)
	go getHuobiLastPrice(symbolLowwer)
	getPriceThread.Add(1)
	go getBianLastPrice(symbolUpper)
	getPriceThread.Wait()

	currDiffPrice, huobiIsGreaterThanBian := getDiffPrice(huobiPrice, bianPrice)
	// TODO When the platform funds are seriously unbalanced, the threshold transfer funds will be lowered according to the situation.
	targetDiffPrice := huobiPrice * config.PlatformDiff
	if bianPrice != 0 && huobiPrice != 0 && currDiffPrice >= targetDiffPrice {
		log.Println(fmt.Sprintf("diff price is %.10f, the Huobi Price is greater than the Bian Price: %t", currDiffPrice, huobiIsGreaterThanBian))
		// base bian, 币安买卖都是数量, 火币买入用交易额，卖出用数量
		if huobiIsGreaterThanBian {
			// TODO sell huobi, buy bian
			// Trading on both platforms when the transaction is successfully completed
			if huobiAccount.Btc*bianPrice > bianAccount.Usdt {
				// todo huobi sell bianAccount.Usdt/huobiPrice, bian buy bianAccount.Usdt
				huobiSellCount := bianAccount.Usdt / huobiPrice
				bianBuyCount := bianAccount.Usdt / bianPrice
				// TODO Trading when the transaction amount is less than 15 USD
				if bianAccount.Usdt < 15 {
					log.Println("Did not reach the minimum order transaction amount, no transaction")
					return
				}
				log.Println(fmt.Sprintf("Sell ​%.10f BTC on the Huobi\nBuy %.10f BTC on the Binance", huobiSellCount, bianBuyCount))
				tradeTest(huobiSellCount, bianBuyCount, huobiIsGreaterThanBian)
			} else {
				// todo huobi sell huobiAccount.Btc, bian buy huobiAccount.btc*bianPrice
				huobiSellCount := huobiAccount.Btc
				bianBuyCount := huobiAccount.Btc
				if huobiAccount.Btc*bianPrice < 15 {
					log.Println("Did not reach the minimum order transaction amount, no transaction")
					return
				}
				// TODO Trading when the transaction amount is less than 15 USD
				log.Println(fmt.Sprintf("Sell ​%.10f BTC on the Huobi\nBuy %.10f BTC on the Binance", huobiSellCount, bianBuyCount))
				tradeTest(huobiSellCount, bianBuyCount, huobiIsGreaterThanBian)
			}
		} else {
			// TODO buy huobi, sell bian
			if huobiAccount.Usdt < bianAccount.Btc*huobiPrice {
				// todo huobi buy huobiAccount.Usdt, bian sell huobiAccount.Usdt/bianPrice
				huobiBuy := huobiAccount.Usdt
				bianSellCount := huobiAccount.Usdt / bianPrice
				if huobiBuy < 15 {
					log.Println("Did not reach the minimum order transaction amount, no transaction")
					return
				}
				log.Println(fmt.Sprintf("Spend ​$%.10f on the Huobi to buy BTC\nSell %.10f BTC on the Binance", huobiBuy, bianSellCount))
				tradeTest(huobiBuy, bianSellCount, huobiIsGreaterThanBian)
			} else {
				// todo huobi buy bianAccount.btc*huobiPrice, bian sell bianAccount.Btc
				huobiBuy := bianAccount.Btc * huobiPrice
				bianSellCount := bianAccount.Btc
				if huobiBuy < 15 {
					log.Println("Did not reach the minimum order transaction amount, no transaction")
					return
				}
				log.Println(fmt.Sprintf("Spend ​$%.10f on the Huobi to buy BTC\nSell %.10f BTC on the Binance", huobiBuy, bianSellCount))
				tradeTest(huobiBuy, bianSellCount, huobiIsGreaterThanBian)
			}
		}
	}
}

var tradeTestThread = sync.WaitGroup{}
var testHuoBiPrice, testBianPrice float64

func tradeTest(huobiValue, bianValue float64, huobiIsGreaterThanBian bool) {
	startTime := utils.UnixMillis(time.Now())
	symbol := "BTCUSDT"
	symbolUpper := strings.ToUpper(symbol)
	symbolLowwer := strings.ToLower(symbol)
	tradeTestThread.Add(1)
	go tradeHuobiTest(symbolLowwer)
	tradeTestThread.Add(1)
	go tradeBianTest(symbolUpper)
	tradeTestThread.Wait()

	costTime := utils.UnixMillis(time.Now()) - startTime
	if huobiIsGreaterThanBian {
		log.Println(fmt.Sprintf("Successful Transaction:\nHuo Bi: Sell %.10f BTC,  Get $%.10f, Average Price: $%.10f\nBinance: Buy %.10f BTC, Take $%.10f, Average Price: $%.10f\nTrading time %d milliseconds",
			huobiValue, huobiValue*testHuoBiPrice, testHuoBiPrice,
			bianValue, bianValue*testBianPrice, testBianPrice,
			costTime))
	} else {
		log.Println(fmt.Sprintf("Successful Transaction:\nHuo  Bi:  Buy %.10f BTC, Take $%.10f, Average Price: $%.10f\nBinance: Sell %.10f BTC,  Get $%.10f, Average Price: $%.10f\nTrading time %d milliseconds",
			huobiValue/testHuoBiPrice, huobiValue, testHuoBiPrice,
			bianValue, bianValue*testBianPrice, testBianPrice,
			costTime))
	}
}

func tradeHuobiTest(symbol string) {
	defer tradeTestThread.Done()
	testHuoBiPrice = api.HuobiLastPrice(symbol).Tick.Data[0].Price
}

func tradeBianTest(symbol string) {
	defer getPriceThread.Done()
	testBianPrice, _ = strconv.ParseFloat(api.BianLastPrice(symbol).Price, 64)
}

// binance-exchange/go-binance/service_websocket.go
func DepthWebsocket(symbol string) (chan *models.DepthEvent, chan struct{}, error) {
	url := fmt.Sprintf("wss://stream.binance.com:9443/ws/%s@depth", symbol)
	_, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	done := make(chan struct{})

	return nil, done, nil
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

func updateAccountBalance() {
	huobiAccount = getHuobiBalance()
	bianAccount = getBianBalance()
}

func getHuobiBalance() models.Account {
	huobiAccountId := api.GetAccounts().Data[0].ID
	huobiBalanceReturn := api.GetAccountBalance(strconv.FormatInt(huobiAccountId, 10))
	return formatHuobiBalance(huobiBalanceReturn)
}

func getBianBalance() models.Account {
	bianAccount := api.BianAccountInfo()
	if bianAccount.Err != "" {
		log.Println("请求失败")
		return models.Account{Platform: "binance"}
	}
	return formatBianBalance(bianAccount)
}

func formatBianBalance(bianAccount models.BianAccount) models.Account {
	account := models.Account{Platform: "binance"}
	bianBalances := bianAccount.Balances
	for _, balance := range bianBalances {
		switch balance.Name {
		case "USDT":
			num, _ := strconv.ParseFloat(balance.Free, 64)
			account.Usdt = num
			locked, _ := strconv.ParseFloat(balance.Locked, 64)
			account.UsdtLocked = locked
		case "BTC":
			num, _ := strconv.ParseFloat(balance.Free, 64)
			account.Btc = num
			locked, _ := strconv.ParseFloat(balance.Locked, 64)
			account.BtcLocked = locked
		case "ETH":
			num, _ := strconv.ParseFloat(balance.Free, 64)
			account.Eth = num
			locked, _ := strconv.ParseFloat(balance.Locked, 64)
			account.EthLocked = locked
		case "BNB":
			num, _ := strconv.ParseFloat(balance.Free, 64)
			account.Bnb = num
			locked, _ := strconv.ParseFloat(balance.Locked, 64)
			account.BnbLocked = locked
		case "EOS":
			num, _ := strconv.ParseFloat(balance.Free, 64)
			account.Eos = num
			locked, _ := strconv.ParseFloat(balance.Locked, 64)
			account.EosLocked = locked
		case "XRP":
			num, _ := strconv.ParseFloat(balance.Free, 64)
			account.Xrp = num
			locked, _ := strconv.ParseFloat(balance.Locked, 64)
			account.XrpLocked = locked
		}
	}
	return account
}

func formatHuobiBalance(balanceReturn models.BalanceReturn) models.Account {
	account := models.Account{Platform: "huobi"}

	balances := balanceReturn.Data.List
	for _, balance := range balances {
		num, _ := strconv.ParseFloat(balance.Balance, 64)
		switch balance.Currency {
		case "usdt":
			if balance.Type == "trade" {
				account.Usdt = num
			} else {
				account.UsdtLocked = num
			}
		case "btc":
			if balance.Type == "trade" {
				account.Btc = num
			} else {
				account.BtcLocked = num
			}
		case "eth":
			if balance.Type == "trade" {
				account.Eth = num
			} else {
				account.EthLocked = num
			}
		case "bnb":
			if balance.Type == "trade" {
				account.Bnb = num
			} else {
				account.BnbLocked = num
			}
		case "eos":
			if balance.Type == "trade" {
				account.Eos = num
			} else {
				account.EosLocked = num
			}
		case "xrp":
			if balance.Type == "trade" {
				account.Xrp = num
			} else {
				account.XrpLocked = num
			}
		}
	}

	return account
}
