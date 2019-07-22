package models

import "time"

// OrderBook represents Bids and Asks.
type BianDepth struct {
	LastUpdateID int        `json:"lastUpdateId"`
	Bids         [][]string `json:"bids"` // 买盘, [price(成交价), amount(成交量)], 按price降序排列
	Asks         [][]string `json:"asks"` // 卖盘, [price(成交价), amount(成交量)], 按price升序排列
	Err          string     `json:"err"`
}

type BianTrade struct {
	ID           int    `json:"id"`
	Price        string `json:"price"`
	Quantity     string `json:"qty"`
	QuoteQty     string `json:"quoteQty"`
	Timestamp    int64  `json:"time"`
	IsBuyerMaker bool   `json:"isBuyerMaker"`
	IsBestMatch  bool   `json:"isBestMatch"`
	Err          string `json:"err"`
}

// AggTrade represents aggregated trade.
type BianAggTrade struct {
	ID           int    `json:"a"`
	Price        string `json:"p"`
	Quantity     string `json:"q"`
	FirstTradeID int    `json:"f"`
	LastTradeID  int    `json:"l"`
	Timestamp    int64  `json:"T"`
	IsBuyerMaker bool   `json:"m"`
	IsBestMatch  bool   `json:"M"`
	Err          string `json:"err"`
}

type BianAvgPrice struct {
	Mins  int    `json:"mins"`
	Price string `json:"price"`
	Err   string `json:"err"`
}

// Ticker24 represents data for 24hr ticker.
type BianTicker24 struct {
	Symbol             string `json:"symbol"`
	PriceChange        string `json:"priceChange"`
	PriceChangePercent string `json:"priceChangePercent"`
	WeightedAvgPrice   string `json:"weightedAvgPrice"`
	PrevClosePrice     string `json:"prevClosePrice"`
	LastPrice          string `json:"lastPrice"`
	LastQty            string `json:"lastQty"`
	BidPrice           string `json:"bidPrice"`
	BidQty             string `json:"bidQty"`
	AskPrice           string `json:"askPrice"`
	AskQty             string `json:"askQty"`
	OpenPrice          string `json:"openPrice"`
	HighPrice          string `json:"highPrice"`
	LowPrice           string `json:"lowPrice"`
	Volume             string `json:"volume"`
	OpenTime           int64  `json:"openTime"`
	CloseTime          int64  `json:"closeTime"`
	FirstID            int    `json:"firstId"`
	LastID             int    `json:"lastId"`
	Count              int    `json:"count"`
	Err                string `json:"err"`
}

// PriceTicker represents ticker data for price.
type BianLastPrice struct {
	Symbol string `json:"symbol"`
	Price  string `json:"price"`
	Err    string `json:"err"`
}

// BookTicker represents book ticker data.
type BianBestTicker struct {
	Symbol   string `json:"symbol"`
	BidPrice string `json:"bidPrice"`
	BidQty   string `json:"bidQty"`
	AskPrice string `json:"askPrice"`
	AskQty   string `json:"askQty"`
	Err      string `json:"err"`
}

// ProcessedOrder represents data from processed order.
type BianFastestOrder struct {
	Symbol        string `json:"symbol"`
	OrderID       int64  `json:"orderId"`
	ClientOrderID string `json:"clientOrderId"`
	TransactTime  int64  `json:"transactTime"`
	Err           string `json:"err"`
}

// ExecutedOrder represents data about executed order.
type BianOrderStatus struct {
	Symbol              string          `json:"symbol"`
	OrderID             int             `json:"orderId"`
	ClientOrderID       string          `json:"clientOrderId"`
	Price               string          `json:"price"`
	OrigQty             string          `json:"origQty"`
	ExecutedQty         string          `json:"executedQty"`
	CummulativeQuoteQty string          `json:"cummulativeQuoteQty"`
	Status              OrderStatus     `json:"status"`
	TimeInForce         BianTimeInForce `json:"timeInForce"`
	Type                BianOrderType   `json:"type"`
	Side                BianOrderSide   `json:"side"`
	StopPrice           string          `json:"stopPrice"`
	IcebergQty          string          `json:"icebergQty"`
	Time                int64           `json:"time"`
	UpdateTime          int64           `json:"updateTime"`
	IsWorking           bool            `json:"isWorking"`
	Err                 string          `json:"err"`
}

// Account represents user's account information.
type BianAccount struct {
	MakerCommision  int64  `json:"makerCommission"`
	TakerCommision  int64  `json:"takerCommission"`
	BuyerCommision  int64  `json:"buyerCommission"`
	SellerCommision int64  `json:"sellerCommission"`
	CanTrade        bool   `json:"canTrade"`
	CanWithdraw     bool   `json:"canWithdraw"`
	CanDeposit      bool   `json:"canDeposit"`
	UpdateTime      int64  `json:"updateTime"`
	AccountType     string `json:"accountType"`
	Balances        []BianBalance
	Err             string `json:"err"`
	Code            int    `json:"code"`
	Msg             string `json:"msg"`
}

// BianBalance groups balance-related information.
type BianBalance struct {
	Symbol string `json:"asset"`
	Free   string `json:"free"`
	Locked string `json:"locked"`
}

// **************** Enum *****************

// Binance BianInterval represents interval enum.
type BianInterval string

// OrderStatus represents order status enum.
type OrderStatus string

// BianOrderType represents order type enum.
type BianOrderType string

// OrderSide represents order side enum.
type BianOrderSide string

type BianOrderResponType string

// BianTimeInForce represents timeInForce enum.
type BianTimeInForce string

var (
	BianMinute         = BianInterval("1m")
	BianThreeMinutes   = BianInterval("3m")
	BianFiveMinutes    = BianInterval("5m")
	BianFifteenMinutes = BianInterval("15m")
	BianThirtyMinutes  = BianInterval("30m")
	BianHour           = BianInterval("1h")
	BianTwoHours       = BianInterval("2h")
	BianFourHours      = BianInterval("4h")
	BianSixHours       = BianInterval("6h")
	BianEightHours     = BianInterval("8h")
	BianTwelveHours    = BianInterval("12h")
	BianDay            = BianInterval("1d")
	BianThreeDays      = BianInterval("3d")
	BianWeek           = BianInterval("1w")
	BianMonth          = BianInterval("1M")

	StatusNew             = OrderStatus("NEW")
	StatusPartiallyFilled = OrderStatus("PARTIALLY_FILLED")
	StatusFilled          = OrderStatus("FILLED")
	StatusCancelled       = OrderStatus("CANCELED")
	StatusPendingCancel   = OrderStatus("PENDING_CANCEL")
	StatusRejected        = OrderStatus("REJECTED")
	StatusExpired         = OrderStatus("EXPIRED")

	TypeLimit       = BianOrderType("LIMIT")
	TypeMarket      = BianOrderType("MARKET")
	StopLoss        = BianOrderType("STOP_LOSS")
	StopLossLimit   = BianOrderType("STOP_LOSS_LIMIT")
	TakeProfit      = BianOrderType("TAKE_PROFIT")
	TakeProfitLimit = BianOrderType("TAKE_PROFIT_LIMIT")
	LimitMaker      = BianOrderType("LIMIT_MAKER")

	SideBuy  = BianOrderSide("BUY")
	SideSell = BianOrderSide("SELL")

	AckResponse    = BianOrderResponType("ACK")
	ResultResponse = BianOrderResponType("RESULT")
	FullResponse   = BianOrderResponType("FULL")

	GTC = BianTimeInForce("GTC")
	IOC = BianTimeInForce("IOC")
)

// **************************************

// OrderBook represents Bids and Asks.
type OrderBook struct {
	LastUpdateID int `json:"lastUpdateId"`
	Bids         []*Order
	Asks         []*Order
}

type DepthEvent struct {
	WSEvent
	UpdateID int
	OrderBook
}

// Order represents single order information.
type Order struct {
	Price    float64
	Quantity float64
}

// OrderBookRequest represents OrderBook request data.
type OrderBookRequest struct {
	Symbol string
	Limit  int
}

// AggTrade represents aggregated trade.
type AggTrade struct {
	ID             int
	Price          float64
	Quantity       float64
	FirstTradeID   int
	LastTradeID    int
	Timestamp      time.Time
	BuyerMaker     bool
	BestPriceMatch bool
}

type AggTradeEvent struct {
	WSEvent
	AggTrade
}

// AggTradesRequest represents AggTrades request data.
type AggTradesRequest struct {
	Symbol    string
	FromID    int64
	StartTime int64
	EndTime   int64
	Limit     int
}

// KlinesRequest represents Klines request data.
type KlinesRequest struct {
	Symbol    string
	Interval  BianInterval
	Limit     int
	StartTime int64
	EndTime   int64
}

// Kline represents single Kline information.
type Kline struct {
	OpenTime                 time.Time
	Open                     float64
	High                     float64
	Low                      float64
	Close                    float64
	Volume                   float64
	CloseTime                time.Time
	QuoteAssetVolume         float64
	NumberOfTrades           int
	TakerBuyBaseAssetVolume  float64
	TakerBuyQuoteAssetVolume float64
}

type KlineEvent struct {
	WSEvent
	Interval     BianInterval
	FirstTradeID int64
	LastTradeID  int64
	Final        bool
	Kline
}

// TickerRequest represents Ticker request data.
type TickerRequest struct {
	Symbol string
}

// Ticker24 represents data for 24hr ticker.
type Ticker24 struct {
	PriceChange        float64
	PriceChangePercent float64
	WeightedAvgPrice   float64
	PrevClosePrice     float64
	LastPrice          float64
	BidPrice           float64
	AskPrice           float64
	OpenPrice          float64
	HighPrice          float64
	LowPrice           float64
	Volume             float64
	OpenTime           time.Time
	CloseTime          time.Time
	FirstID            int
	LastID             int
	Count              int
}

// NewOrderRequest represents NewOrder request data.
type NewOrderRequest struct {
	Symbol           string
	Side             BianOrderSide
	Type             BianOrderType
	TimeInForce      BianTimeInForce
	Quantity         float64
	Price            float64
	NewClientOrderID string
	StopPrice        float64
	IcebergQty       float64
	Timestamp        time.Time
}

// QueryOrderRequest represents QueryOrder request data.
type QueryOrderRequest struct {
	Symbol            string
	OrderID           int64
	OrigClientOrderID string
	RecvWindow        time.Duration
	Timestamp         time.Time
}

// CancelOrderRequest represents CancelOrder request data.
type CancelOrderRequest struct {
	Symbol            string
	OrderID           int64
	OrigClientOrderID string
	NewClientOrderID  string
	RecvWindow        time.Duration
	Timestamp         time.Time
}

// CanceledOrder represents data about canceled order.
type CanceledOrder struct {
	Symbol            string
	OrigClientOrderID string
	OrderID           int64
	ClientOrderID     string
}

// OpenOrdersRequest represents OpenOrders request data.
type OpenOrdersRequest struct {
	Symbol     string
	RecvWindow time.Duration
	Timestamp  time.Time
}

// AllOrdersRequest represents AllOrders request data.
type AllOrdersRequest struct {
	Symbol     string
	OrderID    int64
	Limit      int
	RecvWindow time.Duration
	Timestamp  time.Time
}

// AccountRequest represents Account request data.
type AccountRequest struct {
	RecvWindow time.Duration
	Timestamp  time.Time
}

type AccountEvent struct {
	WSEvent
	BianAccount
}

// MyTradesRequest represents MyTrades request data.
type MyTradesRequest struct {
	Symbol     string
	Limit      int
	FromID     int64
	RecvWindow time.Duration
	Timestamp  time.Time
}

// Trade represents data about trade.
type Trade struct {
	ID              int64
	Price           float64
	Qty             float64
	Commission      float64
	CommissionAsset string
	Time            time.Time
	IsBuyer         bool
	IsMaker         bool
	IsBestMatch     bool
}

// WithdrawRequest represents Withdraw request data.
type WithdrawRequest struct {
	Asset      string
	Address    string
	Amount     float64
	Name       string
	RecvWindow time.Duration
	Timestamp  time.Time
}

// WithdrawResult represents Withdraw result.
type WithdrawResult struct {
	Success bool
	Msg     string
}

// HistoryRequest represents history-related calls request data.
type HistoryRequest struct {
	Asset      string
	Status     *int
	StartTime  time.Time
	EndTime    time.Time
	RecvWindow time.Duration
	Timestamp  time.Time
}

// Deposit represents Deposit data.
type Deposit struct {
	InsertTime time.Time
	Amount     float64
	Asset      string
	Status     int
}

// Withdrawal represents withdrawal data.
type Withdrawal struct {
	Amount    float64
	Address   string
	TxID      string
	Asset     string
	ApplyTime time.Time
	Status    int
}

// Stream represents stream information.
//
// Read web docs to get more information about using streams.
type Stream struct {
	ListenKey string
}

type WSEvent struct {
	Type   string
	Time   time.Time
	Symbol string
}

type DepthWebsocketRequest struct {
	Symbol string
}

type KlineWebsocketRequest struct {
	Symbol   string
	Interval BianInterval
}

type TradeWebsocketRequest struct {
	Symbol string
}

type UserDataWebsocketRequest struct {
	ListenKey string
}
