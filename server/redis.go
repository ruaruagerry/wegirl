package server

import (
	"fmt"
	"wegirl/servercfg"
	"time"

	"github.com/garyburd/redigo/redis"
	log "github.com/sirupsen/logrus"
)

const (
	redisconnectkey = "redisconnectkey"
)

var (
	pool *redis.Pool

	// LuaUNumberAdd script for add UNumber
	LuaUNumberAdd *redis.Script
)

func redisStartup() {
	if servercfg.RedisServer == "" {
		log.Panic("Must specify the RedisServer address in config json")
		return
	}

	pool = newPool(servercfg.RedisServer)

	createLuaScript()

	result := checkRedisKey(fmt.Sprintf("%d", servercfg.ServerID))
	if !result {
		log.Panic("check redis key failed")
	}
}

// newPool 新建redis连接池
func newPool(addr string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial:        func() (redis.Conn, error) { return redis.Dial("tcp", addr) },
	}
}

// createLuaScript lua脚本
func createLuaScript() {
	scriptV := `local value = redis.call('hget', KEYS[1], 'dfHMW')
		if type(value) ~= 'string' then
			value = 0
		else
			value = tonumber(value)
		end
		if value < tonumber(KEYS[2]) then
			redis.call('HSET', KEYS[1], 'dfHMW', KEYS[2])
		end
		value = redis.call('hget', KEYS[1], 'dfHML')
		if type(value) ~= 'string' then
			value = 0
		else
			value = tonumber(value)
		end
		if value > tonumber(KEYS[2]) then
			redis.call('HSET', KEYS[1], 'dfHML', KEYS[2])
		end
		return 0`

	LuaUNumberAdd = redis.NewScript(2, scriptV)
}

// checkRedisKey 测试redis key的正确性
func checkRedisKey(key string) bool {
	if key == "" {
		return false
	}

	con := pool.Get()
	exist, err := redis.Bool(con.Do("EXISTS", redisconnectkey))
	if err != nil {
		return false
	}

	if exist {
		// 没有uuid生成唯一连接标识
		rediskey, err := redis.String(con.Do("GET", redisconnectkey))
		if err != nil {
			return false
		}

		if key != rediskey {
			log.Printf("key:%s rediskey:%s", key, rediskey)
			return false
		}
	} else {
		_, err := con.Do("SET", redisconnectkey, key)
		if err != nil {
			return false
		}
	}

	return true
}
