package ChatPlugin

import (
	"MCDaemon-go/command"
	"MCDaemon-go/config"
	"MCDaemon-go/container"
	"MCDaemon-go/lib"
	"context"
	"sync"

	"golang.org/x/net/websocket"
)

var (
	WSsvr           *WSServer
	WSrsPool        []WebSocketRS
	once            *sync.Once
	packageChan     chan *msgPackage
	pluginCtx       context.Context
	pluginConcel    context.CancelFunc
	LocalServerName string
)

const ChanMaxSize int = 20
const LocalServerId int = -2

func init() {
	once = &sync.Once{}
	pluginCtx, pluginConcel = context.WithCancel(context.Background())
	packageChan = make(chan *msgPackage, ChanMaxSize)
}

//每个服务器的消息包
type msgPackage struct {
	From int //服务器id
	Msg  *Message
}

type ChatPlugin struct{}

func (this *ChatPlugin) Handle(c *command.Command, s lib.Server) {
	if len(c.Argv) < 1 {
		c.Argv = append(c.Argv, "help")
	}
	switch c.Argv[0] {
	case "start":
		command.Group.AddPlayer("ServersChat", c.Player)
	case "stop":
		command.Group.DelPlayer("ServersChat", c.Player)
	case "chat_xxx_say": //发送消息
		plugincfg := config.GetPluginCfg(false)
		contentColor := plugincfg.Section("LinkChat").Key("server_content_color").String()
		playerColor := plugincfg.Section("LinkChat").Key("server_player_color").String()
		serverColor := plugincfg.Section("LinkChat").Key("server_name_color").String()
		msg := &Message{
			ServerName:      &LocalServerName,
			Player:          &(c.Player),
			Message:         &(c.Argv[1]),
			ServerNameColor: &serverColor,
			PlayerColor:     &playerColor,
			MessageColor:    &contentColor,
		}
		packageChan <- &msgPackage{
			From: LocalServerId,
			Msg:  msg,
		}
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
	plugincfg := config.GetPluginCfg(true)
	isStart, _ := plugincfg.Section("LinkChat").Key("isStart").Bool()
	whitelistArr := plugincfg.Section("LinkChat.whitelist").Key("whitelist").ValueWithShadows()
	whitelist := make(map[string]interface{})
	LocalServerName = plugincfg.Section("LinkChat").Key("server_name").String()
	for _, v := range whitelistArr {
		whitelist[v] = 1
	}
	//开启websocket服务器
	if isStart {
		//创建服务器实例
		port, _ := plugincfg.Section("LinkChat").Key("server_port").Int()
		WSSvrCtx, cancel := context.WithCancel(pluginCtx)
		WSsvr = &WSServer{
			ServerId:       -1,
			Port:           port,
			Suburl:         plugincfg.Section("LinkChat").Key("server_sub_url").String(),
			ReceiveMessage: make(chan *Message, ChanMaxSize),
			SendMessage:    make(chan *Message, ChanMaxSize),
			ConnPool:       make(map[string]*websocket.Conn),
			RWPool:         &sync.RWMutex{},
			WhiteList:      whitelist,
			ServerName:     LocalServerName,
			Ctx:            WSSvrCtx,
			Cancel:         cancel,
		}

		go WSsvr.Start()   //开启服务器
		go WSsvr.SendJob() //开启发送消息任务
		//加入websocket读发池
		WSrsPool = append(WSrsPool, WSsvr)
	}

	//连接websocket服务器
	Servername := plugincfg.Section("LinkChat.server").KeyStrings()
	for k, v := range Servername {
		WSCliCtx, cancel := context.WithCancel(pluginCtx)
		WSCli := &WSClient{
			ServerId:       k,
			ServerName:     v,
			addr:           plugincfg.Section("LinkChat.server").Key(v).String(),
			origin:         "TISMCDGO://" + plugincfg.Section("LinkChat.server").Key("server_name").String(),
			ReceiveMessage: make(chan *Message, ChanMaxSize),
			SendMessage:    make(chan *Message, ChanMaxSize),
			Ctx:            WSCliCtx,
			Cancel:         cancel,
		}
		go WSCli.Start() //连接服务器
		WSrsPool = append(WSrsPool, WSCli)
	}
}

/**
多协程读取websocket读发池中收到的消息
*/
func read() {
	for _, rs := range WSrsPool {
		go rs.Read(packageChan)
	}
}

/**
向网络服务器发送消息
*/
func sendNetServer(pkg *msgPackage) {
	for _, rs := range WSrsPool {
		//给除接受消息服务器的其他服务器发送
		if rs.GetId() != pkg.From {
			serverName := rs.GetName()
			pkg.Msg.ServerName = &serverName
			go rs.Send(pkg.Msg)
		}
	}
}

/**
向本机游戏服务器发送消息
*/
func sendLocalServer(pkg *msgPackage) {
	if pkg.From != LocalServerId {
		serverPool := container.GetInstance().Servers
		for _, server := range serverPool {
			// 发送给用户组消息
			players := command.Group.GetPlayer()["ServersChat"]
			for _, player := range players {
				server.Tell(player,
					command.Text{"[" + pkg.Msg.GetServerName() + "]", pkg.Msg.GetServerNameColor()},
					command.Text{pkg.Msg.GetPlayer(), pkg.Msg.GetPlayerColor()},
					command.Text{":", "white"},
					command.Text{pkg.Msg.GetMessage(), pkg.Msg.GetMessageColor()},
				)
			}
		}
	}
}
