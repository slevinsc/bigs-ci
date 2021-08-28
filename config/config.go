package config

import (
	"encoding/json"
	"os"
	"path/filepath"

	"bigs-ci/lib/logger"
)

var Config *Configure

type Configure struct {
	DockerAddr  string `json:"docker_addr"`
	GitlabAddr  string `json:"gitlab_addr"`
	MysqlConfig struct {
		Master struct {
			DSN     string `json:"dsn"`
			User    string `json:"user"`
			Pass    string `json:"pass"`
			Host    string `json:"host"`
			Port    int    `json:"port"`
			DBName  string `json:"dbname"`
			MaxOpen int    `json:"max_open"`
			MaxIdle int    `json:"max_idle"`
		} `json:"master"`
		Slave struct {
			DSN     string `json:"dsn"`
			User    string `json:"user"`
			Pass    string `json:"pass"`
			Host    string `json:"host"`
			Port    int    `json:"port"`
			DBName  string `json:"dbname"`
			MaxOpen int    `json:"max_open"`
			MaxIdle int    `json:"max_idle"`
		} `json:"slave"`
	} `json:"mysql_config"`
	RedisConf redisConf `json:"redis_conf"`
}

type redisConf struct {
		Addr     string `json:"addr" binding:"required"`
		Password string `json:"password" binding:"required"`
		DB       int    `json:"db" binding:"required"`
}

func Init() {
	loadConfigFromJson()

}

func loadConfigFromJson() {
	path, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	filePtr, err := os.Open(path + "/config/config.json")
	if err != nil {
		panic(err)
	}
	defer filePtr.Close()

	decoder := json.NewDecoder(filePtr)
	err = decoder.Decode(&Config)
	if err != nil {
		panic(err)
	}
	logger.DebugAsJson(Config)
}
