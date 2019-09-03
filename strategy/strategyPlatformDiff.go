package strategy

import (
	"AutoTrading/api"
	"AutoTrading/config"
	"AutoTrading/models"
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
	targetDiffPrice := huobiPrice * config.PlatformDiffPoint

	// When the platform funds are seriously unbalanced, the threshold transfer funds will be lowered according to the situation.
	bianUsdtValue := bianAccount.Usdt / bianPrice
	logPrefix := ""
	if huobiIsGreaterThanBian && bianAccount.Btc < bianUsdtValue && bianAccount.Btc/bianUsdtValue < 0.5 {
		targetDiffPrice = targetDiffPrice / 3
		logPrefix = "平衡资金, "
	} else if !huobiIsGreaterThanBian && bianAccount.Btc > bianUsdtValue && bianUsdtValue/bianAccount.Btc < 0.5 {
		targetDiffPrice = targetDiffPrice / 3
		logPrefix = "平衡资金, "
	}
	if bianPrice == 0 {
		log.Println("未获取到 Binance 的最新价格信息")
		return
	}
	if huobiPrice == 0 {
		log.Println("未获取到 Huobi 的最新价格信息")
		return
	}
	if currDiffPrice >= targetDiffPrice {
		if huobiIsGreaterThanBian {
			log.Println(fmt.Sprintf("差价 %.2f USD, 火币的价格大于币安的价格", currDiffPrice))
		} else {
			log.Println(fmt.Sprintf("差价 %.2f USD, 火币的价格小于币安的价格", currDiffPrice))
		}

		// base bian, 币安买卖都是数量, 火币买入用交易额，卖出用数量
		if huobiIsGreaterThanBian {
			// sell huobi, buy bian
			// Trading on both platforms when the transaction is successfully completed
			if logPrefix != "" {
				bianBuyBtcCount := bianUsdtValue - (bianAccount.Btc+bianUsdtValue)/2
				log.Println(fmt.Sprintf("%s开始交易...\n火币卖 %.6f BTC\n币安买 %.6f BTC", logPrefix, bianBuyBtcCount, bianBuyBtcCount))
				if isTest {
					tradeTest(bianBuyBtcCount, bianBuyBtcCount, huobiIsGreaterThanBian)
				} else {
					diffTrade(bianBuyBtcCount, bianBuyBtcCount, huobiIsGreaterThanBian, logPrefix)
				}
			} else if huobiAccount.Btc*bianPrice > bianAccount.Usdt {
				// huobi sell bianAccount.Usdt/huobiPrice, bian buy bianAccount.Usdt
				huobiSellBtcCount := bianAccount.Usdt / huobiPrice
				bianBuyBtcCount := bianAccount.Usdt / bianPrice
				// Trading when the transaction amount is less than 15 USD
				if bianAccount.Usdt < 15 {
					log.Println("Did not reach the minimum order transaction amount, no transaction")
					return
				}
				log.Println(fmt.Sprintf("%s开始交易...\n火币卖 %.6f BTC\n币安买 %.6f BTC", logPrefix, huobiSellBtcCount, bianBuyBtcCount))
				if isTest {
					tradeTest(huobiSellBtcCount, bianBuyBtcCount, huobiIsGreaterThanBian)
				} else {
					diffTrade(huobiSellBtcCount, bianBuyBtcCount, huobiIsGreaterThanBian, logPrefix)
				}
			} else {
				// huobi sell huobiAccount.Btc, bian buy huobiAccount.btc*bianPrice
				huobiSellBtcCount := huobiAccount.Btc
				bianBuyBtcCount := huobiAccount.Btc
				if huobiAccount.Btc*bianPrice < 15 {
					log.Println("Did not reach the minimum order transaction amount, no transaction")
					return
				}
				// Trading when the transaction amount is less than 15 USD
				log.Println(fmt.Sprintf("%s开始交易...\n火币卖 %.6f BTC\n币安买 %.6f BTC", logPrefix, huobiSellBtcCount, bianBuyBtcCount))
				if isTest {
					tradeTest(huobiSellBtcCount, bianBuyBtcCount, huobiIsGreaterThanBian)
				} else {
					diffTrade(huobiSellBtcCount, bianBuyBtcCount, huobiIsGreaterThanBian, logPrefix)
				}
			}
		} else {
			// buy huobi, sell bian
			if logPrefix != "" {
				bianSellCount := bianAccount.Btc - (bianAccount.Btc+bianUsdtValue)/2
				huobiBuy := bianSellCount * bianPrice
				log.Println(fmt.Sprintf("%s开始交易...\n火币花费 ​%.8f USD 买 BTC\n币安卖 %.6f BTC", logPrefix, huobiBuy, bianSellCount))
				if isTest {
					tradeTest(huobiBuy, bianSellCount, huobiIsGreaterThanBian)
				} else {
					diffTrade(huobiBuy, bianSellCount, huobiIsGreaterThanBian, logPrefix)
				}
			} else if huobiAccount.Usdt < bianAccount.Btc*huobiPrice {
				// huobi buy huobiAccount.Usdt, bian sell huobiAccount.Usdt/bianPrice
				huobiBuy := huobiAccount.Usdt
				bianSellBtcCount := huobiAccount.Usdt / bianPrice
				if huobiBuy < 15 {
					log.Println("Did not reach the minimum order transaction amount, no transaction")
					return
				}
				log.Println(fmt.Sprintf("%s开始交易...\n火币花费 ​%.8f USD 买 BTC\n币安卖 %.6f BTC", logPrefix, huobiBuy, bianSellBtcCount))
				if isTest {
					tradeTest(huobiBuy, bianSellBtcCount, huobiIsGreaterThanBian)
				} else {
					diffTrade(huobiBuy, bianSellBtcCount, huobiIsGreaterThanBian, logPrefix)
				}
			} else {
				// huobi buy bianAccount.btc*huobiPrice, bian sell bianAccount.Btc
				huobiBuy := bianAccount.Btc * huobiPrice
				bianSellBtcCount := bianAccount.Btc
				if huobiBuy < 15 {
					log.Println("Did not reach the minimum order transaction amount, no transaction")
					return
				}
				log.Println(fmt.Sprintf("%s开始交易...\n火币花费 ​%.8f USD 买 BTC\n币安卖 %.6f BTC", logPrefix, huobiBuy, bianSellBtcCount))
				if isTest {
					tradeTest(huobiBuy, bianSellBtcCount, huobiIsGreaterThanBian)
				} else {
					diffTrade(huobiBuy, bianSellBtcCount, huobiIsGreaterThanBian, logPrefix)
				}
			}
		}
	}
}

func diffTrade(huobiValue, bianValue float64, huobiIsGreaterThanBian bool, logPrefix string) {
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
	huobiQty, huobiAvgPrice := getHuobiAvgPrice(huobiOrderResult.Data)
	for huobiQty == 0 {
		time.Sleep(time.Duration(100) * time.Millisecond)
		huobiOrderResult = api.HuobiOrderQueryDetail(huobiResultOrderId)
		huobiQty, huobiAvgPrice = getHuobiAvgPrice(huobiOrderResult.Data)
	}
	bianQty, _ := strconv.ParseFloat(bianOrderResult.ExecutedQty, 64)
	_, bianAvgPrice := getBianAvgPrice(bianOrderResult.Fills)
	var logStr string
	if huobiIsGreaterThanBian {
		logStr = fmt.Sprintf("%s交易成功:\n火币卖 %.6f BTC, 获得 %.8f USD, 均价: %.2f USD, 订单号: %s\n币安买 %.6f BTC, 花费 %.8f USD, 均价: %.2f USD, 订单号: %d\n交易用时 %d 毫秒", logPrefix,
			huobiQty, huobiQty*huobiAvgPrice, huobiAvgPrice, huobiResultOrderId,
			bianQty, bianQty*bianAvgPrice, bianAvgPrice, bianOrderResult.OrderID,
			costTime)
		log.Println(logStr)
	} else {
		logStr = fmt.Sprintf("%s交易成功:\n火币买 %.6f BTC, 花费 %.8f USD, 均价: %.2f USD, 订单号: %s\n币安卖 %.6f BTC, 获得 %.8f USD, 均价: %.2f USD, 订单号: %d\n交易用时 %d 毫秒", logPrefix,
			huobiQty, huobiQty*huobiAvgPrice, huobiAvgPrice, huobiResultOrderId,
			bianQty, bianQty*bianAvgPrice, bianAvgPrice, bianOrderResult.OrderID,
			costTime)
		log.Println(logStr)
	}

	// TODO Save to database
	// TODO Save log to file
	// TODO Calculated income

	// reminder when the Commission is insufficient
	if config.Ifttt.Enabled {
		log.Println("发送 IFTTT 消息提醒...")
		api.IftttNotice("完成一笔交易", logStr, "")
	}

	log.Println("火币交易详情: ")
	log.Println(huobiOrderResult)
	log.Println("币安交易详情: ")
	log.Println(bianOrderResult)

	updateAccountBalance()
}

// return qty, avg price
func getHuobiAvgPrice(fills []models.HuobiFill) (float64, float64) {
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
func getBianAvgPrice(fills []models.BianFill) (float64, float64) {
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
		log.Println(fmt.Sprintf("Successful Transaction:\nHuo Bi: Sell %.6f BTC,  Get %.8f USD, Average Price: %.2f USD\nBinance: Buy %.6f BTC, Take %.8f USD, Average Price: %.2f USD\nTrading time %d milliseconds",
			huobiValue, huobiValue*testHuoBiPrice, testHuoBiPrice,
			bianValue, bianValue*testBianPrice, testBianPrice,
			costTime))
	} else {
		log.Println(fmt.Sprintf("Successful Transaction:\nHuo  Bi:  Buy %.6f BTC, Take %.8f USD, Average Price: %.2f USD\nBinance: Sell %.6f BTC,  Get %.8f USD, Average Price: %.2f USD\nTrading time %d milliseconds",
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

func getBianLastPrice(symbol string) float64 {
	startTime := utils.UnixMillis(time.Now())
	defer getPriceThread.Done()
	bianPrice, _ = strconv.ParseFloat(api.BianLastPrice(symbol).Price, 64)
	bianPrice = utils.Decimal(bianPrice, "2")
	costTime := utils.UnixMillis(time.Now()) - startTime
	if costTime > 150 {
		longTime = true
	}
	log.Println(fmt.Sprintf("The last %s price in Bian : %.2f USD, start at %d, takes %dms", symbol, bianPrice, startTime, costTime))
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
	log.Println(fmt.Sprintf("The last %s price in Huobi: %.2f USD, start at %d, takes %dms", symbol, huobiPrice, startTime, costTime))
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
