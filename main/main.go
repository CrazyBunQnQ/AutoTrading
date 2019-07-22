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

func queryStrategyLBHS(symbol, platform string) models.StrategyLowBuyHighSell {
	data := models.StrategyLowBuyHighSell{Symbol: symbol, Platform: platform}
	if created, id, err := o.ReadOrCreate(&data, "Symbol", "Platform"); err == nil {
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
		quantity := models.Quantity{Symbol: balance.Symbol, Platform: platform}
		if created, id, err := o.ReadOrCreate(&quantity, "Symbol", "Platform"); err == nil {
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
		switch balance.Symbol {
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

	id, err := o.Update(&account)
	a := models.Account{Id: id}
	err = o.Read(&a)
	fmt.Printf("ERR: %v\n", err)
	return a
}
