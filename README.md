# AutoTrading

AutoTrading is a process that automates trading.

It obtains the current price of mainstream currencies from Binance, Otcbtc, Okex, huobi, and other platforms. Calculate the price difference of each coin on different platforms. Determine if it's worth trading at this point based on a particular algorithm.

If you decide to trade, trade on both platforms to gain revenue.

```bash
GOROOT=/usr/local/Cellar/go/1.12.7/libexec #gosetup
GOPATH=/Users/baojunjie/go #gosetup
CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' -o /Users/baojunjie/go/src/AutoTrading/out/main /Users/baojunjie/go/src/AutoTrading/main/main.go
cd out
docker build -t crazybun/autotrading:0.1 .
docker run -it -p 8000:8000 crazybun/autotrading:0.1
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

```bash
# 打包资源，生成 linux 平台二进制可执行文件，并上传到服务器
go-bindata -o=./config/asset.go -pkg=config resources/... && CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' -o /Users/baojunjie/go/src/AutoTrading/out/main /Users/baojunjie/go/src/AutoTrading/main/main.go && scp out/main myserver:~/trade
```