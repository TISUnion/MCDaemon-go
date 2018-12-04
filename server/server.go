package server

import (
	"bufio"
	"io"
	"log"
	"os/exec"
)

//单例模式
var (
	is_set bool
	svr    *Server
	err    error
)

type Server struct {
	Stdout *bufio.Reader  //子进程输出
	Stdin  *bufio.Writer  //子进程输入
	Cmd    *exec.Cmd      //子进程实例
	stdin  io.WriteCloser //用于关闭输入管道
	stdout io.ReadCloser  //用于关闭输出管道
}

func init() {
	is_set = false
}

//获取实例
func GetServerInstance() *Server {
	if !is_set {
		svr = &Server{}
	}
	return svr
}

//根据参数初始化服务器
func (svr *Server) Init(argv []string) {
	svr.Cmd = exec.Command("java", argv...)
	svr.stdout, err = svr.Cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	svr.Stdout = bufio.NewReader(svr.stdout)
	svr.stdin, err = svr.Cmd.StdinPipe()
	if err != nil {
		log.Fatal(err)
	}
	svr.Stdin = bufio.NewWriter(svr.stdin)
}

//运行子进程
func (svr *Server) run_process() {
	svr.Cmd.Start()
}

//关闭服务器
func (svr *Server) Close() {
	svr.stdin.Close()
	svr.stdout.Close()
	svr.Cmd.Process.Kill()
}
