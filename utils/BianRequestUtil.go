package utils

import (
	"AutoTrading/config"
	"io/ioutil"
	"net/http"
	"time"
)

func BianGetRequest(strUrl string, params map[string]string, sign bool) string {
	httpClient := &http.Client{}

	request, err := http.NewRequest("POST", strUrl, nil)
	if nil != err {
		return err.Error()
	}

	q := request.URL.Query()
	for key, val := range params {
		q.Add(key, val)
	}

	//hmacSigner := &models.HmacSigner{
	//	Key: []byte(config.BianConf.SecretKeyPrivate),
	//}
	//q.Add("signature", hmacSigner.Sign([]byte(q.Encode())))

	request.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/74.0.3729.169 Safari/537.36")
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Accept-Language", "zh-cn")
	request.Header.Add("X-MBX-APIKEY", config.BianConf.ApiKeyPrivate)

	request.URL.RawQuery = q.Encode()

	response, err := httpClient.Do(request)
	defer response.Body.Close()
	if nil != err {
		return err.Error()
	}

	body, err := ioutil.ReadAll(response.Body)
	if nil != err {
		return err.Error()
	}

	return string(body)
}

func BianSign() {

}

func UnixMillis(t time.Time) int64 {
	return t.UnixNano() / int64(time.Millisecond)
}
