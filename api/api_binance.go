package api

import (
	"AutoTrading/config"
	"AutoTrading/models"
	"AutoTrading/utils"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

// strSymbol: Transaction pair, btcusdt, bccbtc......
// strPeriod: KLine type, 1min, 5min, 15min......
// nSize: Get quantity, [1-2000]
// return: KLineReturn  Object
func GetBianKLine(symbol string, interval models.Interval, limit int, startTime, endTime int64) []interface{} {
	var bianKLines []interface{}

	mapParams := make(map[string]string)
	mapParams["symbol"] = symbol
	mapParams["interval"] = string(interval)
	if limit != 0 {
		mapParams["limit"] = strconv.Itoa(limit)
	}
	if startTime != 0 {
		mapParams["startTime"] = strconv.FormatInt(startTime, 10)
	}
	if endTime != 0 {
		mapParams["endTime"] = strconv.FormatInt(endTime, 10)
	}

	strRequestUrl := "/api/v1/klines"
	strUrl := config.BianConf.BaseUrl + strRequestUrl

	jsonKLineReturn := utils.HttpGetRequest(strUrl, mapParams)
	json.Unmarshal([]byte(jsonKLineReturn), &bianKLines)

	return bianKLines
}

func BianDepth(strSymbol string, limit int) models.BianDepth {
	marketDepth := models.BianDepth{}

	mapParams := make(map[string]string)
	mapParams["symbol"] = strSymbol
	mapParams["limit"] = strconv.Itoa(limit)

	strRequestUrl := "/api/v1/depth"
	strUrl := config.BianConf.BaseUrl + strRequestUrl

	jsonMarketDepthReturn := utils.HttpGetRequest(strUrl, mapParams)
	json.Unmarshal([]byte(jsonMarketDepthReturn), &marketDepth)

	return marketDepth
}

func BianTrades(symbol string, limit string) string {
	return bianApiJsonResult("/api/v1/trades?symbol=" + symbol + "&limit=" + limit)
}

// Exchange information
func ExchangeInfo() string {
	return bianApiJsonResult("/api/v1/exchangeInfo")
}

// Ping Test server connectivity
func Ping() bool {
	resp, err := http.Get(fullBianApi("/api/v1/ping"))
	if err != nil {
		log.Println(err)
		return false
	}
	defer resp.Body.Close()
	return true
}

// Time Get server time
func Time() int64 {
	resp, err := http.Get(fullBianApi("/api/v1/time"))
	if err != nil {
		log.Println(err)
		return 0
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
		return 0
	}
	log.Println(string(body))
	return 0
	//return json.Get(body, "serverTime").ToInt64()
}

func bianApiJsonResult(fullApi string) string {
	return httpGetJsonStr(fullBianApi(fullApi))
}

func fullBianApi(api string) string {
	return config.BianConf.BaseUrl + api
}
