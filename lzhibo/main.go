package main

import (
	"time"
	"fmt"
	"strings"
	"./base"
	"github.com/gomodule/redigo/redis"
)

const (
	oneTimer = 60
	fiveTimer = 60*5
	halfTimer = 60*30
	getTaskTimer = 60*31

)

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

func redisClient() redis.Conn {
	c, _ := redis.Dial("tcp", "127.0.0.1:6379")
	return c
}


func getData()  {
	fmt.Println("start")
	count := newSafe()
	count.Map["one|timer"] = time.Now().Unix()+oneTimer // 保存数据
	count.Map["five|timer"] = time.Now().Unix()+fiveTimer // 保存数据
	count.Map["half|timer"] = time.Now().Unix()+halfTimer
	count.Map["restart"] = time.Now().Unix()+getTaskTimer
	reRunTask := time.Now().Unix()+getTaskTimer // 重新跑任务
	first := true
	for {
		if time.Now().Unix() >= reRunTask || first{
			go longLink(count)
			reRunTask = time.Now().Unix()+getTaskTimer
		}
			first = false
	}
}


func longLink(count *base.SafeMap)  {
	fmt.Println("start all longlink")
	c  := redisClient()
	ch := make(chan string)
	rooms := getTask()

	for _, i := range rooms{
		go base.CountConnect(i, count, c)
	}
	for i:=1;i<len(rooms) ;i++  {
		<-ch
	}
}

func getTask() []string{
	c := redisClient()
	u, _ := redis.String(c.Do("GET", "douyu|task"))

	defer c.Close()
	return strings.Split(u, "|")
}
