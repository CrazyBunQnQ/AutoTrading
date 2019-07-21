package main

import (
	"AutoTrading/api"
	"AutoTrading/utils"
	"io"
	"net/http"
	"strconv"
	"time"
)

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
