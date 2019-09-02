package api

import (
	"AutoTrading/config"
	"bytes"
	"encoding/json"
	"fmt"
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

func IftttNotice(title, text, link string) {
	//拼接url
	u := fmt.Sprintf(config.Ifttt.WebhooksUrl, config.Ifttt.EventName, config.Ifttt.Key)
	log.Printf("ifttt url=%s", u)

	//post
	values := map[string]string{"value1": title, "value2": text, "value3": link}
	jsonStr, _ := json.Marshal(values)
	resp, err := http.Post(u, "application/json", bytes.NewBuffer(jsonStr))
	if err != nil {
		log.Fatalf("http.Post error:%v", err)
	}
	d, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("ioutil.ReadAll error:%v", err)
	}
	log.Printf("resp:\n%s", d)
}
