package api

import (
	jsoniter "github.com/json-iterator/go"
	"io/ioutil"
	"log"
	"net/http"
)

func OtcbtcDepth(symbol string, limit string) string {
	resp, err := http.Get(fullOtcbtcApi("/api/v2/depth?market=" + symbol + "&limit=" + limit))
	if err != nil {
		log.Println(err)
		return ""
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
		return ""
	}
	return jsoniter.Get(body).ToString()
}

func OtcbtcTickers(symbol string) string {
	resp, err := http.Get(fullOtcbtcApi("/api/v2/tickers") + symbol)
	if err != nil {
		log.Println(err)
		return ""
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
		return ""
	}
	return jsoniter.Get(body).ToString()
}

func fullOtcbtcApi(api string) string {
	return OtcbtcConf.BaseUrl + api
}
