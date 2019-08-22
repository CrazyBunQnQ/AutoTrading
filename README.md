# AutoTrading

AutoTrading is a process that automates trading.

It obtains the current price of mainstream currencies from Binance, Otcbtc, Okex, huobi, and other platforms. Calculate the price difference of each coin on different platforms. Determine if it's worth trading at this point based on a particular algorithm.

If you decide to trade, trade on both platforms to gain revenue.

```bash
GOROOT=/usr/local/Cellar/go/1.12.7/libexec #gosetup
GOPATH=/Users/baojunjie/go #gosetup
CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' -o /Users/baojunjie/go/src/AutoTrading/out/main /Users/baojunjie/go/src/AutoTrading/main/main.go
cd out
```

```bash
~/autotrading/main
```

```bash
firewall-cmd --zone=public --add-port=80/tcp --permanent
firewall-cmd --reload && firewall-cmd --zone=public --list-ports
```

```bash
go-bindata -o=./config/asset.go -pkg=config resources/...
```