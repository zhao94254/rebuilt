package base


import (
	"github.com/gomodule/redigo/redis"
	"fmt"
	"strings"
)


func getTask() []string{
	c, err := redis.Dial("tcp", "127.0.0.1:6379")
	if err != nil {
		fmt.Println("Connect to redis error", err)
		return nil
	}
	u, err := redis.String(c.Do("GET", "douyu|task"))

	defer c.Close()
	return strings.Split(u, "|")
}