package config

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"os"
)

type conf struct {
	Database database `yaml:"database"`
	Binance  binance  `yaml:"binance"`
	Huobi    huobi    `yaml:"huobi"`
	Okex     okex     `yaml:"okex"`
	Otcbtc   otcbtc   `yaml:"otcbtc"`
	Version  string   `yaml:"version"`
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

const configFile = "../config/config_private.yaml"

//const configFile = "config/config_private.yaml"

var DBConf database
var BianConf binance
var HuoBiConf huobi
var OkexConf okex
var OtcbtcConf otcbtc

func init() {
	yamlFile, _ := ioutil.ReadFile(configFile)
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
}
