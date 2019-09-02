# AutoTrading

AutoTrading is a process that automates trading.

It obtains the current price of mainstream currencies from Binance, Otcbtc, Okex, huobi, and other platforms. Calculate the price difference of each coin on different platforms. Determine if it's worth trading at this point based on a particular algorithm.

If you decide to trade, trade on both platforms to gain revenue.

```bash
# 打包资源
go-bindata -o=./config/asset.go -pkg=config resources/...
```

```bash
GOROOT=/usr/local/Cellar/go/1.12.7/libexec #gosetup
GOPATH=/Users/baojunjie/go #gosetup
# 生成 linux 平台二进制文件
CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' -o /Users/baojunjie/go/src/AutoTrading/out/main /Users/baojunjie/go/src/AutoTrading/main/main.go
```

```bash
# 打包资源，生成 linux 平台二进制可执行文件，并上传到服务器 一气呵成
go-bindata -o=./config/asset.go -pkg=config resources/... && CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' -o /Users/baojunjie/go/src/AutoTrading/out/main /Users/baojunjie/go/src/AutoTrading/main/main.go && scp out/main myserver:~/trade
```

```bash
~/autotrading/main
```

```bash
firewall-cmd --zone=public --add-port=80/tcp --permanent
firewall-cmd --reload && firewall-cmd --zone=public --list-ports
```
