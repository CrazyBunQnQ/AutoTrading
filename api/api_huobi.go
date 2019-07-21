package api

import (
	"AutoTrading/config"
	"AutoTrading/models"
	"AutoTrading/utils"
	"encoding/json"
	"fmt"
	"strconv"
)

//------------------------------------------------------------------------------------------
// Trade API

// strSymbol: Transaction pair, btcusdt, bccbtc......
// strPeriod: KLine type, 1min, 5min, 15min......
// nSize: Get quantity, [1-2000]
// return: KLineReturn  Object
func GetHuobiKLine(strSymbol, strPeriod string, nSize int) models.KLineReturn {
	kLineReturn := models.KLineReturn{}

	mapParams := make(map[string]string)
	mapParams["symbol"] = strSymbol
	mapParams["period"] = strPeriod
	mapParams["size"] = strconv.Itoa(nSize)

	strRequestUrl := "/market/history/kline"
	strUrl := config.HuoBiConf.MarketUrl + strRequestUrl

	jsonKLineReturn, err := utils.HttpGetRequest(strUrl, mapParams)
	if err != "" {
		errJson := "{\"err\": \"" + err + "\"}"
		json.Unmarshal([]byte(errJson), &kLineReturn)
	} else {
		json.Unmarshal([]byte(jsonKLineReturn), &kLineReturn)
	}

	return kLineReturn
}

// Get aggregated quotes
// strSymbol: Transaction pair, btcusdt, bccbtc......
// return: TickReturn Object
func GetTicker(strSymbol string) models.TickerReturn {
	tickerReturn := models.TickerReturn{}

	mapParams := make(map[string]string)
	mapParams["symbol"] = strSymbol

	strRequestUrl := "/market/detail/merged"
	strUrl := config.HuoBiConf.MarketUrl + strRequestUrl

	jsonTickReturn, err := utils.HttpGetRequest(strUrl, mapParams)
	if err != "" {
		errJson := "{\"err\": \"" + err + "\"}"
		json.Unmarshal([]byte(errJson), &tickerReturn)
	} else {
		json.Unmarshal([]byte(jsonTickReturn), &tickerReturn)
	}

	return tickerReturn
}

// Get transaction depth information
// strSymbol: Transaction pair, btcusdt, bccbtc......
// strType: Depth type, step0、step1......stpe5 (Merge depth 0-5, 0 does not merge)
// return: HuobiDepthReturn Object
func HuobiDepth(strSymbol, strType string) models.HuobiDepthReturn {
	marketDepthReturn := models.HuobiDepthReturn{}

	mapParams := make(map[string]string)
	mapParams["symbol"] = strSymbol
	mapParams["type"] = strType

	strRequestUrl := "/market/depth"
	strUrl := config.HuoBiConf.MarketUrl + strRequestUrl

	jsonMarketDepthReturn, err := utils.HttpGetRequest(strUrl, mapParams)
	if err != "" {
		errJson := "{\"err\": \"" + err + "\"}"
		json.Unmarshal([]byte(errJson), &marketDepthReturn)
	} else {
		json.Unmarshal([]byte(jsonMarketDepthReturn), &marketDepthReturn)
	}

	return marketDepthReturn
}

// Get transaction details
// strSymbol: Transaction pair, btcusdt, bccbtc......
// return: TradeDetailReturn Object
func GetTradeDetail(strSymbol string) models.TradeDetailReturn {
	tradeDetailReturn := models.TradeDetailReturn{}

	mapParams := make(map[string]string)
	mapParams["symbol"] = strSymbol

	strRequestUrl := "/market/trade"
	strUrl := config.HuoBiConf.MarketUrl + strRequestUrl

	jsonTradeDetailReturn, err := utils.HttpGetRequest(strUrl, mapParams)
	if err != "" {
		errJson := "{\"err\": \"" + err + "\"}"
		json.Unmarshal([]byte(errJson), &tradeDetailReturn)
	} else {
		json.Unmarshal([]byte(jsonTradeDetailReturn), &tradeDetailReturn)
	}

	return tradeDetailReturn
}

// Get recent transaction history in bulk
// strSymbol: Transaction pair, btcusdt, bccbtc......
// nSize: Get the number of transaction records, range 1-2000
// return: TradeReturn Object
func HuobiTrade(strSymbol string, nSize int) models.TradeReturn {
	tradeReturn := models.TradeReturn{}

	mapParams := make(map[string]string)
	mapParams["symbol"] = strSymbol
	mapParams["size"] = strconv.Itoa(nSize)

	strRequestUrl := "/market/history/trade"
	strUrl := config.HuoBiConf.MarketUrl + strRequestUrl

	jsonTradeReturn, err := utils.HttpGetRequest(strUrl, mapParams)
	if err != "" {
		errJson := "{\"err\": \"" + err + "\"}"
		json.Unmarshal([]byte(errJson), &tradeReturn)
	} else {
		json.Unmarshal([]byte(jsonTradeReturn), &tradeReturn)
	}

	return tradeReturn
}

// Get Market Detail 24 hour volume data
// strSymbol: Transaction pair, btcusdt, bccbtc......
// return: MarketDetailReturn Object
func GetMarketDetail(strSymbol string) models.MarketDetailReturn {
	marketDetailReturn := models.MarketDetailReturn{}

	mapParams := make(map[string]string)
	mapParams["symbol"] = strSymbol

	strRequestUrl := "/market/detail"
	strUrl := config.HuoBiConf.MarketUrl + strRequestUrl

	jsonMarketDetailReturn, err := utils.HttpGetRequest(strUrl, mapParams)
	if err != "" {
		errJson := "{\"err\": \"" + err + "\"}"
		json.Unmarshal([]byte(errJson), &marketDetailReturn)
	} else {
		json.Unmarshal([]byte(jsonMarketDetailReturn), &marketDetailReturn)
	}

	return marketDetailReturn
}

//------------------------------------------------------------------------------------------
// Public API

// Query all transactions and precision supported by the system
// return: SymbolsReturn Object
func GetSymbols() models.SymbolsReturn {
	symbolsReturn := models.SymbolsReturn{}

	strRequestUrl := "/v1/common/symbols"
	strUrl := config.HuoBiConf.TradeUrl + strRequestUrl

	jsonSymbolsReturn, err := utils.HttpGetRequest(strUrl, nil)
	if err != "" {
		errJson := "{\"err\": \"" + err + "\"}"
		json.Unmarshal([]byte(errJson), &symbolsReturn)
	} else {
		json.Unmarshal([]byte(jsonSymbolsReturn), &symbolsReturn)
	}

	return symbolsReturn
}

// Query all currencies supported by the system
// return: CurrencysReturn Object
func GetCurrencys() models.CurrencysReturn {
	currencysReturn := models.CurrencysReturn{}

	strRequestUrl := "/v1/common/currencys"
	strUrl := config.HuoBiConf.TradeUrl + strRequestUrl

	jsonCurrencysReturn, err := utils.HttpGetRequest(strUrl, nil)
	if err != "" {
		errJson := "{\"err\": \"" + err + "\"}"
		json.Unmarshal([]byte(errJson), &currencysReturn)
	} else {
		json.Unmarshal([]byte(jsonCurrencysReturn), &currencysReturn)
	}

	return currencysReturn
}

// Query system current timestamp
// return: TimestampReturn Object
func GetTimestamp() models.TimestampReturn {
	timestampReturn := models.TimestampReturn{}

	strRequest := "/v1/common/timestamp"
	strUrl := config.HuoBiConf.TradeUrl + strRequest

	jsonTimestampReturn, err := utils.HttpGetRequest(strUrl, nil)
	if err != "" {
		errJson := "{\"err\": \"" + err + "\"}"
		json.Unmarshal([]byte(errJson), &timestampReturn)
	} else {
		json.Unmarshal([]byte(jsonTimestampReturn), &timestampReturn)
	}

	return timestampReturn
}

//------------------------------------------------------------------------------------------
// User Assets API

// Query all accounts of the current user, query according to the included private key
// return: AccountsReturn Object
func GetAccounts() models.AccountsReturn {
	accountsReturn := models.AccountsReturn{}

	strRequest := "/v1/account/accounts"
	jsonAccountsReturn, err := utils.ApiKeyGet(make(map[string]string), strRequest)
	if err != "" {
		errJson := "{\"err\": \"" + err + "\"}"
		json.Unmarshal([]byte(errJson), &accountsReturn)
	} else {
		json.Unmarshal([]byte(jsonAccountsReturn), &accountsReturn)
	}

	return accountsReturn
}

// Check account balance based on account ID
// nAccountID: Account ID, if you don't know, you can get it by GetAccounts(), you can only have spot account, C2C account, futures account.
// return: BalanceReturn Object
func GetAccountBalance(strAccountID string) models.BalanceReturn {
	balanceReturn := models.BalanceReturn{}

	strRequest := fmt.Sprintf("/v1/account/accounts/%s/balance", strAccountID)
	jsonBanlanceReturn, err := utils.ApiKeyGet(make(map[string]string), strRequest)
	if err != "" {
		errJson := "{\"err\": \"" + err + "\"}"
		json.Unmarshal([]byte(errJson), &balanceReturn)
	} else {
		json.Unmarshal([]byte(jsonBanlanceReturn), &balanceReturn)
	}

	return balanceReturn
}

//------------------------------------------------------------------------------------------
// Transaction API

// Order
// placeRequestParams: Order information
// return: PlaceReturn Object
func Place(placeRequestParams models.PlaceRequestParams) models.PlaceReturn {
	placeReturn := models.PlaceReturn{}

	mapParams := make(map[string]string)
	mapParams["account-id"] = placeRequestParams.AccountID
	mapParams["amount"] = placeRequestParams.Amount
	if 0 < len(placeRequestParams.Price) {
		mapParams["price"] = placeRequestParams.Price
	}
	if 0 < len(placeRequestParams.Source) {
		mapParams["source"] = placeRequestParams.Source
	}
	mapParams["symbol"] = placeRequestParams.Symbol
	mapParams["type"] = placeRequestParams.Type

	strRequest := "/v1/order/orders/place"
	jsonPlaceReturn := utils.ApiKeyPost(mapParams, strRequest)
	json.Unmarshal([]byte(jsonPlaceReturn), &placeReturn)

	return placeReturn
}

// Request to cancel an order request
// strOrderID: Order ID
// return: PlaceReturn Object
func SubmitCancel(strOrderID string) models.PlaceReturn {
	placeReturn := models.PlaceReturn{}

	strRequest := fmt.Sprintf("/v1/order/orders/%s/submitcancel", strOrderID)
	jsonPlaceReturn := utils.ApiKeyPost(make(map[string]string), strRequest)
	json.Unmarshal([]byte(jsonPlaceReturn), &placeReturn)

	return placeReturn
}
