package ChatPlugin

import (
	"MCDaemon-go/lib"
	"context"
	"fmt"

	"github.com/golang/protobuf/proto"
	"golang.org/x/net/websocket"
)

type WSClient struct {
	ServerId       int    //服务器id
	ServerName     string //服务器名称
	addr           string
	origin         string
	ws             *websocket.Conn //websocket连接
	ReceiveMessage chan *Message   //接受到的消息
	Ctx            context.Context //上下文
	Cancel         context.CancelFunc
}

func (this *WSClient) Start() error {
	var err error

	this.ws, err = websocket.Dial("ws://"+this.addr, "", this.origin)
	if err != nil {
		lib.WriteDevelopLog("error", err.Error())
		return err
	}
	defer this.ws.Close()
	defer lib.WriteDevelopLog("error", "连接")
	this.Send(&Message{
		ServerName: &LocalServerName,
		State:      &FirstTouch,
	})
	for {
		msg := make([]byte, 5096)
		slen, err := this.ws.Read(msg) //此处阻塞，等待有数据可读
		msg = msg[:slen]
		if err != nil {
			lib.WriteDevelopLog("error", fmt.Sprint("读取错误：", err))
			//如果连接出错，则释放连接
			break
		}
		newMessage := &Message{}
		err = proto.Unmarshal(msg, newMessage)
		if err != nil {
			lib.WriteDevelopLog("error", fmt.Sprint("解码错误:", err, "内容：", msg))
			break
		}
		if newMessage.GetState() != 0 {
			lib.WriteDevelopLog("error", "聊天服务器连接失败：不再白名单内！")
			break
		}
		this.ReceiveMessage <- newMessage
	}
	this.Cancel()
	return err
}

func (this *WSClient) Send(msg *Message) {
	data, err := proto.Marshal(msg)
	if err != nil {
		lib.WriteDevelopLog("error", fmt.Sprint("加密错误：", err))
		return
	}
	err = websocket.Message.Send(this.ws, data)
	if err != nil {
		lib.WriteDevelopLog("error", fmt.Sprint(this.GetName, "发送信息错误：", err))
	}
}

func (this *WSClient) Read() {
	for {
		select {
		case <-this.Ctx.Done():
			return
		case msg := <-this.ReceiveMessage:
			packageChan <- &msgPackage{
				From: this.ServerId,
				Msg:  msg,
			}
		}
	}
}

func (this *WSClient) GetId() int {
	return this.ServerId
}

func (this *WSClient) GetName() string {
	return this.ServerName
}

func (this *WSClient) IsAlive() bool {
	if this.ws != nil {
		return true
	}
	this.Cancel()
	return false
}
