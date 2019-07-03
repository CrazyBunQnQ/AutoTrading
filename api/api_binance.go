package api

import (
	jsoniter "github.com/json-iterator/go"
	"io/ioutil"
	"log"
	"net/http"
)

func Depth(symbol string, limit int32) string {
	resp, err := http.Get(fullApi("/api/v1/depth"))
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

// Exchange information
func ExchangeInfo() string {
	resp, err := http.Get(fullApi("/api/v1/exchangeInfo"))
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

// Ping Test server connectivity
func Ping() bool {
	resp, err := http.Get(fullApi("/api/v1/ping"))
	if err != nil {
		log.Println(err)
		return false
	}
	defer resp.Body.Close()
	return true
}

// Time Get server time
func Time() int64 {
	resp, err := http.Get(fullApi("/api/v1/time"))
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

func fullApi(api string) string {
	return BianConf.BaseUrl + api
}
