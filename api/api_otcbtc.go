package api

import (
	jsoniter "github.com/json-iterator/go"
)

func OtcbtcTrades(symbol string, limit string) jsoniter.Any {
	return otcbtcApiJsonResult("/api/v2/trades?market=" + symbol + "&limit=" + limit)
}

func OtcbtcDepth(symbol string, limit string) jsoniter.Any {
	return otcbtcApiJsonResult("/api/v2/depth?market=" + symbol + "&limit=" + limit)
}

func OtcbtcTickers(symbol string) jsoniter.Any {
	return otcbtcApiJsonResult("/api/v2/tickers/" + symbol)
}

func otcbtcApiJsonResult(fullApi string) jsoniter.Any {
	return httpGetJsonStr(fullOtcbtcApi(fullApi))
}

func fullOtcbtcApi(api string) string {
	return OtcbtcConf.BaseUrl + api
}
