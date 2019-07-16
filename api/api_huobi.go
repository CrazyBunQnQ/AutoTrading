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
func GetKLine(strSymbol, strPeriod string, nSize int) models.KLineReturn {
	kLineReturn := models.KLineReturn{}

	mapParams := make(map[string]string)
	mapParams["symbol"] = strSymbol
	mapParams["period"] = strPeriod
	mapParams["size"] = strconv.Itoa(nSize)

	strRequestUrl := "/market/history/kline"
	strUrl := config.HuoBiConf.MarketUrl + strRequestUrl

	jsonKLineReturn := utils.HttpGetRequest(strUrl, mapParams)
	json.Unmarshal([]byte(jsonKLineReturn), &kLineReturn)

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

	jsonTickReturn := utils.HttpGetRequest(strUrl, mapParams)
	json.Unmarshal([]byte(jsonTickReturn), &tickerReturn)

	return tickerReturn
}

// Get transaction depth information
// strSymbol: Transaction pair, btcusdt, bccbtc......
// strType: Depth type, step0、step1......stpe5 (Merge depth 0-5, 0 does not merge)
// return: MarketDepthReturn Object
func GetMarketDepth(strSymbol, strType string) models.MarketDepthReturn {
	marketDepthReturn := models.MarketDepthReturn{}

	mapParams := make(map[string]string)
	mapParams["symbol"] = strSymbol
	mapParams["type"] = strType

	strRequestUrl := "/market/depth"
	strUrl := config.HuoBiConf.MarketUrl + strRequestUrl

	jsonMarketDepthReturn := utils.HttpGetRequest(strUrl, mapParams)
	json.Unmarshal([]byte(jsonMarketDepthReturn), &marketDepthReturn)

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

	jsonTradeDetailReturn := utils.HttpGetRequest(strUrl, mapParams)
	json.Unmarshal([]byte(jsonTradeDetailReturn), &tradeDetailReturn)

	return tradeDetailReturn
}

// Get recent transaction history in bulk
// strSymbol: Transaction pair, btcusdt, bccbtc......
// nSize: Get the number of transaction records, range 1-2000
// return: TradeReturn Object
func GetTrade(strSymbol string, nSize int) models.TradeReturn {
	tradeReturn := models.TradeReturn{}

	mapParams := make(map[string]string)
	mapParams["symbol"] = strSymbol
	mapParams["size"] = strconv.Itoa(nSize)

	strRequestUrl := "/market/history/trade"
	strUrl := config.HuoBiConf.MarketUrl + strRequestUrl

	jsonTradeReturn := utils.HttpGetRequest(strUrl, mapParams)
	json.Unmarshal([]byte(jsonTradeReturn), &tradeReturn)

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

	jsonMarketDetailReturn := utils.HttpGetRequest(strUrl, mapParams)
	json.Unmarshal([]byte(jsonMarketDetailReturn), &marketDetailReturn)

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

	jsonSymbolsReturn := utils.HttpGetRequest(strUrl, nil)
	json.Unmarshal([]byte(jsonSymbolsReturn), &symbolsReturn)

	return symbolsReturn
}

// Query all currencies supported by the system
// return: CurrencysReturn Object
func GetCurrencys() models.CurrencysReturn {
	currencysReturn := models.CurrencysReturn{}

	strRequestUrl := "/v1/common/currencys"
	strUrl := config.HuoBiConf.TradeUrl + strRequestUrl

	jsonCurrencysReturn := utils.HttpGetRequest(strUrl, nil)
	json.Unmarshal([]byte(jsonCurrencysReturn), &currencysReturn)

	return currencysReturn
}

// Query system current timestamp
// return: TimestampReturn Object
func GetTimestamp() models.TimestampReturn {
	timestampReturn := models.TimestampReturn{}

	strRequest := "/v1/common/timestamp"
	strUrl := config.HuoBiConf.TradeUrl + strRequest

	jsonTimestampReturn := utils.HttpGetRequest(strUrl, nil)
	json.Unmarshal([]byte(jsonTimestampReturn), &timestampReturn)

	return timestampReturn
}

//------------------------------------------------------------------------------------------
// User Assets API

// Query all accounts of the current user, query according to the included private key
// return: AccountsReturn Object
func GetAccounts() models.AccountsReturn {
	accountsReturn := models.AccountsReturn{}

	strRequest := "/v1/account/accounts"
	jsonAccountsReturn := utils.ApiKeyGet(make(map[string]string), strRequest)
	json.Unmarshal([]byte(jsonAccountsReturn), &accountsReturn)

	return accountsReturn
}

// Check account balance based on account ID
// nAccountID: Account ID, if you don't know, you can get it by GetAccounts(), you can only have spot account, C2C account, futures account.
// return: BalanceReturn Object
func GetAccountBalance(strAccountID string) models.BalanceReturn {
	balanceReturn := models.BalanceReturn{}

	strRequest := fmt.Sprintf("/v1/account/accounts/%s/balance", strAccountID)
	jsonBanlanceReturn := utils.ApiKeyGet(make(map[string]string), strRequest)
	json.Unmarshal([]byte(jsonBanlanceReturn), &balanceReturn)

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
