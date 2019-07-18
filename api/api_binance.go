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
func BianKLine(symbol string, interval models.Interval, limit int, startTime, endTime int64) []interface{} {
	var result []interface{}

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

	jsonReturn := utils.HttpGetRequest(strUrl, mapParams)
	json.Unmarshal([]byte(jsonReturn), &result)

	return result
}

// Get transaction depth information
// symbol: Transaction pair, btcusdt, bccbtc......
// limit: Depth type, Default 100; Maximum 1000. Optional values: [5, 10, 20, 50, 100, 500, 1000]
// return: HuobiDepthReturn Object
func BianDepth(symbol string, limit int) models.BianDepth {
	result := models.BianDepth{}

	mapParams := make(map[string]string)
	mapParams["symbol"] = symbol
	if limit != 0 {
		mapParams["limit"] = strconv.Itoa(limit)
	}

	strRequestUrl := "/api/v1/depth"
	strUrl := config.BianConf.BaseUrl + strRequestUrl

	jsonReturn := utils.HttpGetRequest(strUrl, mapParams)
	json.Unmarshal([]byte(jsonReturn), &result)

	return result
}

// Get recent transaction history in bulk
// symbol: Transaction pair, btcusdt, bccbtc......
// limit: Get the number of transaction records, Default 500; max 1000.
// return: TradeReturn Object
func BianTrade(symbol string, limit int) []models.BianTrade {
	result := []models.BianTrade{}

	mapParams := make(map[string]string)
	mapParams["symbol"] = symbol
	if limit != 0 {
		mapParams["limit"] = strconv.Itoa(limit)
	}

	strRequestUrl := "/api/v1/trades"
	strUrl := config.BianConf.BaseUrl + strRequestUrl

	jsonReturn := utils.HttpGetRequest(strUrl, mapParams)
	json.Unmarshal([]byte(jsonReturn), &result)

	return result
}

// Get recent transaction history in bulk
// symbol: Transaction pair, btcusdt, bccbtc......
// limit: Get the number of transaction records, Default 500; max 1000.
// return: TradeReturn Object
func BianAggTrade(symbol string, limit, fromId int, startTime, endTime int64) []models.BianAggTrade {
	result := []models.BianAggTrade{}

	mapParams := make(map[string]string)
	mapParams["symbol"] = symbol
	if limit != 0 {
		mapParams["limit"] = strconv.Itoa(limit)
	}
	if fromId != 0 {
		mapParams["fromId"] = strconv.Itoa(fromId)
	}
	if startTime != 0 {
		mapParams["startTime"] = strconv.FormatInt(startTime, 64)
	}
	if endTime != 0 {
		mapParams["endTime"] = strconv.FormatInt(endTime, 64)
	}

	strRequestUrl := "/api/v1/aggTrades"
	strUrl := config.BianConf.BaseUrl + strRequestUrl

	jsonReturn := utils.HttpGetRequest(strUrl, mapParams)
	json.Unmarshal([]byte(jsonReturn), &result)

	return result
}

func BianAvgPrice(symbol string) models.BianAvgPrice {
	result := models.BianAvgPrice{}

	mapParams := make(map[string]string)
	mapParams["symbol"] = symbol

	strRequestUrl := "/api/v3/avgPrice"
	strUrl := config.BianConf.BaseUrl + strRequestUrl

	jsonReturn := utils.HttpGetRequest(strUrl, mapParams)
	json.Unmarshal([]byte(jsonReturn), &result)

	return result
}

func BianTicker24(symbol string) models.BianTicker24 {
	result := models.BianTicker24{}

	mapParams := make(map[string]string)
	mapParams["symbol"] = symbol

	strRequestUrl := "/api/v1/ticker/24hr"
	strUrl := config.BianConf.BaseUrl + strRequestUrl

	jsonReturn := utils.HttpGetRequest(strUrl, mapParams)
	json.Unmarshal([]byte(jsonReturn), &result)

	return result
}

func BianTicker24All() []models.BianTicker24 {
	result := []models.BianTicker24{}

	mapParams := make(map[string]string)

	strRequestUrl := "/api/v1/ticker/24hr"
	strUrl := config.BianConf.BaseUrl + strRequestUrl

	jsonReturn := utils.HttpGetRequest(strUrl, mapParams)
	json.Unmarshal([]byte(jsonReturn), &result)

	return result
}

func BianLastPrice(symbol string) models.BianLastPrice {
	result := models.BianLastPrice{}

	mapParams := make(map[string]string)
	mapParams["symbol"] = symbol

	strRequestUrl := "/api/v3/ticker/price"
	strUrl := config.BianConf.BaseUrl + strRequestUrl

	jsonReturn := utils.HttpGetRequest(strUrl, mapParams)
	json.Unmarshal([]byte(jsonReturn), &result)

	return result
}

func BianLastAllPrice() []models.BianLastPrice {
	result := []models.BianLastPrice{}

	mapParams := make(map[string]string)

	strRequestUrl := "/api/v3/ticker/price"
	strUrl := config.BianConf.BaseUrl + strRequestUrl

	jsonReturn := utils.HttpGetRequest(strUrl, mapParams)
	json.Unmarshal([]byte(jsonReturn), &result)

	return result
}

// ************************************************************

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
