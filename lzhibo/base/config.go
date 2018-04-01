package base

import (
	"github.com/gomodule/redigo/redis"
	"strings"
)

const (
	BufferSize   = 1024
	ServerAddr   = "openbarrage.douyutv.com:8601"
	PostCode     = 689
	PullCode     = 690
	OneTimer     = 60
	FiveTimer    = 60 * 5
	GetTaskTimer = 60 * 30
	HalfTimer    = 60 * 30
	Taskurl      = "http://127.0.0.1:5000/task"
)

func RedisClient() redis.Conn {
	c, _ := redis.Dial("tcp", "127.0.0.1:6379")
	return c
}

func GetTask() []string {
	c := RedisClient()
	u, _ := redis.String(c.Do("GET", "douyu|task"))

	defer c.Close()
	return strings.Split(u, "|")
}