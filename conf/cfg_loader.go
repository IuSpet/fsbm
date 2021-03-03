package conf

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"
)

const SuffixPath = "/conf/deploy.json"

type config struct {
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

var GlobalConfig config

func Init() {
	currentPath, _ := os.Getwd()
	index := strings.Index(currentPath,"fsbm") + 4
	path := currentPath[:index] + SuffixPath
	data, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(data, &GlobalConfig)
	if err != nil {
		panic(err)
	}
}
