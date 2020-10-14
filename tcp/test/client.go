package main

import (
	"encoding/json"
	"fmt"
	protoutil "github.com/gogo/protobuf/proto"
	"io"
	"log"
	"net"
	"xj_web_server/httpserver/wss/proto"
	"xj_web_server/util"
	"time"
)

const (
	//addr = "test1.baoxinton.com:8000"
	//addr = "192.168.1.8:8000"
	//addr = "47.107.188.43:13001"
	addr = "127.0.0.1:13001"
)

func main() {
	tcpClient()
}
func tcpClient() {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		fmt.Println("连接服务端失败:", err.Error())
		return
	}
	fmt.Println("已连接服务器")
	defer conn.Close()
	go sender(conn)
	read(conn)

}

func sender(conn net.Conn) {
	for {
		//words := "{\"token\":1,\"name\":\"golang\",\"Message\":\"message\"}"
		connMsgStr := &proto.Hall_C_Msg{
			Token: "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiIxMzg1ODEiLCJuYmYiOjE1NzkyNDQ1MTB9.KxaDwRMnt9NlVRuge5huQc1fuYpEqZfiNvZHhrg1z4GQwBw2J83CaCQhVAxIrhWBTnXCMlK_K_v3qSXBcVkGejlKT-ZS7DogQHexE7qE60OQWMlRmRElsrtvVxIxl-7uU9vQ1N6LCwZ8xv3YqHrQfiTUcPzohlPhXjB9_JznP_S3OOzsoHxvnmq5GD6xeGrO219FON2_L6KdljqXMLbChA2HYVRGf3fMeHu_PPLjCWAn0TiOnrE18CSKeejr-7kKJKZDO3QGuIohi_dUtCjAX-Uok2NXYrjg-sNdTgR8GgaELIk3AiNZ4-rqTiqvdLA-7ME1tRW_gx09_DF17qkV5A",
		}

		dataMsg, _ := protoutil.Marshal(connMsgStr)
		_, err := conn.Write(util.CreateCmd(0x00, dataMsg, 138581))
		if err != nil {
			log.Println("write:", err)
			return
		}
		break
		time.Sleep(10 * time.Second)
		//log.Println("写数据:", util.CreateCmd(0x00, dataMsg, 100), n)
	}
}

func read(conn net.Conn) {
	for {
		var message = make([]byte, 1024)
		n, err := conn.Read(message)
		if err != nil && err != io.EOF || len(message) == 0 {
			log.Println("read:", err)
			return
		}
		message = message[:n]
		//校验
		if !util.IsCheckCmd(message) {
			break
		}
		var msg proto.Hall_S_Msg
		err = protoutil.Unmarshal(message[8:len(message)-2], &msg)
		bytes, _ := json.Marshal(msg)
		log.Print("recvString:\n", string(bytes), msg.AnnouncementList)
	}
}
