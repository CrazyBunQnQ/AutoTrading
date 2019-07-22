package main

import (
	"AutoTrading/api"
	"AutoTrading/models"
	"AutoTrading/utils"
	"fmt"
	"github.com/astaxie/beego/orm"
	"io"
	"net/http"
	"strconv"
	"time"
)

var o orm.Ormer

func init() {
	orm.RegisterDataBase("default", "mysql", "root:zy26T$b7V8i3g4mW@tcp(127.0.0.1:3306)/autotrade?charset=utf8mb4", 30)
	orm.RegisterModel(new(models.StrategyLowBuyHighSell))
	orm.RegisterModel(new(models.Account))
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

func updateAccount() models.Account {
	bianAccount := api.BianAccountInfo()
	account := models.Account{Platform: "binance"}
	bianBalances := bianAccount.Balances
	for _, balance := range bianBalances {
		switch balance.Symbol {
		case "USDT":
			num, _ := strconv.ParseFloat(balance.Free, 64)
			account.Usdt = num
		case "BTC":
			num, _ := strconv.ParseFloat(balance.Free, 64)
			account.Btc = num
		case "ETH":
			num, _ := strconv.ParseFloat(balance.Free, 64)
			account.Eth = num
		case "BNB":
			num, _ := strconv.ParseFloat(balance.Free, 64)
			account.Bnb = num
		case "EOS":
			num, _ := strconv.ParseFloat(balance.Free, 64)
			account.Eos = num
		case "XRP":
			num, _ := strconv.ParseFloat(balance.Free, 64)
			account.Xrp = num
		}
	}

	id, err := o.Insert(&account)
	a := models.Account{Id: id}
	err = o.Read(&a)
	fmt.Printf("ERR: %v\n", err)
	return a
}
