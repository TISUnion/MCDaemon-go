package ChatPlugin

import (
	"MCDaemon-go/lib"

	"github.com/golang/protobuf/proto"
	"golang.org/x/net/websocket"
)

type WSClient struct {
	addr            string
	origin          string
	ws              *websocket.Conn //websocket连接
	minecraftServer lib.Server      //服务器实例接口
	ReceiveMessage  chan *Message   //接受到的消息
	SendMessage     chan *Message   //要发送的消息
}

func (this *WSClient) Start() error {
	defer this.ws.Close()
	var err error
	this.ws, err = websocket.Dial(this.addr, "", this.origin)
	if err != nil {
		return err
	}
	for {
		msg := make([]byte, 5096)
		_, err = this.ws.Read(msg[:]) //此处阻塞，等待有数据可读
		if err != nil {
			//如果连接出错，则释放连接
			break
		}
		newMessage := &Message{}
		err = proto.Unmarshal(msg, newMessage)
		if err != nil {
			break
		}
		this.ReceiveMessage <- newMessage
	}
	return err
}

func (this *WSClient) Send(msg *Message) {
	data, err := proto.Marshal(msg)
	if err != nil {
		return
	}
	websocket.Message.Send(this.ws, data)
}

func (this *WSClient) Read() *Message {
	return <-this.ReceiveMessage
}
