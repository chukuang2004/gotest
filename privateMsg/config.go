package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type MysqlServer struct {
	Host     string `json:"host"`
	Database string `json:"database"`
	User     string `json:"user"`
	Password string `json:"password"`
}

type RedisServer struct {
	Host     string `json:"host"`
	Password string `json:"password"`
}
type Config struct {
	ListenAddr     string      `json:"listen_addr"`
	Mysql          MysqlServer `json:"mysql"`
	Redis          RedisServer `json:"redis"`
	RedisIdle      int         `json:"redis_idle"`
	RedisMaxActive int         `json:"redis_max_active"`
	LogLevel       string      `json:"loglevel"`
}

var config *Config = nil

func InitConfiguration(physicalPath string) (*Config, error) {
	file, err := os.Open(physicalPath)
	if err != nil {
		return nil, err
	}
	decoder := json.NewDecoder(file)

	var c Config
	err = decoder.Decode(&c)
	if err != nil {
		return nil, err
	}

	config = &c
	return config, nil
}

func GetListenAddr() string {

	if config == nil {
		return ""
	}

	return config.ListenAddr
}

func GetLogLevel() string {

	if config == nil {
		return ""
	}

	return config.LogLevel
}

// root:123456@tcp(10.1.9.102:3306)/apprank?timeout=5s
func GetMysqlDSN() string {

	if config == nil {
		return ""
	}

	return fmt.Sprintf("%s:%s@tcp(%s)/%s?timeout=5s", config.Mysql.User, config.Mysql.Password, config.Mysql.Host, config.Mysql.Database)
}

func GetRedisAddr() string {

	if config == nil {
		return ""
	}

	return config.Redis.Host
}
