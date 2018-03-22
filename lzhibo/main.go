package main

import (
	"time"
	"fmt"
	"net/http"
	"github.com/garyburd/redigo/redis"
	"strings"
	"./base"
)

const taskurl  = "http://127.0.0.1:5000/task"

// 按照一分钟 五分钟 一个小时 这三个时长

func main()  {
	//start()
	getData()
}

func newSafe() *base.SafeMap {
	sm := new(base.SafeMap)
	sm.Map = make(map[string]int64)
	return sm
}

func getData()  {

	count := newSafe()
	count.Map["timer"] = time.Now().Unix()+30 // 保存数据
	timerA := time.Now().Unix()+60 // 更新task花费的时间
	timerB := time.Now().Unix()+59 // 重新跑任务
	for {
		if time.Now().Unix() >= timerA{
			resp, _ := http.Get(taskurl) // 更新tasks
			fmt.Println(resp)
			timerA = time.Now().Unix()+60
		}
		if time.Now().Unix() >= timerB{
			longLink(count)
			timerB = time.Now().Unix()+59
		}
	}
}

func longLink(count *base.SafeMap)  {
	fmt.Println("start longlink")
	ch := make(chan string)
	rooms := getTask()
	for _, i := range rooms{
		go base.CountConnect(i, count)
	}
	for i:=1;i<len(rooms) ;i++  {
		<-ch
	}
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