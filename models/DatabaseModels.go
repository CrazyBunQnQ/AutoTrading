package models

import (
	"fmt"
)

func (account Account) String() string {
	return fmt.Sprintf("\nid:%s\nusdt:%.10f\nbtc:%.10f\neth:%.10f\nbnb:%.10f\neos:%.10f\nxrp:%.10f\n[createTime:%s updateTime:%s]\n",
		account.Id, account.Usdt, account.Btc, account.Eth, account.Bnb, account.Eos, account.Xrp, account.CreateTime, account.UpdateTime)
}

type Account struct {
	Id         string
	Usdt       float64
	Btc        float64
	Eth        float64
	Bnb        float64
	Eos        float64
	Xrp        float64
	CreateTime string //time.Time []uint8
	UpdateTime string //time.Time []uint8
}

func (strategyLowBuyHighSell StrategyLowBuyHighSell) String() string {
	return fmt.Sprintf("id:%\nsymbol:%s\nplatform:%sspend:%.10f\npositionAverage:%.10f\nlastSpend:%.10f\nprofitPoint:%.3f\ntargetSellPrice:%.10f\nbuyPoint:%.3f\ntargetBuyPrice:%.10f\nmonthAverage:%.10f[createTime:%s updateTime:%s]\n",
		strategyLowBuyHighSell.Id, strategyLowBuyHighSell.Symbol, strategyLowBuyHighSell.Platform, strategyLowBuyHighSell.Spend, strategyLowBuyHighSell.PositionAverage, strategyLowBuyHighSell.LastSpend,
		strategyLowBuyHighSell.TargetProfitPoint, strategyLowBuyHighSell.TargetSellPrice, strategyLowBuyHighSell.TargetBuyPoint, strategyLowBuyHighSell.TargetBuyPrice, strategyLowBuyHighSell.MonthAverage,
		strategyLowBuyHighSell.CreateTime, strategyLowBuyHighSell.UpdateTime)
}

type StrategyLowBuyHighSell struct {
	Id                int
	Symbol            string
	Platform          string
	Spend             float64
	PositionAverage   float64
	LastSpend         float64
	TargetProfitPoint float32
	TargetSellPrice   float64
	TargetBuyPoint    float32
	TargetBuyPrice    float64
	MonthAverage      float64
	CreateTime        string //time.Time []uint8
	UpdateTime        string //time.Time []uint8
}
