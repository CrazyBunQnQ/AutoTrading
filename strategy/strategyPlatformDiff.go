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
var huobiResultOrderId string
var bianOrderResult models.BianFullOrder
var symbol = "BTCUSDT"
var symbolUpper = strings.ToUpper(symbol)
var symbolLowwer = strings.ToLower(symbol)
var longTime = false

func RunPlatformDiffStrategy(isTest bool) {
	// Query the current balance of each platform account
	updateAccountBalance()
	for true {
		startPlatformDiffStrategy(isTest)
		time.Sleep(time.Duration(1) * time.Second)
	}
}

func startPlatformDiffStrategy(isTest bool) {

	getPriceThread.Add(1)
	go getHuobiLastPrice(symbolLowwer)
	getPriceThread.Add(1)
	go getBianLastPrice(symbolUpper)
	getPriceThread.Wait()
	if longTime {
		longTime = false
		//log.Println("Take too long to give up this transaction...")
		return
	}

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
		log.Println(fmt.Sprintf("Diff price is %.2f USD, the Huobi Price is greater than the Bian Price: %t", currDiffPrice, huobiIsGreaterThanBian))
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
				log.Println(fmt.Sprintf("Trade start...\nSell ​%.2f BTC on the Huobi\n Buy %.2f BTC on the Binance", huobiSellBtcCount, bianBuyBtcCount))
				if isTest {
					tradeTest(huobiSellBtcCount, bianBuyBtcCount, huobiIsGreaterThanBian)
				} else {
					diffTrade(huobiSellBtcCount, bianBuyBtcCount, huobiIsGreaterThanBian)
				}
			} else {
				// todo huobi sell huobiAccount.Btc, bian buy huobiAccount.btc*bianPrice
				huobiSellBtcCount := huobiAccount.Btc
				bianBuyBtcCount := huobiAccount.Btc
				if huobiAccount.Btc*bianPrice < 15 {
					log.Println("Did not reach the minimum order transaction amount, no transaction")
					return
				}
				// Trading when the transaction amount is less than 15 USD
				log.Println(fmt.Sprintf("Trade start...\nSell ​%.2f BTC on the Huobi\n Buy %.2f BTC on the Binance", huobiSellBtcCount, bianBuyBtcCount))
				if isTest {
					tradeTest(huobiSellBtcCount, bianBuyBtcCount, huobiIsGreaterThanBian)
				} else {
					diffTrade(huobiSellBtcCount, bianBuyBtcCount, huobiIsGreaterThanBian)
				}
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
				log.Println(fmt.Sprintf("Trade start...\nSpend ​%.8f USD on the Huobi to buy BTC\nSell %.2f BTC on the Binance", huobiBuy, bianSellBtcCount))
				if isTest {
					tradeTest(huobiBuy, bianSellBtcCount, huobiIsGreaterThanBian)
				} else {
					diffTrade(huobiBuy, bianSellBtcCount, huobiIsGreaterThanBian)
				}
			} else {
				// todo huobi buy bianAccount.btc*huobiPrice, bian sell bianAccount.Btc
				huobiBuy := bianAccount.Btc * huobiPrice
				bianSellBtcCount := bianAccount.Btc
				if huobiBuy < 15 {
					log.Println("Did not reach the minimum order transaction amount, no transaction")
					return
				}
				log.Println(fmt.Sprintf("Trade start...\nSpend ​%.8f USD on the Huobi to buy BTC\nSell %.2f BTC on the Binance", huobiBuy, bianSellBtcCount))
				if isTest {
					tradeTest(huobiBuy, bianSellBtcCount, huobiIsGreaterThanBian)
				} else {
					diffTrade(huobiBuy, bianSellBtcCount, huobiIsGreaterThanBian)
				}
			}
		}
	} else if currDiffPrice < targetBalancedDiffPrice {
		// When Binance's funds are seriously out of balance, it will trade to balance the capital when the price difference is lowest
		bianUsdtValue := bianAccount.Usdt / bianPrice
		if huobiIsGreaterThanBian && bianAccount.Btc < bianUsdtValue && bianAccount.Btc/bianUsdtValue < 0.5 {
			// sell huobi, buy bian
			bianBuyBtcCount := bianUsdtValue - (bianAccount.Btc+bianUsdtValue)/2
			log.Println(fmt.Sprintf("Balanced funds, Trade start...\nSell %.2f BTC on the Huobi\n Buy %.2f BTC on the Binance", bianBuyBtcCount, bianBuyBtcCount))
			if isTest {
				tradeTest(bianBuyBtcCount, bianBuyBtcCount, huobiIsGreaterThanBian)
			} else {
				diffTrade(bianBuyBtcCount, bianBuyBtcCount, huobiIsGreaterThanBian)
			}
		} else if !huobiIsGreaterThanBian && bianAccount.Btc > bianUsdtValue && bianUsdtValue/bianAccount.Btc < 0.5 {
			// buy huobi, sell bian
			bianSellCount := bianAccount.Btc - (bianAccount.Btc+bianUsdtValue)/2
			huobiBuy := bianSellCount * bianPrice
			log.Println(fmt.Sprintf("Balanced funds, Trade start...\nSpend ​%.8f USD on the Huobi to buy BTC\nSell %.2f BTC on the Binance", huobiBuy, bianSellCount))
			if isTest {
				tradeTest(huobiBuy, bianSellCount, huobiIsGreaterThanBian)
			} else {
				diffTrade(huobiBuy, bianSellCount, huobiIsGreaterThanBian)
			}
		}
	}
}

func diffTrade(huobiValue, bianValue float64, huobiIsGreaterThanBian bool) {
	startTime := utils.UnixMillis(time.Now())
	var huobiSide string
	var bianSide models.BianOrderSide
	if huobiIsGreaterThanBian {
		huobiSide = "sell"
		bianSide = models.SideBuy
	} else {
		huobiSide = "buy"
		bianSide = models.SideSell
	}
	tradeThread.Add(1)
	go diffTradeHuobi(symbolLowwer, huobiSide, huobiValue*0.9)
	tradeThread.Add(1)
	go diffTradeBian(symbolUpper, bianSide, bianValue*0.9)
	tradeThread.Wait()

	costTime := utils.UnixMillis(time.Now()) - startTime

	// Query order information
	huobiOrderResult := api.HuobiOrderQueryDetail(huobiResultOrderId)
	huobiQty, huobiAvgPrice := huobiAvgPrice(huobiOrderResult.Data)
	bianQty, _ := strconv.ParseFloat(bianOrderResult.ExecutedQty, 64)
	_, bianAvgPrice := bianAvgPrice(bianOrderResult.Fills)
	var logStr string
	if huobiIsGreaterThanBian {
		logStr = fmt.Sprintf("Successful Transaction:\nHuo Bi: Sell %.2f BTC,  Get %.8f USD, Average Price: %.8f USD, OrderID: %s\nBinance: Buy %.2f BTC, Take %.8f USD, Average Price: %.8f USD, OrderID: %d\nTrading time %d milliseconds",
			huobiQty, huobiQty*huobiAvgPrice, huobiAvgPrice, huobiResultOrderId,
			bianQty, bianQty*bianAvgPrice, bianAvgPrice, bianOrderResult.OrderID,
			costTime)
		log.Println(logStr)
	} else {
		logStr = fmt.Sprintf("Successful Transaction:\nHuo  Bi:  Buy %.2f BTC, Take %.8f USD, Average Price: %.8f USD, OrderID: %s\nBinance: Sell %.2f BTC,  Get %.8f USD, Average Price: %.8f USD, OrderID: %d\nTrading time %d milliseconds",
			huobiQty, huobiQty*huobiAvgPrice, huobiAvgPrice, huobiResultOrderId,
			bianQty, bianQty*bianAvgPrice, bianAvgPrice, bianOrderResult.OrderID,
			costTime)
		log.Println(logStr)
	}
	// reminder when the Commission is insufficient
	if config.Ifttt.Enabled {
		api.IftttNotice("完成一笔差价交易", logStr, "")
	}

	updateAccountBalance()
}

// return qty, avg price
func huobiAvgPrice(fills []models.HuobiFill) (float64, float64) {
	var totalQty float64 = 0
	var totalUsdt float64 = 0
	for _, fill := range fills {
		curQty, _ := strconv.ParseFloat(fill.FilledAmount, 64)
		curPrice, _ := strconv.ParseFloat(fill.Price, 64)
		totalQty = totalQty + curQty
		totalUsdt = totalUsdt + curQty*curPrice
	}
	return totalQty, totalUsdt / totalQty
}

// return qty, avg price
func bianAvgPrice(fills []models.BianFill) (float64, float64) {
	var totalQty float64 = 0
	var totalUsdt float64 = 0
	for _, fill := range fills {
		curQty, _ := strconv.ParseFloat(fill.Qty, 64)
		curPrice, _ := strconv.ParseFloat(fill.Price, 64)
		totalQty = totalQty + curQty
		totalUsdt = totalUsdt + curQty*curPrice
	}
	return totalQty, totalUsdt / totalQty
}

//amount: 限价表示下单数量, 市价买单时表示买多少钱, 市价卖单时表示卖多少币
func diffTradeHuobi(symbol, side string, amount float64) {
	defer tradeThread.Done()
	if side == "sell" {
		amount = utils.Decimal(amount, "6")
	} else {
		amount = utils.Decimal(amount, "8")
	}
	result := api.HuobiOrderByMarket(huobiAccountId, symbol, side, amount)
	log.Println("Huobi order result:\n" + result.Data)
	if result.Status == "error" {
		log.Println(fmt.Sprintf("Trade Error on Huobi: %s -> %s", result.ErrCode, result.ErrMsg))
	} else {
		huobiResultOrderId = result.Data
	}
}

func diffTradeBian(symbol string, side models.BianOrderSide, amount float64) {
	defer tradeThread.Done()
	bianOrderResult = api.BianOrderByMarket(symbol, side, utils.Decimal(amount, "6"), 0)
	if bianOrderResult.Err != "" {
		log.Println(fmt.Sprintf("Trade Error on Binance: %s", bianOrderResult.Err))
	}
}

var tradeTestThread = sync.WaitGroup{}
var testHuoBiPrice, testBianPrice float64

func tradeTest(huobiValue, bianValue float64, huobiIsGreaterThanBian bool) {
	startTime := utils.UnixMillis(time.Now())
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
	bianPrice = utils.Decimal(bianPrice, "2")
	costTime := utils.UnixMillis(time.Now()) - startTime
	if costTime > 150 {
		longTime = true
	}
	log.Println(fmt.Sprintf("The last %s price in Bian : %.10f USD, start at %d, takes %dms", symbol, bianPrice, startTime, costTime))
	return bianPrice
}

func getHuobiLastPrice(symbol string) float64 {
	startTime := utils.UnixMillis(time.Now())
	defer getPriceThread.Done()
	// 减去偏移量
	huobiPrice = api.HuobiLastPrice(symbol).Tick.Data[0].Price - config.PlatformOffset
	costTime := utils.UnixMillis(time.Now()) - startTime
	if costTime > 150 {
		longTime = true
	}
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
			account.Usdt = utils.Decimal(num, "8")
			locked, _ := strconv.ParseFloat(balance.Locked, 64)
			account.UsdtLocked = utils.Decimal(locked, "8")
		case "BTC":
			num, _ := strconv.ParseFloat(balance.Free, 64)
			account.Btc = utils.Decimal(num, "6")
			locked, _ := strconv.ParseFloat(balance.Locked, 64)
			account.BtcLocked = utils.Decimal(locked, "6")
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
				account.Usdt = utils.Decimal(num, "8")
			} else {
				account.UsdtLocked = utils.Decimal(num, "8")
			}
		case "btc":
			if balance.Type == "trade" {
				account.Btc = utils.Decimal(num, "6")
			} else {
				account.BtcLocked = utils.Decimal(num, "6")
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
