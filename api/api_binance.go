package api

import (
	jsoniter "github.com/json-iterator/go"
	"io/ioutil"
	"log"
	"net/http"
)

// Ping Test server connectivity
func Ping() bool {
	resp, err := http.Get("https://api.binance.com/api/v1/ping")
	if err != nil {
		log.Println(err)
		return false
	}
	defer resp.Body.Close()
	return true
}

// Time Get server time
func Time() int64 {
	resp, err := http.Get("https://api.binance.com/api/v1/time")
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
