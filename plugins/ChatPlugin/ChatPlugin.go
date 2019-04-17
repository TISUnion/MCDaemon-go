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
	once      *sync.Once
)

const ChanMaxSize int = 20

func init() {
	once = &sync.Once{}
}

//每个服务器的消息包
type msgPackage struct {
	From string
	Msg  *Message
}

type ChatPlugin struct{}

func (this *ChatPlugin) Handle(c *command.Command, s lib.Server) {
	if len(c.Argv) < 1 {
		c.Argv = append(c.Argv, "help")
	}
	switch c.Argv[0] {
	case "start":
	case "stop":
	default:
		text := "!!Chat start 开启跨服聊天模式\\n!!Chat stop 关闭跨服聊天模式"
		s.Tell(c.Player, text)
	}
}

func (this *ChatPlugin) Init(s lib.Server) {

	once.Do(func() {
		//开启服务器和连接服务器
		start()
		WSsvr.minecraftServer = s
		//读取消息并发送到本地服务器和镜像内
		read()
	})
}

func (this *ChatPlugin) Close() {
}

/**
开启和连接websocket服务器
*/
func start() {
	plugincfg = config.GetPluginCfg(true)
	isStart, _ := plugincfg.Section("LinkChat").Key("isStart").Bool()
	whitelistArr := plugincfg.Section("LinkChat.whitelist").Key("whitelist").ValueWithShadows()
	whitelist := make(map[string]interface{})
	for _, v := range whitelistArr {
		whitelist[v] = 1
	}
	//开启websocket服务器
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
			ServerName:     plugincfg.Section("LinkChat").Key("server_name").String(),
		}
		go WSsvr.Start()
		//加入websocket读发池
		WSrsPool = append(WSrsPool, WSsvr)
	}

	//连接websocket服务器
	Servername := plugincfg.Section("LinkChat.server").KeyStrings()
	for _, v := range Servername {
		WSCli := &WSClient{
			ServerName:     v,
			addr:           plugincfg.Section("LinkChat.server").Key(v).String(),
			origin:         "TISMCDGO://" + plugincfg.Section("LinkChat.server").Key("server_name").String(),
			ReceiveMessage: make(chan *Message, ChanMaxSize),
			SendMessage:    make(chan *Message, ChanMaxSize),
		}
		go WSCli.Start()
		WSrsPool = append(WSrsPool, WSCli)
	}
}

/**
多协程读取websocket读发池中收到的消息
*/
func read() {
}

/**
向其他服务器发送消息
*/
func sendNetServer(msg *Message) {
}

/**
向本机游戏服务器发送消息
*/
func sendLocalServer(msg *Message) {

}
