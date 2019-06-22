package ChatPlugin

import (
	"MCDaemon-go/command"
	"MCDaemon-go/config"
	"MCDaemon-go/container"
	"MCDaemon-go/lib"
	"context"
	"fmt"
	"strings"
	"sync"
)

var (
	WSsvr           *WSServer
	WSrsPool        []WebSocketRS
	once            *sync.Once
	packageChan     chan *msgPackage
	pluginCtx       context.Context
	pluginConcel    context.CancelFunc
	LocalServerName string
	IsStart         bool
	FirstTouch      int64 = 1024
	NotInWhitelist  int64 = 1
	ContainerName   string
)

const ChanMaxSize int = 20
const LocalServerId int = -1
const MinecraftServerId int = -2

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
			From: MinecraftServerId,
			Msg:  msg,
		}
	default:
		text := "!!Chat start 开启跨服聊天模式\\n!!Chat stop 关闭跨服聊天模式"
		s.Tell(c.Player, text)
	}
}

func (this *ChatPlugin) Init(s lib.Server) {
	once.Do(func() {
		ContainerName = s.GetName()
		//开启服务器和连接服务器
		start()
		if IsStart {
			WSsvr.minecraftServer = s
		}
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
	IsStart, _ = plugincfg.Section("LinkChat").Key("is_start").Bool()
	whitelistArr := plugincfg.Section("LinkChat.whitelist").Key("whitelist").ValueWithShadows()
	whitelist := make(map[string]interface{})
	LocalServerName = plugincfg.Section("LinkChat").Key("server_name").String()

	for _, v := range whitelistArr {
		whitelist[v] = 1
	}
	//开启websocket服务器
	if IsStart {
		//创建服务器实例
		port, _ := plugincfg.Section("LinkChat").Key("server_port").Int()
		WSSvrCtx, cancel := context.WithCancel(pluginCtx)
		WSsvr = &WSServer{
			ServerId:       LocalServerId,
			Port:           port,
			Suburl:         plugincfg.Section("LinkChat").Key("server_sub_url").String(),
			ReceiveMessage: make(chan *msgPackage, ChanMaxSize),
			SendMessage:    make(chan *Message, ChanMaxSize),
			RWPool:         &sync.RWMutex{},
			WhiteList:      whitelist,
			ServerName:     LocalServerName,
			Ctx:            WSSvrCtx,
			Cancel:         cancel,
			Alive:          true,
		}
		go WSsvr.Start() //开启服务器
		//加入websocket读写池
		WSrsPool = append(WSrsPool, WSsvr)
	}

	//连接websocket服务器
	Servername := plugincfg.Section("LinkChat.server").KeyStrings()
	clientID := 0
	for _, v := range Servername {
		clientID++
		WSCliCtx, cancel := context.WithCancel(pluginCtx)
		WSCli := &WSClient{
			ServerId:       clientID,
			ServerName:     LocalServerName,
			addr:           plugincfg.Section("LinkChat.server").Key(v).String(),
			origin:         "TISMCDGO://" + plugincfg.Section("LinkChat").Key("server_name").String(),
			ReceiveMessage: make(chan *Message, ChanMaxSize),
			Ctx:            WSCliCtx,
			Cancel:         cancel,
		}
		//发送初探消息

		go WSCli.Start() //连接服务器
		WSrsPool = append(WSrsPool, WSCli)
	}

	//开启消息发送
	go SendMessage()
}

/**
多协程读取websocket读发池中收到的消息
*/
func read() {
	if WSrsPool != nil {
		for _, rs := range WSrsPool {
			go rs.Read()
		}
	}
}

/**
向网络服务器发送消息
*/
func sendNetServer(pkg *msgPackage) {
	var releaseIndexs []int
	if WSrsPool != nil {
		for k, rs := range WSrsPool {
			//死链接，加入待释放池
			if !rs.IsAlive() {
				releaseIndexs = append(releaseIndexs, k)
				continue
			}
			if rs.GetId() != pkg.From {
				// serverName := rs.GetName()
				// pkg.Msg.ServerName = &serverName
				go rs.Send(pkg.Msg)
			}
		}
	}
	//删除死链接
	if len(releaseIndexs) != 0 {
		for _, index := range releaseIndexs {
			lib.WriteDevelopLog("info", fmt.Sprint(WSrsPool[index].GetName(), "已死亡"))
			WSrsPool = append(WSrsPool[:index], WSrsPool[index+1:]...)
		}
	}
}

/**
向本机游戏服务器发送消息
*/
func sendLocalServer(pkg *msgPackage) {
	//是本机发送就不发
	if WSrsPool != nil {
		serverPool := container.GetInstance().Servers
		for serverName, server := range serverPool {
			if pkg.From == MinecraftServerId && ContainerName == serverName {
				continue
			}
			// 发送给用户组消息
			players := command.Group.GetPlayer()["ServersChat"]
			for _, player := range players {
				server.Tell(player,
					command.Text{"[" + pkg.Msg.GetServerName() + "]", pkg.Msg.GetServerNameColor()},
					command.Text{pkg.Msg.GetPlayer(), pkg.Msg.GetPlayerColor()},
					command.Text{":", "white"},
					command.Text{strings.Replace(pkg.Msg.GetMessage(), "\r", "", -1), pkg.Msg.GetMessageColor()},
				)
			}
		}
	}
}

/**
接受消息并发送
*/
func SendMessage() {
	for {
		select {
		case <-pluginCtx.Done():
			lib.WriteDevelopLog("warn", "消息发送上下文已关闭")
			return
		case message := <-packageChan:
			sendNetServer(message)
			sendLocalServer(message)
		}
	}
}
