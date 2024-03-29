package models

import (
	"fmt"
	"time"
)

type Platform struct {
	Id         int64     `orm:"pk;auto"`
	Name       string    `orm:size(15);unique`
	NameCn     string    `orm:size(10);null`
	CreateTime time.Time `orm:"auto_now_add;type(datetime)"`
	UpdateTime time.Time `orm:"auto_now;type(datetime)"`
}

func (account Account) String() string {
	return fmt.Sprintf("\nid:%d\nplatform:%s\nusdt:%.10f\nbtc:%.10f\neth:%.10f\nbnb:%.10f\neos:%.10f\nxrp:%.10f\n[createTime:%s updateTime:%s]\n",
		account.Id, account.Platform, account.Usdt, account.Btc, account.Eth, account.Bnb, account.Eos, account.Xrp, account.CreateTime, account.UpdateTime)
}

type Account struct {
	Id         int64     `orm:"pk;auto"`
	Platform   string    `orm:size(15)`
	Usdt       float64   `orm:"digits(18);decimals(10)"`
	UsdtLocked float64   `orm:"digits(18);decimals(10)"`
	Btc        float64   `orm:"digits(18);decimals(10)"`
	BtcLocked  float64   `orm:"digits(18);decimals(10)"`
	Eth        float64   `orm:"digits(18);decimals(10)"`
	EthLocked  float64   `orm:"digits(18);decimals(10)"`
	Bnb        float64   `orm:"digits(18);decimals(10)"`
	BnbLocked  float64   `orm:"digits(18);decimals(10)"`
	Eos        float64   `orm:"digits(18);decimals(10)"`
	EosLocked  float64   `orm:"digits(18);decimals(10)"`
	Xrp        float64   `orm:"digits(18);decimals(10)"`
	XrpLocked  float64   `orm:"digits(18);decimals(10)"`
	CreateTime time.Time `orm:"auto_now_add;type(datetime)"`
	UpdateTime time.Time `orm:"auto_now;type(datetime)"`
}

func (quantity Quantity) String() string {
	return fmt.Sprintf("\nid:%d\nname:%s\nplatform:%s\nfree:%.10f\nlocked:%.10f\n[createTime:%s updateTime:%s]\n",
		quantity.Id, quantity.Name, quantity.Platform, quantity.Free, quantity.Locked, quantity.CreateTime, quantity.UpdateTime)
}

type Quantity struct {
	Id         int64     `orm:"pk;auto"`
	Name       string    `orm:"size(15)"`
	Platform   string    `orm:"size(15)"`
	Free       float64   `orm:"digits(18);decimals(10)"`
	Locked     float64   `orm:"digits(18);decimals(10)"`
	CreateTime time.Time `orm:"auto_now_add;type(datetime)"`
	UpdateTime time.Time `orm:"auto_now;type(datetime)"`
}

func (strategyLowBuyHighSell StrategyLowBuyHighSell) String() string {
	return fmt.Sprintf("id:%d\nsymbol:%s\ncoinName:%s\nplatform:%s\nquantity:%.10f\nspend:%.10f\npositionAverage:%.10f\nlastSpend:%.10f\nspendCoefficient:%.3f\nprofitPoint:%.3f\ntargetSellPrice:%.10f\nbuyPoint:%.3f\ntargetBuyPrice:%.10f\nmonthAverage:%.10f\nenabled:%t\nacturlCost:%.10f\n[createTime:%s updateTime:%s]\n",
		strategyLowBuyHighSell.Id, strategyLowBuyHighSell.Symbol, strategyLowBuyHighSell.CoinName, strategyLowBuyHighSell.Platform, strategyLowBuyHighSell.Quantity, strategyLowBuyHighSell.Spend, strategyLowBuyHighSell.SpendCoefficient, strategyLowBuyHighSell.PositionAverage, strategyLowBuyHighSell.LastSpend,
		strategyLowBuyHighSell.TargetProfitPoint, strategyLowBuyHighSell.TargetSellPrice, strategyLowBuyHighSell.TargetBuyPoint, strategyLowBuyHighSell.TargetBuyPrice, strategyLowBuyHighSell.MonthAverage, strategyLowBuyHighSell.Enabled, strategyLowBuyHighSell.ActualCost,
		strategyLowBuyHighSell.CreateTime, strategyLowBuyHighSell.UpdateTime)
}

type StrategyLowBuyHighSell struct {
	Id                int64     `orm:"pk;auto"`
	Symbol            string    `orm:"size(15)"`
	CoinName          string    `orm:"size(11)"`
	Platform          string    `orm:"size(15)"`
	Quantity          float64   `orm:"digits(18);decimals(10)"`
	Spend             float64   `orm:"digits(18);decimals(10)"`
	PositionAverage   float64   `orm:"digits(18);decimals(10)"`
	LastSpend         float64   `orm:"digits(18);decimals(10)"`
	SpendCoefficient  float64   `orm:"digits(4);decimals(3)"`
	TargetProfitPoint float64   `orm:"digits(4);decimals(3);default(1.025)"`
	TargetSellPrice   float64   `orm:"digits(18);decimals(10)"`
	TargetBuyPoint    float64   `orm:"digits(4);decimals(3);default(0.95)"`
	TargetBuyPrice    float64   `orm:"digits(18);decimals(10)"`
	MonthAverage      float64   `orm:"digits(18);decimals(10)"`
	LastOrderId       int64     `orm:default(0)`
	Enabled           bool      `orm:default(false)`
	ActualCost        float64   `orm:"digits(18);decimals(10)"`
	CreateTime        time.Time `orm:"auto_now_add;type(datetime)"`
	UpdateTime        time.Time `orm:"auto_now;type(datetime)"`
}

type OrderHistory struct {
	OrderId       int64     `orm:"pk"`
	Symbol        string    `orm:"size(15)"`
	ClientOrderId string    `orm:"size(50)"`
	Price         float64   `orm:"digits(18);decimals(10)"`
	OrigQty       float64   `orm:"digits(18);decimals(10)"`
	Status        string    `orm:"size(20)"`
	TimeInForce   string    `orm:"size(10)"`
	Type          string    `orm:"size(10)"`
	Side          string    `orm:"size(10)"`
	StopPrice     float64   `orm:"digits(18);decimals(10)"`
	IcebergQty    float64   `orm:"digits(18);decimals(10)"`
	Time          time.Time `orm:"auto_now_add;type(datetime)"`
	UpdateTime    time.Time `orm:"auto_now;type(datetime)"`
	IsWorking     bool      `orm:default(true)`
	Platform      string    `orm:"size(15)"`
}
