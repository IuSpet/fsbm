package conf

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"
)

const SuffixPath = "/conf/deploy.json"

type Env_ int

const (
	TEST Env_ = iota
	PRODUCT
)

var env Env_

type AllConfig struct {
	Product EnvConfig `json:"product"`
	Test    EnvConfig `json:"test"`
}

type EnvConfig struct {
	Mysql mysqlConfig
	Redis redisConfig
}

type mysqlConfig struct {
	Ip       string `json:"ip"`
	Port     string `json:"port"`
	DbName   string `json:"db_name"`
	User     string `json:"user"`
	Password string `json:"password"`
}

type redisConfig struct {
	Addr     string `json:"addr"`
	Password string `json:"password"`
	DB       int    `json:"db"`
}

var allCfg AllConfig
var GlobalConfig EnvConfig

func init() {
	product := os.Getenv("FSBM_PRODUCT")
	if product != "" {
		env = PRODUCT
	} else {
		env = TEST
	}
}

func Init() {
	currentPath, _ := os.Getwd()
	index := strings.Index(currentPath, "fsbm") + 4
	path := currentPath[:index] + SuffixPath
	data, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(data, &allCfg)
	if err != nil {
		panic(err)
	}
	if env == PRODUCT {
		GlobalConfig = allCfg.Product
	} else if env == TEST {
		GlobalConfig = allCfg.Test
	} else {
		panic("invalid environment")
	}
}

func GetEnv() Env_ {
	return env
}
