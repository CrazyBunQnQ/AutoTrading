package api

import (
	"io/ioutil"
	"log"
	"net/http"
)

func httpGetJsonStr(fullUrl string) string {
	resp, err := http.Get(fullUrl)
	if err != nil {
		log.Println("api resp error: " + fullUrl)
		log.Println(err)
		return ""
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("api get body erroi: " + fullUrl)
		log.Fatal(err)
		return ""
	}
	return string(body)
}
