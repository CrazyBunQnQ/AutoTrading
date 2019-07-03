package api

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
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
	Port     int32  `yaml:"port"`
	UserName string `yaml:"uname"`
	Password string `yaml:"pwd"`
}

type binance struct {
	BaseUrl          string `yaml:"baseurl"`
	UserName         string `yaml:"uname"`
	Password         string `yaml:"pwd"`
	ApiKeyPrivate    string `yaml:"api_key_private"`
	SecretKeyPrivate string `yaml:"secret_key_private"`
	ApiKeyPublic     string `yaml:"api_key_public"`
	SecretKeyPublic  string `yaml:"secret_key_public"`
}

type huobi struct {
}

type okex struct {
}

type otcbtc struct {
}

const configFile = "config_private.yaml"

func Conf() conf {
	yamlFile, _ := ioutil.ReadFile(configFile)
	var c conf
	err := yaml.Unmarshal(yamlFile, &c)
	if err != nil {
		log.Fatalf("cannot unmarshal data: %v", err)
	}
	return c
}