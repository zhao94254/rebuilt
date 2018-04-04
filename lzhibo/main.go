package main

import (
	"time"
	"fmt"
	"./base"
	"net/http"
)

func main() {
	getData()
}

func getData() {
	fmt.Println("start")
	count := base.NewSafe()
	count.Map["one|timer"] = time.Now().Unix() + base.OneTimer   // 保存数据
	count.Map["five|timer"] = time.Now().Unix() + base.FiveTimer // 保存数据
	count.Map["half|timer"] = time.Now().Unix() + base.HalfTimer //
	count.Map["restart"] = time.Now().Unix() + base.GetTaskTimer
	for {
		resp, _ := http.Get(base.Taskurl) // 更新tasks
		fmt.Println("Req new task", resp)
		longLink(count)
	}
}

func longLink(count *base.SafeMap) {
	fmt.Println("start all longlink")
	c := base.RedisClient()
	ch := make(chan string)
	rooms := base.GetTask()

	for _, i := range rooms {
		go base.CountConnect(i, count, c, ch)
	}
	for i := 1; i < len(rooms); i++ {
		fmt.Println(<-ch)
	}
	return
}
