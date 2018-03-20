package main

import (
	"time"
	"fmt"
	"net/http"
	"github.com/garyburd/redigo/redis"
	"strings"
)

const taskurl  = "http://127.0.0.1:5000/task"

func main()  {
	start()
}

func start()  {
	timerA := time.Now().Unix()+10
	timerB := time.Now().Unix()+5
	fmt.Println(timerA, timerA)
	for {
		if time.Now().Unix() >= timerA{
			resp, _ := http.Get(taskurl) // 更新一次
			fmt.Println(resp)
			timerA = time.Now().Unix()+10
		}
		if time.Now().Unix() >= timerB{
			fmt.Println(getTask())
			timerB = time.Now().Unix()+5
		}


	}
}

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