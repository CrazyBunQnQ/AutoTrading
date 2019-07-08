package api

import (
	jsoniter "github.com/json-iterator/go"
	"io/ioutil"
	"log"
	"net/http"
)

func BianTrades(symbol string, limit string) jsoniter.Any {
	return bianApiJsonResult("/api/v1/trades?symbol=" + symbol + "&limit=" + limit)
}

func BianDepth(symbol string, limit string) jsoniter.Any {
	return bianApiJsonResult("/api/v1/depth?symbol=" + symbol + "&limit=" + limit)
}

// Exchange information
func ExchangeInfo() jsoniter.Any {
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
	return jsoniter.Get(body, "serverTime").ToInt64()
}

func bianApiJsonResult(fullApi string) jsoniter.Any {
	return httpGetJsonStr(fullBianApi(fullApi))
}

func fullBianApi(api string) string {
	return BianConf.BaseUrl + api
}
