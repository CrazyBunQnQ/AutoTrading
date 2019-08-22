# AutoTrading

AutoTrading is a process that automates trading.

It obtains the current price of mainstream currencies from Binance, Otcbtc, Okex, huobi, and other platforms. Calculate the price difference of each coin on different platforms. Determine if it's worth trading at this point based on a particular algorithm.

If you decide to trade, trade on both platforms to gain revenue.

```bash
go-bindata -o=./config/asset.go -pkg=config resources/...
```