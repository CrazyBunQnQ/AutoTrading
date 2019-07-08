package api

import (
	jsoniter "github.com/json-iterator/go"
	"io/ioutil"
	"log"
	"net/http"
)

func httpGetJsonStr(fullUrl string) jsoniter.Any {
	resp, err := http.Get(fullUrl)
	if err != nil {
		log.Println("api resp error: " + fullUrl)
		log.Println(err)
		return nil
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("api get body erroi: " + fullUrl)
		log.Fatal(err)
		return nil
	}
	return jsoniter.Get(body)
}
