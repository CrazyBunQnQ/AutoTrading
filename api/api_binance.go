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

func BianTrades(symbol string, limit string) string {
	return bianApiJsonResult("/api/v1/trades?symbol=" + symbol + "&limit=" + limit)
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
