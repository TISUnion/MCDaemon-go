package server

import (
	"MCDaemon-go/config"
	"MCDaemon-go/container"
	"MCDaemon-go/lib"
	parser "MCDaemon-go/parsers"
	plugin "MCDaemon-go/plugins"
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os/exec"
	"strconv"
	"sync"
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
	pluginList        plugin.PluginMap //插件列表
	disablePluginList plugin.PluginMap //禁用插件列表
	parserList        []lib.Parser     //语法解析器列表
	port              string           //启动服务器端口
	unqiueLock        sync.Mutex       //堵塞插件执行池锁
}

//根据参数初始化服务器
func (svr *Server) Init(name string, argv []string, workDir string) {
	svr.name = name
	//创建子进程实例
	svr.Cmd = exec.Command("java", argv...)
	svr.Cmd.Dir = workDir
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

	//初始化插件列表
	svr.pluginList, svr.disablePluginList = plugin.CreatePluginsList(false)
	svr.parserList = parser.CreateParserList()
	//执行插件init
	for _, v := range svr.pluginList {
		v.Init(svr)
	}
	//设置端口
	if svr.port == "" {
		svr.port = "25565"
	}
}

//运行子进程
func (svr *Server) run_process() {
	svr.Cmd.Start()
}

func (svr *Server) Getinfo() string {
	return fmt.Sprintf("镜像名称：%s\\n,端口：%s\\n", svr.name, svr.port)
}

//写入日志
func (svr *Server) WriteLog(level string, msg string) {
	lib.WriteRuntimeLog(level, msg, svr.name)
}

//重启服务器
func (svr *Server) Restart() {
	c := container.GetInstance()
	c.Group.Add(1)
	//关闭
	c.Del(svr.name)
	//启动
	workDir := svr.Cmd.Dir
	c.Add(svr.name, workDir, svr)
	c.Group.Done()
}

//启动服务器
func (svr *Server) Start(name string, Argv []string, workDir string) {
	//初始化
	svr.Init(name, Argv, workDir)
	//等待加载地图
	if svr.WaitEndLoading() {
		//正式运行MCD
		svr.Run()
	} else {
		//没加载成功就释放同步锁
		c := container.GetInstance()
		c.Group.Done()
	}
}

//重新读取配置文件
func (svr *Server) ReloadConf() {
	svr.pluginList, svr.disablePluginList = plugin.CreatePluginsList(true)
}

//复制一个镜像服务器（用于镜像启动）
func (svr *Server) Clone() lib.Server {
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
	cloneServer.port = port
	return cloneServer
}

//获取端口
func (svr *Server) GetPort() string {
	return svr.port
}

//获取服务器实例名称//获取端口
func (svr *Server) GetName() string {
	return svr.name
}

//以容器形式关闭服务器
func (svr *Server) CloseInContainer() {
	c := container.GetInstance()
	//关闭
	c.Del(svr.name)
}

//关闭服务器
func (svr *Server) Close() {
	// 关闭插件
	svr.RunPluginClose()
	svr.Execute("/stop")
}

//在容器中注销该服务器
func (svr *Server) End() {
	// 关闭插件
	svr.RunPluginClose()
	//释放同步锁
	c := container.GetInstance()
	c.Group.Done()
}
