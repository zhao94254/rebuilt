package base

import (
	"fmt"
	"net"
	"encoding/binary"
	"bytes"
	"errors"
	"strings"
	"time"
	"github.com/gomodule/redigo/redis"
)

// 主要的处理逻辑



func PostData(msg string) []byte {
	// 构造需要发送的二进制数据
	length := 9 + len(msg) // 长度4字节 + 类型2字节 + 加密字段1字节 + 保留字段1字节 + 结尾字段1字节
	buffer := bytes.NewBuffer([]byte{})
	binary.Write(buffer, binary.LittleEndian, int32(length))
	binary.Write(buffer, binary.LittleEndian, int32(length))
	binary.Write(buffer, binary.LittleEndian, int16(PostCode))
	binary.Write(buffer, binary.LittleEndian, int8(0))
	binary.Write(buffer, binary.LittleEndian, int8(0))
	binary.Write(buffer, binary.LittleEndian, []byte(msg))
	binary.Write(buffer, binary.LittleEndian, int8(0))
	return buffer.Bytes()
}

func JoinRoom(roomid string) []byte {
	// 选择要链接的房间号
	msg := fmt.Sprintf("type@=loginreq/roomid@=%s/", roomid)
	return PostData(msg)
}

func JoinMsg(roomid string) []byte {
	msg := fmt.Sprintf("type@=joingroup/rid@=%s/gid@=-9999/", roomid)
	return PostData(msg)
}

func PreParse(conn net.Conn) (string, error) {
	var header = make([]byte, 12)
	var buffer = make([]byte, BufferSize)
	//var msgLen  int32
	_, err := conn.Read(header)
	if err != nil {
		return "", errors.New("预解析失败")
	}
	conn.Read(buffer)
	return string(buffer), nil
}

func ParseData(conn net.Conn) map[string]interface{} {
	// 解析， 将二进制数据转化为可读的
	Parsed := make(map[string]interface{})
	str, err := PreParse(conn)
	if err != nil {
		// fmt.Println(err)
		return nil
	}

	s := strings.Trim(str, "/")
	items := strings.Split(s, "/")
	for _, str := range items {
		k := strings.SplitN(str, "@=", 2)
		if len(k) > 1 {
			Parsed[k[0]] = k[1]
		}
	}
	return Parsed
}

func PreConn(roomid string) net.Conn {
	buffer := make([]byte, BufferSize)
	JoinData := JoinRoom(roomid)
	JoinMsg := JoinMsg(roomid)
	conn, _ := net.Dial("tcp", ServerAddr)
	_, werr := conn.Write(JoinData)
	if werr != nil {
		fmt.Println(werr)
	}
	_, err := conn.Read(buffer)
	if err != nil {
		fmt.Println(errors.New("无法连接房间 " + err.Error()))
	}
	conn.Write(JoinMsg)
	return conn
}

func CountConnect(roomid string, count *SafeMap, redisC redis.Conn, ch chan string) {
	conn := PreConn(roomid)
	timestamp := time.Now().Unix()
	defer conn.Close()
	for {
		parsed := ParseData(conn) // type: dgb - gift, chatmsg - danmu , uenter - enter
		// nn - nickname  level  tx
		if time.Now().Unix()-timestamp > 21 {
			timestamp = time.Now().Unix()
			_, err := conn.Write(PostData(fmt.Sprintf("type@=keeplive/tick@=%s/", timestamp)))
			if err != nil {
				conn.Close()
				return
			}
		}
		if parsed["type"] == "chatmsg" {
			//fmt.Printf("user: %s  danmu: %s level: %s room: %s \n", parsed["nn"], parsed["txt"], parsed["level"], parsed["rid"])
			key := fmt.Sprintf("%s", parsed["rid"])
			count.add(key)
		}
		if count.readMap("one|timer") < time.Now().Unix() { // 按照一分钟 五分钟 半个小时为维度进行保存
			count.setValue("one|timer", time.Now().Unix()+OneTimer)
			oneMinData(count.Map, redisC)
		}
		if count.readMap("five|timer") < time.Now().Unix() { // 按照一分钟 五分钟 半个小时为维度进行保存
			count.setValue("five|timer", time.Now().Unix()+FiveTimer)
			fiveMinData(count.Map, redisC)
		}
		if count.readMap("half|timer") < time.Now().Unix() { // 按照一分钟 五分钟 半个小时为维度进行保存
			count.setValue("half|timer", time.Now().Unix()+HalfTimer)
			halfHourData(count, redisC)
		}
		if count.Map["restart"] < time.Now().Unix(){
			count.setValue("restart",time.Now().Unix() + GetTaskTimer)
			ch <- ""
		}

	}
}

// 将字典的数据保存进去。。

func oneMinData(mapData map[string]int64, redisC redis.Conn) {
	for k, v := range mapData {
		if k == "one|timer" || k == "five|timer" || k == "half|timer" {
			continue
		}
		key := "one|" + k
		hisKey := "onehis|" + k
		hisdata, _ := redis.Int64(redisC.Do("GET", hisKey))
		if hisdata > v { // 重启后新获取的是会小于旧的累计的数据的， 进行重置
			redisC.Do("SET", key, v) // 保存历史的数据
			redisC.Do("SET", hisKey, v)
		} else {
			redisC.Do("SET", key, v-hisdata) // 用历史保存的和现在的相减
			redisC.Do("SET", hisKey, v)
		}
	}
}

func fiveMinData(mapData map[string]int64, redisC redis.Conn) {
	for k, v := range mapData {
		if k == "one|timer" || k == "five|timer" || k == "half|timer" {
			continue
		}
		key := "five|" + k
		hisKey := "fivehis|" + k
		hisdata, _ := redis.Int64(redisC.Do("GET", hisKey))
		fmt.Println(hisdata, v)
		if hisdata > v { // 重启后新获取的是会小于旧的累计的数据的， 进行重置
			redisC.Do("SET", key, v)
			redisC.Do("SET", hisKey, v)
		} else {
			redisC.Do("SET", key, v-hisdata)
			redisC.Do("SET", hisKey, v)
		}
	}
}

func halfHourData(count *SafeMap, redisC redis.Conn) {
	fmt.Println("half count", count.Map)
	for k, v := range count.Map {
		if k == "one|timer" || k == "five|timer" || k == "half|timer" || k == "restart" {
			continue
		}
		key := "half|" + k
		redisC.Do("SET", key, v)
		count.setValue(k, 0)
	}
}
