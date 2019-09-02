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
	"time"
)

// ************************* public API **********************

// strSymbol: Transaction pair, btcusdt, bccbtc......
// strPeriod: KLine type, 1min, 5min, 15min......
// nSize: Get quantity, [1-2000]
// return: KLineReturn  Object
// Opening time
// Opening price
// Highest price
// Lowest price
// Closing price (the current price is not the current K line)
// Volume
// Closing time
// Turnover
// Number of transactions
// Active buying volume
// Active buying turnover
// Please ignore this parameter
func BianKLine(symbol string, interval models.BianInterval, limit int, startTime, endTime int64) []interface{} {
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

	jsonReturn, err := utils.HttpGetRequest(strUrl, mapParams)
	if err != "" {
		errJson := "{\"err\": \"" + err + "\"}"
		json.Unmarshal([]byte(errJson), &result)
	} else {
		json.Unmarshal([]byte(jsonReturn), &result)
	}

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

	jsonReturn, err := utils.HttpGetRequest(strUrl, mapParams)
	if err != "" {
		errJson := "{\"err\": \"" + err + "\"}"
		json.Unmarshal([]byte(errJson), &result)
	} else {
		json.Unmarshal([]byte(jsonReturn), &result)
	}

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

	jsonReturn, err := utils.HttpGetRequest(strUrl, mapParams)
	if err != "" {
		errJson := "{\"err\": \"" + err + "\"}"
		json.Unmarshal([]byte(errJson), &result)
	} else {
		json.Unmarshal([]byte(jsonReturn), &result)
	}

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

	jsonReturn, err := utils.HttpGetRequest(strUrl, mapParams)
	if err != "" {
		errJson := "{\"err\": \"" + err + "\"}"
		json.Unmarshal([]byte(errJson), &result)
	} else {
		json.Unmarshal([]byte(jsonReturn), &result)
	}

	return result
}

func BianAvgPrice(symbol string) models.BianAvgPrice {
	result := models.BianAvgPrice{}

	mapParams := make(map[string]string)
	mapParams["symbol"] = symbol

	strRequestUrl := "/api/v3/avgPrice"
	strUrl := config.BianConf.BaseUrl + strRequestUrl

	jsonReturn, err := utils.HttpGetRequest(strUrl, mapParams)
	if err != "" {
		errJson := "{\"err\": \"" + err + "\"}"
		json.Unmarshal([]byte(errJson), &result)
	} else {
		json.Unmarshal([]byte(jsonReturn), &result)
	}

	return result
}

func BianTicker24(symbol string) models.BianTicker24 {
	result := models.BianTicker24{}

	mapParams := make(map[string]string)
	mapParams["symbol"] = symbol

	strRequestUrl := "/api/v1/ticker/24hr"
	strUrl := config.BianConf.BaseUrl + strRequestUrl

	jsonReturn, err := utils.HttpGetRequest(strUrl, mapParams)
	if err != "" {
		errJson := "{\"err\": \"" + err + "\"}"
		json.Unmarshal([]byte(errJson), &result)
	} else {
		json.Unmarshal([]byte(jsonReturn), &result)
	}

	return result
}

func BianTicker24All() []models.BianTicker24 {
	result := []models.BianTicker24{}

	mapParams := make(map[string]string)

	strRequestUrl := "/api/v1/ticker/24hr"
	strUrl := config.BianConf.BaseUrl + strRequestUrl

	jsonReturn, err := utils.HttpGetRequest(strUrl, mapParams)
	if err != "" {
		errJson := "{\"err\": \"" + err + "\"}"
		json.Unmarshal([]byte(errJson), &result)
	} else {
		json.Unmarshal([]byte(jsonReturn), &result)
	}

	return result
}

func BianLastPrice(symbol string) models.BianLastPrice {
	result := models.BianLastPrice{}

	mapParams := make(map[string]string)
	mapParams["symbol"] = symbol

	strRequestUrl := "/api/v3/ticker/price"
	strUrl := config.BianConf.BaseUrl + strRequestUrl

	jsonReturn, err := utils.HttpGetRequest(strUrl, mapParams)
	if err != "" {
		errJson := "{\"err\": \"" + err + "\"}"
		json.Unmarshal([]byte(errJson), &result)
	} else {
		json.Unmarshal([]byte(jsonReturn), &result)
	}

	return result
}

func BianLastAllPrice() []models.BianLastPrice {
	result := []models.BianLastPrice{}

	mapParams := make(map[string]string)

	strRequestUrl := "/api/v3/ticker/price"
	strUrl := config.BianConf.BaseUrl + strRequestUrl

	jsonReturn, err := utils.HttpGetRequest(strUrl, mapParams)
	if err != "" {
		errJson := "{\"err\": \"" + err + "\"}"
		json.Unmarshal([]byte(errJson), &result)
	} else {
		json.Unmarshal([]byte(jsonReturn), &result)
	}

	return result
}

func BianBestTicker(symbol string) models.BianBestTicker {
	result := models.BianBestTicker{}

	mapParams := make(map[string]string)
	mapParams["symbol"] = symbol

	strRequestUrl := "/api/v3/ticker/bookTicker"
	strUrl := config.BianConf.BaseUrl + strRequestUrl

	jsonReturn, err := utils.HttpGetRequest(strUrl, mapParams)
	if err != "" {
		errJson := "{\"err\": \"" + err + "\"}"
		json.Unmarshal([]byte(errJson), &result)
	} else {
		json.Unmarshal([]byte(jsonReturn), &result)
	}

	return result
}

func BianAllBestTicker() []models.BianBestTicker {
	result := []models.BianBestTicker{}

	mapParams := make(map[string]string)

	strRequestUrl := "/api/v3/ticker/bookTicker"
	strUrl := config.BianConf.BaseUrl + strRequestUrl

	jsonReturn, err := utils.HttpGetRequest(strUrl, mapParams)
	if err != "" {
		errJson := "{\"err\": \"" + err + "\"}"
		json.Unmarshal([]byte(errJson), &result)
	} else {
		json.Unmarshal([]byte(jsonReturn), &result)
	}

	return result
}

// ************************* Account API **********************
// Create a limit order
func BianOrderByLimit(symbol string, side models.BianOrderSide, timeInForce models.BianTimeInForce, quantity, price, icebergQty float64) models.BianFastestOrder {
	result := models.BianFastestOrder{}

	mapParams := make(map[string]string)
	mapParams["type"] = string(models.TypeLimit)
	// fastest response type
	mapParams["newOrderRespType"] = string(models.AckResponse)
	mapParams["symbol"] = symbol
	mapParams["side"] = string(side)
	mapParams["timeInForce"] = string(timeInForce)
	mapParams["quantity"] = strconv.FormatFloat(quantity, 'f', -1, 64)
	mapParams["price"] = strconv.FormatFloat(price, 'f', -1, 64)
	mapParams["timestamp"] = strconv.FormatInt(utils.UnixMillis(time.Now()), 10)
	if icebergQty != 0 {
		mapParams["icebergQty"] = strconv.FormatFloat(icebergQty, 'f', -1, 64)
	}

	strRequestUrl := "/api/v3/order"
	strUrl := config.BianConf.BaseUrl + strRequestUrl

	jsonReturn, err := utils.BianPostRequest(strUrl, mapParams, true)
	if err != "" {
		errJson := "{\"err\": \"" + err + "\"}"
		json.Unmarshal([]byte(errJson), &result)
	} else {
		json.Unmarshal([]byte(jsonReturn), &result)
	}
	json.Unmarshal([]byte(jsonReturn), &result)

	return result
}

func BianOrderByMarket(symbol string, side models.BianOrderSide, quantity, icebergQty float64) models.BianFastestOrder {
	result := models.BianFastestOrder{}

	mapParams := make(map[string]string)
	mapParams["type"] = string(models.TypeMarket)
	// fastest response type
	mapParams["newOrderRespType"] = string(models.AckResponse)
	mapParams["symbol"] = symbol
	mapParams["side"] = string(side)
	mapParams["quantity"] = strconv.FormatFloat(quantity, 'f', -1, 64)
	mapParams["timestamp"] = strconv.FormatInt(utils.UnixMillis(time.Now()), 10)
	if icebergQty != 0 {
		mapParams["icebergQty"] = strconv.FormatFloat(icebergQty, 'f', -1, 64)
	}

	strRequestUrl := "/api/v3/order"
	strUrl := config.BianConf.BaseUrl + strRequestUrl

	jsonReturn, err := utils.BianPostRequest(strUrl, mapParams, true)
	if err != "" {
		errJson := "{\"err\": \"" + err + "\"}"
		json.Unmarshal([]byte(errJson), &result)
	} else {
		json.Unmarshal([]byte(jsonReturn), &result)
	}
	json.Unmarshal([]byte(jsonReturn), &result)

	return result
}

func BianOrderQuery(symbol, origClientOrderId string, orderId int64) models.BianOrderStatus {
	result := models.BianOrderStatus{}

	mapParams := make(map[string]string)
	mapParams["symbol"] = symbol
	mapParams["timestamp"] = strconv.FormatInt(utils.UnixMillis(time.Now()), 10)
	if origClientOrderId != "" {
		mapParams["origClientOrderId"] = origClientOrderId
	}
	if orderId != 0 {
		mapParams["orderId"] = strconv.FormatInt(orderId, 10)
	}

	strRequestUrl := "/api/v3/order"
	strUrl := config.BianConf.BaseUrl + strRequestUrl

	jsonReturn, err := utils.BianGetRequest(strUrl, mapParams, true)
	if err != "" {
		errJson := "{\"err\": \"" + err + "\"}"
		json.Unmarshal([]byte(errJson), &result)
	} else {
		json.Unmarshal([]byte(jsonReturn), &result)
	}

	return result
}

func BianOpenOrder(symbol string) []models.BianOrderStatus {
	result := []models.BianOrderStatus{}

	mapParams := make(map[string]string)
	mapParams["timestamp"] = strconv.FormatInt(utils.UnixMillis(time.Now()), 10)
	mapParams["recvWindow"] = strconv.FormatInt(recvWindow(5*time.Second), 10)
	if symbol != "" {
		mapParams["symbol"] = symbol
	}

	strRequestUrl := "/api/v3/openOrders"
	strUrl := config.BianConf.BaseUrl + strRequestUrl

	jsonReturn, err := utils.BianGetRequest(strUrl, mapParams, true)
	if err != "" {
		errJson := "{\"err\": \"" + err + "\"}"
		json.Unmarshal([]byte(errJson), &result)
	} else {
		json.Unmarshal([]byte(jsonReturn), &result)
	}

	return result
}

func BianOrderDelete(symbol, origClientOrderId string, orderId int64) models.BianOrderStatus {
	result := models.BianOrderStatus{}

	mapParams := make(map[string]string)
	mapParams["symbol"] = symbol
	mapParams["timestamp"] = strconv.FormatInt(utils.UnixMillis(time.Now()), 10)
	if origClientOrderId != "" {
		mapParams["origClientOrderId"] = origClientOrderId
	}
	if orderId != 0 {
		mapParams["orderId"] = strconv.FormatInt(orderId, 10)
	}

	strRequestUrl := "/api/v3/order"
	strUrl := config.BianConf.BaseUrl + strRequestUrl

	jsonReturn, err := utils.BianDeleteRequest(strUrl, mapParams, true)
	if err != "" {
		errJson := "{\"err\": \"" + err + "\"}"
		json.Unmarshal([]byte(errJson), &result)
	} else {
		json.Unmarshal([]byte(jsonReturn), &result)
	}

	return result
}

func BianAccountInfo() models.BianAccount {
	result := models.BianAccount{}

	mapParams := make(map[string]string)
	mapParams["timestamp"] = strconv.FormatInt(utils.UnixMillis(time.Now()), 10)
	mapParams["recvWindow"] = strconv.FormatInt(recvWindow(5*time.Second), 10)

	strRequestUrl := "/api/v3/account"
	strUrl := config.BianConf.BaseUrl + strRequestUrl

	jsonReturn, err := utils.BianGetRequest(strUrl, mapParams, true)
	if err != "" {
		errJson := "{\"err\": \"" + err + "\"}"
		json.Unmarshal([]byte(errJson), &result)
	} else {
		json.Unmarshal([]byte(jsonReturn), &result)
	}

	return result
}

func recvWindow(d time.Duration) int64 {
	return int64(d) / int64(time.Millisecond)
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
