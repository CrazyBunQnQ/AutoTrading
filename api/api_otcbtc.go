package api

import (
	"AutoTrading/config"
)

func OtcbtcTrades(symbol string, limit string) string {
	return otcbtcApiJsonResult("/api/v2/trades?market=" + symbol + "&limit=" + limit)
}

func OtcbtcDepth(symbol string, limit string) string {
	return otcbtcApiJsonResult("/api/v2/depth?market=" + symbol + "&limit=" + limit)
}

func OtcbtcTickers(symbol string) string {
	return otcbtcApiJsonResult("/api/v2/tickers/" + symbol)
}

func otcbtcApiJsonResult(fullApi string) string {
	return httpGetJsonStr(fullOtcbtcApi(fullApi))
}

func fullOtcbtcApi(api string) string {
	return config.OtcbtcConf.BaseUrl + api
}
