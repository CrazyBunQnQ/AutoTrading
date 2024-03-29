package config

import (
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

type conf struct {
	Database              database `yaml:"database"`
	Binance               binance  `yaml:"binance"`
	Huobi                 huobi    `yaml:"huobi"`
	Okex                  okex     `yaml:"okex"`
	Otcbtc                otcbtc   `yaml:"otcbtc"`
	PlatformDiffPoint     float64  `yaml:"platform_diff_point"`
	PlatformBalancedPoint float64  `yaml:"platform_balanced_point"`
	PlatformOffset        float64  `yaml:"platform_offset"`
	AmountPerTrade        float64  `yaml:"amount_per_trade"`
	Ifttt                 ifttt    `yaml:"ifttt"`
	Version               string   `yaml:"version"`
}

type database struct {
	Addr     string `yaml:"addr"`
	Port     string `yaml:"port"`
	Schema   string `yaml:"schema"`
	UserName string `yaml:"uname"`
	Password string `yaml:"pwd"`
}

type binance struct {
	BaseUrl          string `yaml:"baseurl"`
	ApiKeyPrivate    string `yaml:"api_key_private"`
	SecretKeyPrivate string `yaml:"secret_key_private"`
	ApiKeyPublic     string `yaml:"api_key_public"`
	SecretKeyPublic  string `yaml:"secret_key_public"`
}

type huobi struct {
	BaseUrl                string `yaml:"baseurl"`
	MarketUrl              string `yaml:"market_url"`
	TradeUrl               string `yaml:"trade_url"`
	HostName               string `yaml:"host_name"`
	AccessKeyPrivate       string `yaml:"access_key_private"`
	SecretKeyPrivate       string `yaml:"secret_key_private"`
	AccessKeyPublic        string `yaml:"access_key_public"`
	SecretKeyPublic        string `yaml:"secret_key_public"`
	EnablePrivateSignature bool   `yaml:"enable_private_signature"`
	PrivateKeyPrime256     string `yaml:"private_key_prime_256"`
}

type okex struct {
}

type otcbtc struct {
	BaseUrl string `yaml:"baseurl"`
}

type ifttt struct {
	Enabled     bool   `yaml:"enabled"`
	Key         string `yaml:"key"`
	WebhooksUrl string `yaml:"webhooks_url"`
	EventName   string `yaml:"event_name"`
}

const configFile = "resources/config_private.yaml"

var DBConf database
var BianConf binance
var HuoBiConf huobi
var OkexConf okex
var OtcbtcConf otcbtc
var PlatformDiffPoint float64
var PlatformBalancedPoint float64
var PlatformOffset float64
var AmountPerTrade float64
var Ifttt ifttt
var Version string

func init() {
	yamlFile, _ := Asset(configFile)

	if yamlFile == nil {
		log.Println("未找到配置文件 " + configFile)
		os.Exit(1)
	}
	var c conf
	err := yaml.Unmarshal(yamlFile, &c)
	if err != nil {
		log.Fatalf("cannot unmarshal data: %v", err)
	}
	DBConf = c.Database
	BianConf = c.Binance
	HuoBiConf = c.Huobi
	OkexConf = c.Okex
	OtcbtcConf = c.Otcbtc
	PlatformDiffPoint = c.PlatformDiffPoint
	PlatformBalancedPoint = c.PlatformBalancedPoint
	AmountPerTrade = c.AmountPerTrade
	PlatformOffset = c.PlatformOffset
	Ifttt = c.Ifttt
	Version = c.Version
}
