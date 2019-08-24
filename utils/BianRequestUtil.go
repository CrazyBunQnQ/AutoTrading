package utils

import (
	"AutoTrading/config"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"io/ioutil"
	"net/http"
)

func BianGetRequest(strUrl string, params map[string]string, sign bool) (string, string) {
	return BianRequest("GET", strUrl, params, sign)
}

func BianPostRequest(strUrl string, params map[string]string, sign bool) (string, string) {
	return BianRequest("POST", strUrl, params, sign)
}

func BianDeleteRequest(strUrl string, params map[string]string, sign bool) (string, string) {
	return BianRequest("DELETE", strUrl, params, sign)
}

func BianRequest(requestType, strUrl string, params map[string]string, sign bool) (string, string) {
	httpClient := &http.Client{}

	request, err := http.NewRequest(requestType, strUrl, nil)
	if nil != err {
		return "", "创建 " + requestType + " 请求失败: " + err.Error()
	}

	q := request.URL.Query()
	for key, val := range params {
		q.Add(key, val)
	}
	if sign {
		q.Add("signature", BianSign([]byte(config.BianConf.SecretKeyPrivate), []byte(q.Encode())))
	}

	request.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/74.0.3729.169 Safari/537.36")
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Accept-Language", "zh-cn")
	request.Header.Add("X-MBX-APIKEY", config.BianConf.ApiKeyPrivate)

	request.URL.RawQuery = q.Encode()

	response, err := httpClient.Do(request)
	defer response.Body.Close()
	if nil != err {
		return "", requestType + " 请求失败: " + err.Error()
	}

	body, err := ioutil.ReadAll(response.Body)
	if nil != err {
		return "", "读取相应内容失败: " + err.Error()
	}

	return string(body), ""
}

// Sign signs provided payload and returns encoded string sum.
func BianSign(keyByte, queryEncodeByte []byte) string {
	mac := hmac.New(sha256.New, keyByte)
	mac.Write(queryEncodeByte)
	return hex.EncodeToString(mac.Sum(nil))
}
