package base

import (
	"fmt"
	"net"
	"encoding/binary"
	"bytes"
	"errors"
	"strings"
	"time"
	"github.com/garyburd/redigo/redis"
)

// 主要的处理逻辑

const (
	BufferSize  = 1024
	ServerAddr  = "openbarrage.douyutv.com:8601"
	PostCode = 689
	PullCode = 690
	wtf = "asd"
)

func PostData(msg string) []byte {
	// 构造需要发送的二进制数据
	length := 9+len(msg) // 长度4字节 + 类型2字节 + 加密字段1字节 + 保留字段1字节 + 结尾字段1字节
	buffer := bytes.NewBuffer([]byte{})
	binary.Write(buffer, binary.LittleEndian, int32(length))
	binary.Write(buffer, binary.LittleEndian, int32(length))
	binary.Write(buffer, binary.LittleEndian, int16(PostCode))
	binary.Write(buffer, binary.LittleEndian, int8( 0))
	binary.Write(buffer, binary.LittleEndian, int8(0))
	binary.Write(buffer, binary.LittleEndian, []byte(msg))
	binary.Write(buffer, binary.LittleEndian, int8(0))
	return buffer.Bytes()
}

func JoinRoom(roomid string)[]byte  {
	// 选择要链接的房间号
	msg := fmt.Sprintf("type@=loginreq/roomid@=%s/", roomid)
	return PostData(msg)
}

func JoinMsg(roomid string)[]byte{
	msg := fmt.Sprintf("type@=joingroup/rid@=%s/gid@=-9999/", roomid)
	return PostData(msg)
}


func PreParse(conn net.Conn) (string, error){
	var header = make([]byte, 12)
	var buffer = make([]byte, BufferSize)
	//var msgLen  int32
	_, err := conn.Read(header)
	if err != nil{
		return "", errors.New("预解析失败")
	}
	conn.Read(buffer)
	return string(buffer), nil
}


func ParseData(conn net.Conn) map[string]interface{} {
	// 解析， 将二进制数据转化为可读的
	Parsed := make(map[string]interface{})
	str, err := PreParse(conn)
	if err != nil{
		fmt.Println(err)
	}

	s := strings.Trim(str, "/")
	items := strings.Split(s, "/")
	for _, str := range items {
		k := strings.SplitN(str, "@=", 2)
		if len(k) >1{
			Parsed[k[0]] = k[1]
		}
	}
	return Parsed
}

func PreConn(roomid string) net.Conn  {
	buffer := make([]byte, BufferSize)
	JoinData := JoinRoom(roomid)
	JoinMsg := JoinMsg(roomid)
	conn, _ := net.Dial("tcp", ServerAddr)
	_, werr := conn.Write(JoinData)
	if werr != nil{
		fmt.Println(werr)
	}
	_, err := conn.Read(buffer)
	if err != nil {
		fmt.Println(errors.New("无法连接房间 " + err.Error()))
	}
	conn.Write(JoinMsg)
	return conn
}


func Connect(roomid string)  {
	conn := PreConn(roomid)
	timestamp := time.Now().Unix()
	for  {
		parsed := ParseData(conn) // type: dgb - gift, chatmsg - danmu , uenter - enter
		// nn - nickname  level  txt
		if time.Now().Unix() - timestamp > 21{
			timestamp = time.Now().Unix()
			_, err := conn.Write(PostData(fmt.Sprintf("type@=keeplive/tick@=%s/", timestamp)))
			if err != nil{
				fmt.Println("心跳失败")
			}
		}
		if parsed["type"] == "chatmsg"{
			fmt.Printf("user: %s  danmu: %s level: %s room: %s \n", parsed["nn"], parsed["txt"], parsed["level"], parsed["rid"])
		}
	}
	conn.Close()
}

func CountConnect(roomid string, count *SafeMap, redisC redis.Conn)  {
	conn := PreConn(roomid)
	timestamp := time.Now().Unix()
	for  {
		parsed := ParseData(conn) // type: dgb - gift, chatmsg - danmu , uenter - enter
		// nn - nickname  level  txt
		if time.Now().Unix() - timestamp > 21{
			timestamp = time.Now().Unix()
			_, err := conn.Write(PostData(fmt.Sprintf("type@=keeplive/tick@=%s/", timestamp)))
			if err != nil{
				fmt.Println("心跳失败")
			}
		}
		if parsed["type"] == "chatmsg"{
			//fmt.Printf("user: %s  danmu: %s level: %s room: %s \n", parsed["nn"], parsed["txt"], parsed["level"], parsed["rid"])
			key := fmt.Sprintf("%s", parsed["rid"])
			count.add(key)
		}

		if count.readMap("one|timer") < time.Now().Unix(){ // 按照一分钟 五分钟 半个小时为维度进行保存
			count.setValue("one|timer", time.Now().Unix()+60)
			oneMinData(count.Map, redisC)
		}
		if count.readMap("five|timer") < time.Now().Unix(){ // 按照一分钟 五分钟 半个小时为维度进行保存
			count.setValue("five|timer", time.Now().Unix()+60*5)
			fiveMinData(count.Map, redisC)
		}

	}
	conn.Close()
}

// 将字典的数据保存进去。。


func oneMinData(mapData map[string]int64 ,redisC redis.Conn)  {
	fmt.Println("one count", mapData)
	for k, v := range mapData{
		key := "one|"+k
		u, _ := redis.Int64(redisC.Do("GET", key))
		if v < u{ // fix bug. mapdata中的数据会重新开始增加，而redis中的不会。
			redisC.Do("SET", key, v)
		}else {
			redisC.Do("SET", key, v-u)
		}

	}
}

func fiveMinData(mapData map[string]int64 ,redisC redis.Conn)  {
	fmt.Println("five count", mapData)
	for k, v := range mapData{
		key := "five|"+k
		u, _ := redis.Int64(redisC.Do("GET", key))
		if v < u{
			redisC.Do("SET", key, v)
		}else {
			redisC.Do("SET", key, v-u)
		}
	}
}