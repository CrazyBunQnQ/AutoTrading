package main

import (
	"AutoTrading/api"
	"AutoTrading/models"
	"AutoTrading/utils"
	"fmt"
	"github.com/astaxie/beego/orm"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"
)

var o orm.Ormer

func init() {
	orm.RegisterDataBase("default", "mysql", "root:zy26T$b7V8i3g4mW@tcp(127.0.0.1:3306)/autotrade?charset=utf8mb4", 30)
	orm.RegisterModel(new(models.StrategyLowBuyHighSell))
	orm.RegisterModel(new(models.Account))
	orm.RegisterModel(new(models.Quantity))
	o = orm.NewOrm()
}

func main() {

	//bianPrice := api.BianTrades("BTCUSDT", "1").Get(0, "price").ToFloat64()
	//log.Println("bian cur btc price: " + strconv.FormatFloat(bianPrice, 'f', -1, 64))
	//huobiPrice := api.GetTrade("btcusdt", 1)
	//log.Println("huobi cur btc price: " + strconv.FormatFloat(huobiPrice.Data[0].Data[0].Price, 'f', -1, 64))
	//otcbtcPrice := api.OtcbtcTrades("btcusdt", "1").Get(0, "price").ToFloat64()
	//log.Println("otcbtc cur btc price: " + strconv.FormatFloat(otcbtcPrice, 'f', -1, 64))
	//log.Println(api.BianOrderQuery("XRPUSDT", "", 207779114))
	http.HandleFunc("/", hello)
	http.ListenAndServe(":8000", nil)
}

func hello(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "hello world\n")
	startTime := utils.UnixMillis(time.Now())
	var status = api.BianOrderQuery("XRPUSDT", "", 207779114)
	costTime := utils.UnixMillis(time.Now()) - startTime

	if status.Price != "" {
		io.WriteString(w, "BianOrderQuery cost time: "+strconv.FormatInt(costTime, 10)+"ms\n")
	} else if status.Err != "" {
		io.WriteString(w, status.Err+"\n")
	} else {
		io.WriteString(w, "无返回")
	}
	startTime = utils.UnixMillis(time.Now())
	api.BianDepth("BTCUSDT", 5)
	costTime = utils.UnixMillis(time.Now()) - startTime
	io.WriteString(w, "BianDepth cost time: "+strconv.FormatInt(costTime, 10)+"ms\n")

	startTime = utils.UnixMillis(time.Now())
	api.HuobiDepth("btcusdt", "step0")
	costTime = utils.UnixMillis(time.Now()) - startTime
	io.WriteString(w, "HuobiDepth cost time: "+strconv.FormatInt(costTime, 10)+"ms\n")
}

func monthlyAvg() {
	// TODO Calculate/get monthly average price every day
}

func updateStrategyLBHS(name, platform string) models.StrategyLowBuyHighSell {
	// TODO Set the usdt case based on the name case
	strategyLBHS := queryStrategyLBHS(name+"USDT", platform)
	quantity := queryQuantity(name, platform)
	strategyLBHS.Quantity = quantity.Free + quantity.Locked
	if strategyLBHS.Spend != 0 {
		strategyLBHS.PositionAverage = strategyLBHS.Spend / strategyLBHS.Quantity
	}
	if strategyLBHS.PositionAverage != 0 {
		strategyLBHS.TargetSellPrice = strategyLBHS.PositionAverage * strategyLBHS.TargetProfitPoint
		strategyLBHS.TargetBuyPrice = strategyLBHS.PositionAverage * strategyLBHS.TargetBuyPoint
	}
	//	TODO other updates
	// monthAverage

	if _, err := o.Update(&strategyLBHS); err == nil {
		fmt.Println("Update a Strategy object:", strategyLBHS.String())
	}
	return strategyLBHS
}

func queryStrategyLBHS(symbol, platform string) models.StrategyLowBuyHighSell {
	data := models.StrategyLowBuyHighSell{Symbol: symbol, Platform: platform}
	if created, id, err := o.ReadOrCreate(&data, "Symbol", "Platform"); err == nil {
		if created {
			fmt.Println("New Insert an object. Id:", id)
		}
	}
	return data
}

func queryQuantity(name, platform string) models.Quantity {
	data := models.Quantity{Name: name, Platform: platform}
	if created, id, err := o.ReadOrCreate(&data, "Name", "Platform"); err == nil {
		if created {
			fmt.Println("New Insert an object. Id:", id)
		}
	}
	return data
}

func updateQuantity() {
	platform := "binance"
	bianAccount := api.BianAccountInfo()
	if bianAccount.Err != "" {
		log.Println("请求失败")
		return
	}
	bianBalances := bianAccount.Balances
	for _, balance := range bianBalances {
		free, _ := strconv.ParseFloat(balance.Free, 64)
		locked, _ := strconv.ParseFloat(balance.Locked, 64)
		if free == 0 && locked == 0 {
			continue
		}
		quantity := models.Quantity{Name: balance.Name, Platform: platform}
		if created, id, err := o.ReadOrCreate(&quantity, "Name", "Platform"); err == nil {
			if created {
				fmt.Println("Insert a Quantity object. Id:", id)
			}
		}
		quantity.Free = free
		quantity.Locked = locked
		if _, err := o.Update(&quantity); err == nil {
			fmt.Println("Update a Quantity object:", quantity.String())
		}
	}
}

func updateAccount() models.Account {
	bianAccount := api.BianAccountInfo()
	account := models.Account{Platform: "binance"}
	if bianAccount.Err != "" {
		log.Println("请求失败")
		return account
	}
	if created, id, err := o.ReadOrCreate(&account, "Platform"); err == nil {
		if created {
			fmt.Println("New Insert an object. Id:", id)
		}
	}
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

	rows, err := o.Update(&account)
	if err != nil {
		fmt.Printf("ERR: %v\n", err)
	}
	if rows == 0 {
		fmt.Printf("ERR: No record update\n")
	}
	return account
}

// Low price buy high price selling strategy
func RunLBHS() {
	// TODO All currencies with a policy status of 1 are participating in this strategy

	name := "XRP"
	symbol := name + "USDT"
	platform := "binance"
	sideBy := 0 // 1: sell, 2:buy

	lbhs := queryStrategyLBHS(symbol, platform)
	targetSellPrice := lbhs.TargetSellPrice

	// Calculate whether the balance is greater than the next cover price
	nextSpend := lbhs.LastSpend * lbhs.SpendCoefficient
	usdtQuantity := queryQuantity("USDT", platform)
	coinQuantity := queryQuantity(name, platform)
	notEnough := usdtQuantity.Free < nextSpend

	if usdtQuantity.Free < nextSpend && usdtQuantity.Free >= nextSpend/2 {
		targetSellPrice = lbhs.PositionAverage * lbhs.TargetProfitPoint / 2
	} else if usdtQuantity.Free < nextSpend/2 && usdtQuantity.Free >= nextSpend/3 {
		targetSellPrice = lbhs.PositionAverage * lbhs.TargetProfitPoint / 3
	} else if usdtQuantity.Free < nextSpend/3 {
		targetSellPrice = lbhs.PositionAverage * lbhs.TargetProfitPoint / 4
	}

	// Get the latest market price
	curPrice, _ := strconv.ParseFloat(api.BianTrade(symbol, 1)[0].Price, 64)
	log.Println(fmt.Sprintf("\nCurrent market price of %s: %.10f\nNext sale price: %.10f\nNext spend: %.10f\nNext buy price: %.10f", name, curPrice, targetSellPrice, nextSpend, lbhs.TargetBuyPrice))

	// Determine if it is higher than the specified value？Or is it lower than the specified value?
	if curPrice > targetSellPrice {
		sideBy = 1
		// Sell ​​target amount at current market price
		if notEnough {
			getUsdtBySell := lbhs.LastSpend
			targetSellQuantity := getUsdtBySell * curPrice
			if coinQuantity.Free > targetSellQuantity {
				// TODO Sell targetSellQuantity coins
			} else {
				getUsdtBySell = coinQuantity.Free * curPrice
				// TODO Sell all free coins
			}
			//update spend and actual cost
			lbhs.ActualCost = lbhs.ActualCost - getUsdtBySell
			lbhs.LastSpend = lbhs.LastSpend / lbhs.SpendCoefficient
		} else {
			// TODO The balance is sufficient, and the sales quota is adjusted according to the average market price.

		}
	} else if curPrice < lbhs.TargetBuyPrice {
		if usdtQuantity.Free > nextSpend {
			sideBy = 2
			// TODO Spend nextSpend amount to purchase

			lbhs.Spend = lbhs.Spend + nextSpend
			//update spend and actual cost
			lbhs.ActualCost = lbhs.ActualCost + nextSpend
			lbhs.LastSpend = nextSpend
		} else {
			// TODO Send a message to remind you to top up
		}
	}

	if sideBy == 0 {
		log.Println("There is no trade")
		return // continue
	}

	// Check the actual balance of the account after the transaction
	updateQuantity()
	//usdtQuantity = queryQuantity("USDT", platform)
	coinQuantity = queryQuantity(name, platform)

	// Update strategy, Reset parameters, Calculate the average price of the current position
	lbhs.Quantity = coinQuantity.Free + coinQuantity.Locked
	if sideBy == 1 { // sell
		lbhs.Spend = lbhs.Quantity * curPrice
		lbhs.PositionAverage = curPrice
	} else { // sideBy == 2 buy
		lbhs.PositionAverage = lbhs.Spend / lbhs.Quantity
	}
	lbhs.TargetSellPrice = lbhs.PositionAverage * lbhs.TargetProfitPoint
	lbhs.TargetBuyPrice = lbhs.PositionAverage * lbhs.TargetBuyPoint

	log.Println(lbhs.String())
	// update table
	//if _, err := o.Update(&lbhs); err == nil {
	//	fmt.Println("Update StrategyLowBuyHighSell:", lbhs.String())
	//}
}
