/**
websocket客户端
*/
package ChatPlugin

import (
	"MCDaemon-go/lib"
	"fmt"
	"net/http"
	"strconv"
	"sync"

	"github.com/golang/protobuf/proto"
	"golang.org/x/net/websocket"
)

type WSServer struct {
	Port            int
	Suburl          string                     //子路由
	ReceiveMessage  chan *Message              //接受到的消息
	SendMessage     chan *Message              //要发送的消息
	origin          string                     //源地址
	minecraftServer lib.Server                 //服务器实例接口
	ConnPool        map[string]*websocket.Conn //连接池，键为服务器名称，值为websocket连接
	RWPool          *sync.RWMutex              //连接池读写锁
	WhiteList       map[string]interface{}     //白名单
}

func (this *WSServer) handler(conn *websocket.Conn) {
	var err error
	for {
		var reply []byte
		err = websocket.Message.Receive(conn, &reply)
		if err != nil {
			this.minecraftServer.WriteLog("error", fmt.Sprint("聊天服务器出错：", err))
			conn.Close()
			break
		}
		//将proto消息解码
		newMessage := &Message{}
		err = proto.Unmarshal(reply, newMessage)
		if err != nil {
			this.minecraftServer.WriteLog("warn", fmt.Sprint("非法连接：", conn.RemoteAddr().String()))
			conn.Close()
			break
		}
		serverName := newMessage.GetServerName()
		//加入到连接池中,若不在聊天白名单中，则关闭连接
		if ok := this.appendToConnPool(serverName, conn); !ok {
			conn.Close()
			break
		}
		//将消息加入到接收管道中
		this.ReceiveMessage <- newMessage
	}
}

//将websocket连接加入到连接池中
func (this *WSServer) appendToConnPool(serverName string, conn *websocket.Conn) bool {
	if _, ok := this.WhiteList[serverName]; ok {
		//如果没有进入连接池，则加入到连接池中
		if _, ok := this.readPool(serverName); !ok {
			this.writePool(serverName, conn)
		}
		return true
	}
	return false
}

//发送消息
func (this *WSServer) sendMessage() {
	var err error
	for {
		messageObj := <-this.SendMessage
		messageContent := messageObj.GetMessage()
		for serverName, conn := range this.ConnPool {
			//若出现错误，则从连接池中删除并关闭这条连接
			if err = websocket.Message.Send(conn, messageContent); err != nil {
				this.deletePool(serverName)
				conn.Close()
			}
		}
	}
}

func (this *WSServer) readPool(serverName string) (*websocket.Conn, bool) {
	this.RWPool.RLock()
	defer this.RWPool.RUnlock()
	if val, ok := this.ConnPool[serverName]; ok {
		return val, true
	}
	return nil, false
}

func (this *WSServer) writePool(serverName string, conn *websocket.Conn) {
	this.RWPool.Lock()
	defer this.RWPool.Unlock()
	this.ConnPool[serverName] = conn
}

func (this *WSServer) deletePool(serverName string) {
	this.RWPool.Lock()
	defer this.RWPool.Unlock()
	delete(this.ConnPool, serverName)
}

func (this *WSServer) Start() error {
	url := "localhost:" + strconv.Itoa(this.Port)
	http.Handle("/"+this.Suburl, websocket.Handler(this.handler))
	err := http.ListenAndServe(url, nil)
	if err != nil {
		return err
	}
	return nil
}
