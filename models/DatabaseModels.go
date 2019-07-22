package models

import (
	"fmt"
	"time"
)

type Platform struct {
	Id         int       `orm:"pk;auto"`
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
	Id         int       `orm:"pk;auto"`
	Platform   string    `orm:size(15)`
	Usdt       float64   `orm:"digits(18);decimals(10)"`
	Btc        float64   `orm:"digits(18);decimals(10)"`
	Eth        float64   `orm:"digits(18);decimals(10)"`
	Bnb        float64   `orm:"digits(18);decimals(10)"`
	Eos        float64   `orm:"digits(18);decimals(10)"`
	Xrp        float64   `orm:"digits(18);decimals(10)"`
	CreateTime time.Time `orm:"auto_now_add;type(datetime)"`
	UpdateTime time.Time `orm:"auto_now;type(datetime)"`
}

func (strategyLowBuyHighSell StrategyLowBuyHighSell) String() string {
	return fmt.Sprintf("id:%d\nsymbol:%s\nplatform:%sspend:%.10f\npositionAverage:%.10f\nlastSpend:%.10f\nprofitPoint:%.3f\ntargetSellPrice:%.10f\nbuyPoint:%.3f\ntargetBuyPrice:%.10f\nmonthAverage:%.10f[createTime:%s updateTime:%s]\n",
		strategyLowBuyHighSell.Id, strategyLowBuyHighSell.Symbol, strategyLowBuyHighSell.Platform, strategyLowBuyHighSell.Spend, strategyLowBuyHighSell.PositionAverage, strategyLowBuyHighSell.LastSpend,
		strategyLowBuyHighSell.TargetProfitPoint, strategyLowBuyHighSell.TargetSellPrice, strategyLowBuyHighSell.TargetBuyPoint, strategyLowBuyHighSell.TargetBuyPrice, strategyLowBuyHighSell.MonthAverage,
		strategyLowBuyHighSell.CreateTime, strategyLowBuyHighSell.UpdateTime)
}

type StrategyLowBuyHighSell struct {
	Id                int       `orm:"pk;auto"`
	Symbol            string    `orm:"size(15)"`
	Platform          string    `orm:"size(15)"`
	Spend             float64   `orm:"digits(18);decimals(10)"`
	PositionAverage   float64   `orm:"digits(18);decimals(10)"`
	LastSpend         float64   `orm:"digits(18);decimals(10)"`
	TargetProfitPoint float32   `orm:"digits(4);decimals(3);default(1.025)"`
	TargetSellPrice   float64   `orm:"digits(18);decimals(10)"`
	TargetBuyPoint    float32   `orm:"digits(4);decimals(3);default(0.95)"`
	TargetBuyPrice    float64   `orm:"digits(18);decimals(10)"`
	MonthAverage      float64   `orm:"digits(18);decimals(10)"`
	CreateTime        time.Time `orm:"auto_now_add;type(datetime)"`
	UpdateTime        time.Time `orm:"auto_now;type(datetime)"`
}
