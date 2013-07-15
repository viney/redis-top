package main

import (
	"fmt"
	"github.com/hoisie/redis"
	"sync"
)

var (
	client *redis.Client
	mutex  sync.Mutex
)

func init() {
	mutex.Lock()
	defer mutex.Unlock()

	if client != nil {
		return
	}

	client = &redis.Client{
		Addr:        "127.0.0.1:6379",
		Db:          0, // default db is 0
		Password:    "admin",
		MaxPoolSize: 10000,
	}

	if err := client.Auth("admin"); err != nil {
		fmt.Println("Auth: ", err.Error())
		return
	}
}

// 添加数据到有序集合中
func Zadd(key string, value []byte, score float64) (bool, error) {
	return client.Zadd(key, value, score)
}

// 给有序集合指定成员，增加score
// 返回score结果
func Zincrby(key string, value []byte, score float64) (float64, error) {
	return client.Zincrby(key, value, score)
}

// 获取有序集合中所有成员和score
// 并获取成员对应的score值
func Zrevrange(key string, start int, end int) (map[string]float64, error) {
	bytes, err := client.Zrevrange(key, start, end)
	if err != nil {
		return nil, err
	}

	var ret = make(map[string]float64)
	if len(bytes) > 0 {
		for _, member := range bytes {
			f, err := client.Zscore(key, member)
			if err != nil {
				return nil, err
			}
			ret[string(member)] = f
		}
	}
	return ret, nil
}
