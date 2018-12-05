package server

import (
	parser "MCDaemon-go/parsers"
	plugin "MCDaemon-go/plugins"
	"fmt"
	"io"
	"log"
	"strings"
)

//等待服务器加载完地图
func (svr *Server) WaitEndLoading() {
	var buffer []byte = make([]byte, 4096)
	var retStr string
	//运行子进程
	svr.run_process()
	fmt.Println("正在加载服务器地图...")
	for {
		n, err := svr.Stdout.Read(buffer)
		if err != nil {
			if err == io.EOF {
				log.Fatalln("子进程标准输出出错")
			} else {
				log.Fatalln("获取标准输出出错")
			}
			break
		}
		retStr = string(buffer[:n])
		if strings.Contains(retStr, "[Server thread/INFO]: Done") {
			fmt.Println("服务器地图加载完成!")
			break
		}
	}
}

//正式运行MCD
func (svr *Server) Run() {
	var buffer []byte = make([]byte, 4096)
	var retStr string
	for {
		n, err := svr.Stdout.Read(buffer)
		if err != nil {
			if err == io.EOF {
				log.Fatalln("子进程标准输出出错")
			} else {
				log.Fatalln("获取标准输出出错")
			}
			break
		}
		retStr = string(buffer[:n])
		fmt.Println(retStr)
		for _, val := range parser.ParseList {
			command, ok := val.Parsing(retStr)
			if ok {
				plugin.PluginsList[command.Cmd].Handle(command, svr)
			}
		}
	}
}
