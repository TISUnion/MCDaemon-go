package server

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"strings"
	"unicode/utf8"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

//等待服务器加载完地图
func (svr *Server) WaitEndLoading() bool {
	var buffer []byte = make([]byte, 4096)
	var retStr string
	//运行子进程
	svr.run_process()
	fmt.Println("正在加载服务器地图...")
	for {
		n, err := svr.Stdout.Read(buffer)
		retStr = string(buffer[:n])
		if err != nil {
			fmt.Println(svr.name, "服务器启动失败!")
			return false
		}
		if strings.Contains(retStr, "[Server thread/INFO]: Done") {
			fmt.Println(svr.name, "服务器地图加载完成!")
			break
		}
	}
	return true
}

//正式运行MCD
func (svr *Server) Run() {
	var buffer []byte = make([]byte, 4096)
	var buffercmd []byte = make([]byte, 4096)
	cmdin := bufio.NewReader(os.Stdin)
	go func() {
		for {
			n, err := cmdin.Read(buffercmd)
			if err != nil {
				break
			}
			cmdStr := Buffer2String(buffercmd, n)
			if cmdStr[len(cmdStr)-1] == '\n' {
				cmdStr = cmdStr[0 : len(cmdStr)-1]
			}
			svr.Execute(cmdStr)
		}
	}()

	for {
		n, err := svr.Stdout.Read(buffer)
		if err != nil {
			//如果进程已关闭则执行容器关闭函数
			svr.End()
			break
		}
		svrStr := Buffer2String(buffer, n)

		fmt.Println(svr.name, "服务器:", svrStr)
		// 异步处理语法解析器和运行插件
		go svr.RunParsers(svrStr)
	}

}

// Buffer2String :从内存读取字符串
func Buffer2String(buffer []byte, n int) string {
	var retStr string
	//注意在window下minecraft进程间通信使用gbk2312;
	//而在linux下则是utf-8
	switch runtime.GOOS {
	case "darwin", "linux":
		var tempN int
		var tempStr rune
		for i := 0; i < n; {
			tempStr, tempN = utf8.DecodeRune(buffer[i:])
			retStr += fmt.Sprintf("%c", tempStr)
			i += tempN
		}
	case "windows":
		reader := transform.NewReader(bytes.NewReader(buffer[:n]), simplifiedchinese.GBK.NewDecoder())
		retBt, _ := ioutil.ReadAll(reader)
		retStr = string(retBt)
	}
	return retStr
}
