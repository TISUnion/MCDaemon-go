package ChatPlugin

import (
	"MCDaemon-go/command"
	"MCDaemon-go/config"
	"MCDaemon-go/lib"
	"sync"

	"github.com/go-ini/ini"
	"golang.org/x/net/websocket"
)

var (
	plugincfg *ini.File
	WSsvr     *WSServer
	WSrsPool  []WebSocketRS
)

const ChanMaxSize int = 20

func init() {
	plugincfg = config.GetPluginCfg(true)
	isStart, _ := plugincfg.Section("LinkChat").Key("isStart").Bool()
	whitelistArr := plugincfg.Section("LinkChat.whitelist").Key("whitelist").ValueWithShadows()
	whitelist := make(map[string]interface{})
	for _, v := range whitelistArr {
		whitelist[v] = 1
	}
	if isStart {
		//创建服务器实例
		port, _ := plugincfg.Section("LinkChat").Key("server_port").Int()
		WSsvr = &WSServer{
			Port:           port,
			Suburl:         plugincfg.Section("LinkChat").Key("server_sub_url").String(),
			ReceiveMessage: make(chan *Message, ChanMaxSize),
			SendMessage:    make(chan *Message, ChanMaxSize),
			ConnPool:       make(map[string]*websocket.Conn),
			RWPool:         &sync.RWMutex{},
			WhiteList:      whitelist,
		}
		WSsvr.Start()
		//加入websocket读发池
		WSrsPool = append(WSrsPool, WSsvr)
	}
}

type ChatPlugin struct{}

func (this *ChatPlugin) Handle(c *command.Command, s lib.Server) {
}

func (this *ChatPlugin) Init(s lib.Server) {
	WSsvr.minecraftServer = s
}

func (this *ChatPlugin) Close() {
}
