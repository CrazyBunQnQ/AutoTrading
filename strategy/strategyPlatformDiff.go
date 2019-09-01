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

var bianPrice, bianActualTradePrice, huobiPrice, huobiActualTradePrice float64
var huobiAccountId = api.GetAccounts().Data[0].ID
var bianAccount, huobiAccount models.Account

func RunPlatformDiffStrategy() {
	// Query the current balance of each platform account
	updateAccountBalance()
	for true {
		startPlatformDiffStrategy()
		time.Sleep(time.Duration(3) * time.Second)
	}
}

func startPlatformDiffStrategy() {

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
	targetDiffPrice := huobiPrice * config.PlatformDiffPoint
	targetBalancedDiffPrice := huobiPrice * config.PlatformBalancedPoint
	if bianPrice == 0 {
		log.Println("未获取到 Binance 的最新价格信息")
		return
	}
	if huobiPrice == 0 {
		log.Println("未获取到 Huobi 的最新价格信息")
		return
	}
	if currDiffPrice >= targetDiffPrice {
		log.Println(fmt.Sprintf("Diff price is %.10f USD, the Huobi Price is greater than the Bian Price: %t", currDiffPrice, huobiIsGreaterThanBian))
		// base bian, 币安买卖都是数量, 火币买入用交易额，卖出用数量
		if huobiIsGreaterThanBian {
			// sell huobi, buy bian
			// Trading on both platforms when the transaction is successfully completed
			if huobiAccount.Btc*bianPrice > bianAccount.Usdt {
				// todo huobi sell bianAccount.Usdt/huobiPrice, bian buy bianAccount.Usdt
				huobiSellBtcCount := bianAccount.Usdt / huobiPrice
				bianBuyBtcCount := bianAccount.Usdt / bianPrice
				// Trading when the transaction amount is less than 15 USD
				if bianAccount.Usdt < 15 {
					log.Println("Did not reach the minimum order transaction amount, no transaction")
					return
				}
				log.Println(fmt.Sprintf("Trade start...\nSell ​%.10f BTC on the Huobi\n Buy %.10f BTC on the Binance", huobiSellBtcCount, bianBuyBtcCount))
				tradeTest(huobiSellBtcCount, bianBuyBtcCount, huobiIsGreaterThanBian)
				updateAccountBalance()
			} else {
				// todo huobi sell huobiAccount.Btc, bian buy huobiAccount.btc*bianPrice
				huobiSellBtcCount := huobiAccount.Btc
				bianBuyBtcCount := huobiAccount.Btc
				if huobiAccount.Btc*bianPrice < 15 {
					log.Println("Did not reach the minimum order transaction amount, no transaction")
					return
				}
				// Trading when the transaction amount is less than 15 USD
				log.Println(fmt.Sprintf("Trade start...\nSell ​%.10f BTC on the Huobi\n Buy %.10f BTC on the Binance", huobiSellBtcCount, bianBuyBtcCount))
				tradeTest(huobiSellBtcCount, bianBuyBtcCount, huobiIsGreaterThanBian)
				updateAccountBalance()
			}
		} else {
			// buy huobi, sell bian
			if huobiAccount.Usdt < bianAccount.Btc*huobiPrice {
				// todo huobi buy huobiAccount.Usdt, bian sell huobiAccount.Usdt/bianPrice
				huobiBuy := huobiAccount.Usdt
				bianSellBtcCount := huobiAccount.Usdt / bianPrice
				if huobiBuy < 15 {
					log.Println("Did not reach the minimum order transaction amount, no transaction")
					return
				}
				log.Println(fmt.Sprintf("Trade start...\nSpend ​%.10f USD on the Huobi to buy BTC\nSell %.10f BTC on the Binance", huobiBuy, bianSellBtcCount))
				tradeTest(huobiBuy, bianSellBtcCount, huobiIsGreaterThanBian)
				updateAccountBalance()
			} else {
				// todo huobi buy bianAccount.btc*huobiPrice, bian sell bianAccount.Btc
				huobiBuy := bianAccount.Btc * huobiPrice
				bianSellBtcCount := bianAccount.Btc
				if huobiBuy < 15 {
					log.Println("Did not reach the minimum order transaction amount, no transaction")
					return
				}
				log.Println(fmt.Sprintf("Trade start...\nSpend ​%.10f USD on the Huobi to buy BTC\nSell %.10f BTC on the Binance", huobiBuy, bianSellBtcCount))
				tradeTest(huobiBuy, bianSellBtcCount, huobiIsGreaterThanBian)
				updateAccountBalance()
			}
		}
	} else if currDiffPrice < targetBalancedDiffPrice {
		// When Binance's funds are seriously out of balance, it will trade to balance the capital when the price difference is lowest
		bianUsdtValue := bianAccount.Usdt / bianPrice
		if huobiIsGreaterThanBian && bianAccount.Btc < bianUsdtValue && bianAccount.Btc/bianUsdtValue < 0.5 {
			// sell huobi, buy bian
			bianBuyBtcCount := bianUsdtValue - (bianAccount.Btc+bianUsdtValue)/2
			log.Println(fmt.Sprintf("Balanced funds, Trade start...\nSell ​%.10f BTC on the Huobi\n Buy %.10f BTC on the Binance", bianBuyBtcCount, bianBuyBtcCount))
			tradeTest(bianBuyBtcCount, bianBuyBtcCount, huobiIsGreaterThanBian)
			updateAccountBalance()
		} else if !huobiIsGreaterThanBian && bianAccount.Btc > bianUsdtValue && bianUsdtValue/bianAccount.Btc < 0.5 {
			// buy huobi, sell bian
			bianSellCount := bianAccount.Btc - (bianAccount.Btc+bianUsdtValue)/2
			huobiBuy := bianSellCount * bianPrice
			log.Println(fmt.Sprintf("Balanced funds, Trade start...\nSpend ​%.10f USD on the Huobi to buy BTC\nSell %.10f BTC on the Binance", huobiBuy, bianSellCount))
			tradeTest(huobiBuy, bianSellCount, huobiIsGreaterThanBian)
			updateAccountBalance()
		}
	}
}

func diffTrade(huobiValue, bianValue float64, huobiIsGreaterThanBian bool) {
	startTime := utils.UnixMillis(time.Now())
	symbol := "BTCUSDT"
	symbolUpper := strings.ToUpper(symbol)
	symbolLowwer := strings.ToLower(symbol)
	tradeThread.Add(1)
	go diffTradeHuobi(symbolLowwer)
	tradeThread.Add(1)
	go diffTradeBian(symbolUpper)
	tradeThread.Wait()

	costTime := utils.UnixMillis(time.Now()) - startTime
	if huobiIsGreaterThanBian {
		log.Println(fmt.Sprintf("Successful Transaction:\nHuo Bi: Sell %.10f BTC,  Get %.10f USD, Average Price: %.10f USD\nBinance: Buy %.10f BTC, Take %.10f USD, Average Price: %.10f USD\nTrading time %d milliseconds",
			huobiValue, huobiValue*testHuoBiPrice, testHuoBiPrice,
			bianValue, bianValue*testBianPrice, testBianPrice,
			costTime))
	} else {
		log.Println(fmt.Sprintf("Successful Transaction:\nHuo  Bi:  Buy %.10f BTC, Take %.10f USD, Average Price: %.10f USD\nBinance: Sell %.10f BTC,  Get %.10f USD, Average Price: %.10f USD\nTrading time %d milliseconds",
			huobiValue/testHuoBiPrice, huobiValue, testHuoBiPrice,
			bianValue, bianValue*testBianPrice, testBianPrice,
			costTime))
	}
}

func diffTradeHuobi(symbol string) {
	defer tradeThread.Done()
	huobiActualTradePrice = api.HuobiLastPrice(symbol).Tick.Data[0].Price
}

func diffTradeBian(symbol string) {
	defer tradeThread.Done()
	api.BianOrderByMarket(symbol, models.SideBuy, 0.002, 0)
	bianActualTradePrice, _ = strconv.ParseFloat(api.BianLastPrice(symbol).Price, 64)
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
		log.Println(fmt.Sprintf("Successful Transaction:\nHuo Bi: Sell %.10f BTC,  Get %.10f USD, Average Price: %.10f USD\nBinance: Buy %.10f BTC, Take %.10f USD, Average Price: %.10f USD\nTrading time %d milliseconds",
			huobiValue, huobiValue*testHuoBiPrice, testHuoBiPrice,
			bianValue, bianValue*testBianPrice, testBianPrice,
			costTime))
	} else {
		log.Println(fmt.Sprintf("Successful Transaction:\nHuo  Bi:  Buy %.10f BTC, Take %.10f USD, Average Price: %.10f USD\nBinance: Sell %.10f BTC,  Get %.10f USD, Average Price: %.10f USD\nTrading time %d milliseconds",
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
	defer tradeTestThread.Done()
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

func getBianLastPrice(symbol string) float64 {
	startTime := utils.UnixMillis(time.Now())
	defer getPriceThread.Done()
	bianPrice, _ = strconv.ParseFloat(api.BianLastPrice(symbol).Price, 64)
	costTime := utils.UnixMillis(time.Now()) - startTime
	log.Println(fmt.Sprintf("The last %s price in Bian : %.10f USD, start at %d, takes %dms", symbol, bianPrice, startTime, costTime))
	return bianPrice
}

func getHuobiLastPrice(symbol string) float64 {
	startTime := utils.UnixMillis(time.Now())
	defer getPriceThread.Done()
	huobiPrice = api.HuobiLastPrice(symbol).Tick.Data[0].Price
	costTime := utils.UnixMillis(time.Now()) - startTime
	log.Println(fmt.Sprintf("The last %s price in Huobi: %.10f USD, start at %d, takes %dms", symbol, huobiPrice, startTime, costTime))
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
