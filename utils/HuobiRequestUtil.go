package utils

import (
	"AutoTrading/config"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/url"
	"sort"
	"time"
)

// 进行签名后的HTTP GET请求, 参考官方Python Demo写的
// mapParams: map类型的请求参数, key:value
// strRequest: API路由路径
// return: 请求结果
func ApiKeyGet(mapParams map[string]string, strRequestPath string) (string, string) {
	strMethod := "GET"
	timestamp := time.Now().UTC().Format("2006-01-02T15:04:05")

	mapParams["AccessKeyId"] = config.HuoBiConf.AccessKeyPublic
	mapParams["SignatureMethod"] = "HmacSHA256"
	mapParams["SignatureVersion"] = "2"
	mapParams["Timestamp"] = timestamp

	hostName := config.HuoBiConf.HostName
	mapParams["Signature"] = CreateSign(mapParams, strMethod, hostName, strRequestPath, config.HuoBiConf.SecretKeyPublic)

	if config.HuoBiConf.EnablePrivateSignature == true {
		privateSignature, err := CreatePrivateSignByJWT(mapParams["Signature"])
		if nil == err {
			mapParams["PrivateSignature"] = privateSignature
		} else {
			fmt.Println("signed error: ", err)
		}
	}

	strUrl := config.HuoBiConf.TradeUrl + strRequestPath

	return HttpGetRequest(strUrl, MapValueEncodeURI(mapParams))
}

// 进行签名后的HTTP POST请求, 参考官方Python Demo写的
// mapParams: map类型的请求参数, key:value
// strRequest: API路由路径
// return: 请求结果
func ApiKeyPost(mapParams map[string]string, strRequestPath string) string {
	strMethod := "POST"
	timestamp := time.Now().UTC().Format("2006-01-02T15:04:05")

	mapParams2Sign := make(map[string]string)
	mapParams2Sign["AccessKeyId"] = config.HuoBiConf.AccessKeyPublic
	mapParams2Sign["SignatureMethod"] = "HmacSHA256"
	mapParams2Sign["SignatureVersion"] = "2"
	mapParams2Sign["Timestamp"] = timestamp

	hostName := config.HuoBiConf.HostName

	mapParams2Sign["Signature"] = CreateSign(mapParams2Sign, strMethod, hostName, strRequestPath, config.HuoBiConf.SecretKeyPublic)

	if config.HuoBiConf.EnablePrivateSignature == true {
		privateSignature, err := CreatePrivateSignByJWT(mapParams2Sign["Signature"])

		if nil == err {
			mapParams2Sign["PrivateSignature"] = privateSignature
		} else {
			fmt.Println("signed error:", err)
		}
	}

	strUrl := config.HuoBiConf.TradeUrl + strRequestPath + "?" + Map2UrlQuery(MapValueEncodeURI(mapParams2Sign))

	return HttpPostRequest(strUrl, mapParams)
}

// 构造签名
// mapParams: 送进来参与签名的参数, Map类型
// strMethod: 请求的方法 GET, POST......
// strHostUrl: 请求的主机
// strRequestPath: 请求的路由路径
// strSecretKey: 进行签名的密钥
func CreateSign(mapParams map[string]string, strMethod, strHostUrl, strRequestPath, strSecretKey string) string {
	// 参数处理, 按API要求, 参数名应按ASCII码进行排序(使用UTF-8编码, 其进行URI编码, 16进制字符必须大写)
	mapCloned := make(map[string]string)
	for key, value := range mapParams {
		mapCloned[key] = url.QueryEscape(value)
	}

	strParams := Map2UrlQueryBySort(mapCloned)

	strPayload := strMethod + "\n" + strHostUrl + "\n" + strRequestPath + "\n" + strParams
	return ComputeHmac256(strPayload, strSecretKey)
}

func CreatePrivateSignByJWT(sign string) (string, error) {
	return SignByJWT(config.HuoBiConf.PrivateKeyPrime256, sign)
}

// 对Map按着ASCII码进行排序
// mapValue: 需要进行排序的map
// return: 排序后的map
func MapSortByKey(mapValue map[string]string) map[string]string {
	var keys []string
	for key := range mapValue {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	mapReturn := make(map[string]string)
	for _, key := range keys {
		mapReturn[key] = mapValue[key]
	}

	return mapReturn
}

// 对Map的值进行URI编码
// mapParams: 需要进行URI编码的map
// return: 编码后的map
func MapValueEncodeURI(mapValue map[string]string) map[string]string {
	for key, value := range mapValue {
		valueEncodeURI := url.QueryEscape(value)
		mapValue[key] = valueEncodeURI
	}

	return mapValue
}

// 将map格式的请求参数转换为字符串格式的,并按照Map的key升序排列
// mapParams: map格式的参数键值对
// return: 查询字符串
func Map2UrlQueryBySort(mapParams map[string]string) string {
	var keys []string
	for key := range mapParams {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	var strParams string
	for _, key := range keys {
		strParams += key + "=" + mapParams[key] + "&"
	}

	if 0 < len(strParams) {
		strParams = string([]rune(strParams)[:len(strParams)-1])
	}

	return strParams
}

// HMAC SHA256加密
// strMessage: 需要加密的信息
// strSecret: 密钥
// return: BASE64编码的密文
func ComputeHmac256(strMessage string, strSecret string) string {
	key := []byte(strSecret)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(strMessage))

	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}
