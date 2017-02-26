package main

import (
	"database/sql"
	"fmt"
	"github.com/garyburd/redigo/redis"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

var DB *sql.DB
var redisConn redis.Conn
var redisPool *redis.Pool

func InitDB() error {
	var err error

	dsn := GetMysqlDSN()

	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("Open database error: ", err.Error())
	}

	DB.SetMaxOpenConns(100)
	DB.SetMaxIdleConns(100)
	err = DB.Ping()
	if err != nil {
		return fmt.Errorf("Ping database error: ", err.Error())
	}

	return nil
}

func InitRedis() error {

	var err error

	addr := GetRedisAddr()
	redisConn, err = redis.Dial("tcp", addr)
	if err != nil {
		return fmt.Errorf("dail redis ", err.Error())
	}

	pool := 16
	redisPool = &redis.Pool{
		MaxIdle:     pool,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", addr)
			if err != nil {
				return nil, err
			}

			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}

	return nil
}

func Quit() {

	if DB != nil {
		DB.Close()
	}

	if redisConn != nil {
		redisConn.Close()
	}

	if redisPool != nil {
		redisPool.Close()
	}
}
