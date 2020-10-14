package socket

import (
	"errors"
	protoutil "github.com/gogo/protobuf/proto"
	"net"
	"xj_web_server/httpserver/wss/proto"
	"xj_web_server/util"
	"strings"
	"sync"
	"time"
)

const (
	// 写超时时间
	writeWait = 10 * time.Second
	msgCont   = 1024
)

type Connection struct {
	tcpConn   net.Conn
	inChan    chan []byte
	outChan   chan []byte
	closeChan chan struct{}

	mutex    sync.Mutex
	isClosed bool
}

func InitConn(address string, connFunc func(conn *Connection)) error {
	tcpConn, err := net.Listen("tcp", address)
	if err == nil {
		for {
			tcpConn, err := tcpConn.Accept()
			if err != nil {
				if strings.Contains(err.Error(), "use of closed network connection") {
					break
				}
				util.Logger.Errorf("tcp accept err:%v", err)
				continue
			}
			var conn *Connection
			conn = &Connection{
				tcpConn:   tcpConn,
				inChan:    make(chan []byte, 0),
				outChan:   make(chan []byte, 0),
				closeChan: make(chan struct{}, 1),
			}
			go conn.readLoop()
			go conn.writeLoop()

			go connFunc(conn)

		}
	}

	return err
}

// TODO 修改创建msg的方法
func (conn *Connection) createMsg(cmd byte, data []byte, uid int32) ([]byte, error) {

	//TODO 根据uid查询相应的公私钥

	//1,组成命令包
	createCmd := util.CreateCmd(cmd, data, uid)

	return createCmd, nil
}

func (conn *Connection) ReadMessage() (data []byte, err error) {
	select {
	case data = <-conn.inChan:
	case <-conn.closeChan:
		err = errors.New("connection is closed")
	}
	return
}

func (conn *Connection) WriteMessage(cmd byte, data []byte, uid int32) (err error) {
	msg, err := conn.createMsg(cmd, data, uid)
	if err != nil {
		err = errors.New("msg create err" + err.Error())
	}
	select {
	case conn.outChan <- msg:
	case <-conn.closeChan:
		err = errors.New("connection is closed")
	}
	return
}

// 非法返回
func (conn *Connection) WriteErr(code int32, errorMsg string) (err error) {
	//心跳
	msg := &proto.Hall_S_Fail{
		ErrorCode: code,
		ErrorMsg:  errorMsg,
	}
	dataMsg, _ := protoutil.Marshal(msg)
	return conn.WriteMessage(0x09, dataMsg, 0)
}

// 获取客户端ip
func (conn *Connection) GetIP() string {
	return conn.tcpConn.RemoteAddr().String()
}

func (conn *Connection) Close() error {

	time.Sleep(50 * time.Millisecond)
	//关闭chan 只能执行一次
	conn.mutex.Lock()
	if !conn.isClosed {
		close(conn.closeChan)
		conn.isClosed = true
	}
	conn.mutex.Unlock()

	// 线程安全,可以重复close
	err := conn.tcpConn.Close()
	return err
}

// 读 协程
func (conn *Connection) readLoop() {
	for {
		var data = make([]byte, msgCont)
		if n, err := conn.tcpConn.Read(data); err != nil {
			goto ERROR
		} else {
			// 阻塞这里,等待inChan 空闲
			select {
			case conn.inChan <- data[:n]:
			}
		}
	}
ERROR:
	conn.Close()
}

//写  协程
func (conn *Connection) writeLoop() {
	for {
		select {
		case data := <-conn.outChan:
			//写超时时间
			if err := conn.tcpConn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
				util.Logger.Errorf("write err:%v", err)
				goto ERROR
			}
			if _, err := conn.tcpConn.Write(data); err != nil {
				util.Logger.Errorf("write err:%v", err)
				goto ERROR
			}
		}

	}
ERROR:
	conn.Close()
}
