package server

import (
	"MCDaemon-go/config"
	"MCDaemon-go/lib"
	parser "MCDaemon-go/parsers"
	plugin "MCDaemon-go/plugins"
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/exec"
	"strconv"
	"sync"
	"time"

	"github.com/go-ini/ini"
	"github.com/sirupsen/logrus"
)

var err error

type Server struct {
	name              string           //服务器名称
	Stdout            *bufio.Reader    //子进程输出
	Cmd               *exec.Cmd        //子进程实例
	stdin             io.WriteCloser   //用于关闭输入管道
	stdout            io.ReadCloser    //用于关闭输出管道
	lock              sync.Mutex       //输入管道同步锁
	pulginPool        chan interface{} //插件池
	maxRunPlugins     int              //插件最大并发数
	log               *logrus.Logger   //日志文件
	pluginList        plugin.PluginMap //插件列表
	disablePluginList plugin.PluginMap //禁用插件列表
	parserList        []lib.Parser     //语法解析器列表
}

//根据参数初始化服务器
func (svr *Server) Init(name string, argv []string) {
	//创建子进程实例
	svr.Cmd = exec.Command("java", argv...)
	svr.Cmd.Dir = config.Cfg.Section("MCDeamon").Key("server_path").String()
	svr.stdout, err = svr.Cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	//接管子进程输入输出
	svr.Stdout = bufio.NewReader(svr.stdout)
	svr.stdin, err = svr.Cmd.StdinPipe()
	if err != nil {
		log.Fatal(err)
	}
	//初始化插件执行池参数
	svr.maxRunPlugins, _ = strconv.Atoi(config.Cfg.Section("MCDeamon").Key("maxRunPlugins").String())
	svr.pulginPool = make(chan interface{}, svr.maxRunPlugins)

	//初始化日志类
	svr.log = logrus.New()
	svr.log.SetLevel(logrus.DebugLevel)

	//初始化插件列表
	svr.pluginList, svr.disablePluginList = plugin.CreatePluginsList()
	svr.parserList = parser.CreateParserList()
}

//运行子进程
func (svr *Server) run_process() {
	svr.Cmd.Start()
}

//写入日志
func (svr *Server) WriteLog(level string, msg string) {
	logFile, err := os.OpenFile(fmt.Sprintf("logs/%s.log", time.Now().Format("2006-01-02")), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)
	defer logFile.Close()
	if err != nil {
		fmt.Println("日志写入系统发生错误！ 因为", err)
	}
	svr.log.SetOutput(logFile)
	switch level {
	case "debug":
		svr.log.Debug(msg)
	case "info":
		svr.log.Info(msg)
	case "warn":
		svr.log.Warn(msg)
	case "error":
		svr.log.Error(msg)
	case "fatal":
		svr.log.Fatal(msg)
	}
}

//重启服务器
func (svr *Server) Restart() {
	svr.Close()
	//获取所有启动项配置
	MCDeamon := config.GetStartConfig()
	//初始化
	svr.Init(svr.name, MCDeamon)
	//等待加载地图
	svr.WaitEndLoading()
	//正式运行MCD
	svr.Run()
}

//复制一个镜像服务器（用于镜像启动）
func (svr *Server) Clone(name string, Argv []string) lib.Server {

	cloneServer := &Server{}
	// 得到一个可用的端口
	port, _ := func() (string, bool) {
		listener, err := net.Listen("tcp", "127.0.0.1:0")
		if err != nil {
			return "", false
		}
		defer listener.Close()

		addr := listener.Addr().String()
		_, portString, err := net.SplitHostPort(addr)
		if err != nil {
			return "", false
		}

		return portString, true
	}()
	//修改端口
	cfg, _ := ini.Load("back-up/" + name + "/server.properties")
	cfg.Section("").NewKey("server-port", port)
	//初始化
	cloneServer.Init(name, Argv)
	//等待加载地图
	cloneServer.WaitEndLoading()
	//正式运行MCD
	cloneServer.Run()
	return cloneServer
}

//关闭服务器
func (svr *Server) Close() {
	svr.Execute("/stop")
}
