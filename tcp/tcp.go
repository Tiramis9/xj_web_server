package tcp

import (
	"fmt"
	protoutil "github.com/gogo/protobuf/proto"
	"log"
	"net/http"
	"xj_web_server/config"
	"xj_web_server/httpserver/wss/proto"
	"xj_web_server/model"
	"xj_web_server/socket/tcp"
	"xj_web_server/tcp/hander"
	tcpModel "xj_web_server/tcp/model"
	"xj_web_server/util"
	"xj_web_server/util/uuid"
	"time"
)



func Run(address string) {
	err := socket.InitConn(address, func(conn *socket.Connection) {
		defer conn.Close()
		client(conn)
	})
	if err != nil {
		log.Fatalf("tcp err:%v", err)
	}
}

func client(conn *socket.Connection) {
	//TODO 服务器判断
	currentClient := &model.TCPClient{
		Id:          uuid.NewUUID().Hex32(),
		Token:       "",
		Ip:          conn.GetIP(),
		Conn:        conn,
		MessageChan: make(chan interface{}, 0),
		IsOnLine:    true,
		OnLineTime:  time.Now().Format("2006-01-02T15:04:05.000"),
		OfflineTime: "",
	}
	go heartBeat(currentClient)
	readMsg(currentClient)
}

// 读客户端消息
func readMsg(currentClient *model.TCPClient) {
	defer currentClient.ReleaseClient()
	var errMsg = ""
	for {
		var (
			data []byte
			err  error
		)
		data, err = currentClient.Conn.ReadMessage()
		// TODO  客服端数据 要处理
		//fmt.Println(len(data),string(data), err)
		if err != nil {
			// 连接错误
			util.Logger.Errorf("读取错误：err:%v", err)
			goto ERROR
		}
		cmd := util.IsCheckCmd(data)
		if !cmd {
			util.Logger.Errorf("数据包错误：err:%v", err)
			//TODO  返回错误到客户端
			goto ERROR
		}

		// uid
		var cmdUid int32
		err = util.BytesToInt(&cmdUid, data[4:10])
		if err != nil {
			util.Logger.Errorf("数据包错误：err:%v", err)
			//TODO  返回错误到客户端
			goto ERROR
		}
		// 解析命令
		var msg proto.Hall_C_Msg
		err = protoutil.Unmarshal(data[8:len(data)-2], &msg)
		//log.Print("服务端recvMsg:\n", msg, string(msg.Data))
		//log.Print("服务端recvMsg 命令号:\n", data[1])
		//log.Print("服务端recvMsg uid：\n", cmdUid)
		// 命令号
		switch data[1] {
		case 0:
			token, err := hander.Auth(msg, cmdUid)
			if err != nil {
				util.Logger.Errorf("auth err ：%v", err)
				errMsg = err.Error()
				// 回复错误
				goto ERROR
			}
			currentClient.Token = token
			currentClient.Uid = cmdUid
			if _, ok := model.ClientList.Load(currentClient.Id); !ok {

				//添加客户端
				model.ClientList.Store(currentClient.Token, currentClient)
			}
			hellMsg, fail := GetHellMsg(cmdUid)
			if fail != nil {
				goto ERROR
			}
			dataMsg, _ := protoutil.Marshal(hellMsg)
			err = currentClient.Conn.WriteMessage(0x00, dataMsg, cmdUid)
			if err != nil {
				util.Logger.Errorf("回复错误：%v", err)
				errMsg = err.Error()
				// 回复错误
				goto ERROR
			}
		default:
			goto ERROR
		}
	}

ERROR:
	currentClient.Conn.WriteErr(http.StatusForbidden, errMsg)
}

// 心跳包
func heartBeat(currentClient *model.TCPClient) {
	ticker := time.NewTicker(time.Duration(config.GetWss().HeartbeatTime) * time.Second)
	defer ticker.Stop()
	defer currentClient.ReleaseClient()
	for {
		select {
		case <-ticker.C:
			//心跳
			msg, _ := GetHellMsg(currentClient.Uid)
			dataMsg, _ := protoutil.Marshal(msg)
			err := currentClient.Conn.WriteMessage(0x00, dataMsg, currentClient.Uid)
			if err != nil {
				fmt.Println(currentClient.Conn.GetIP(), err)
				// 某个客户端异常退出
				return
			}
		}
	}
}

func GetHellMsg(uid int32) (*proto.Hall_S_Msg, error) {

	return tcpModel.GetUerByUid(uid)
}
