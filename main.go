package main

import (
	"MCDaemon-go/config"
	"bufio"
	"fmt"
	"io"
	"log"
	"os/exec"
)

var (
	MCDconfig   map[string]string
	commandArgu []string
)

func getres(bufReader *bufio.Reader) {
	var buffer []byte = make([]byte, 4096)
	for {
		n, err := bufReader.Read(buffer)
		if err != nil {
			if err == io.EOF {
				fmt.Printf("pipi has Closed\n")
			} else {
				fmt.Println("Read content failed")
			}
			break
		}
		fmt.Printf("%s\n\n", string(buffer[:n]))
	}
}

func init() {
	//判断eula是否为真
	config.SetEula()

	//加载服务器启动配置
	MCDconfig = config.GetConfig()
	commandArgu = []string{
		MCDconfig["Xmx"],
		MCDconfig["Xms"],
		"-jar",
		fmt.Sprintf("%s/%s", MCDconfig["serverPath"], MCDconfig["serverName"]),
	}
	if MCDconfig["gui"] != "true" {
		commandArgu = append(commandArgu, "nogui")
	}
}

func main() {

	cmd := exec.Command("java", commandArgu...)
	stdout, err := cmd.StdoutPipe()
	bufReader := bufio.NewReader(stdout)
	if err != nil {
		log.Fatal(err)
	}
	go getres(bufReader)
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}
	defer stdout.Close()
	defer cmd.Process.Kill()
	cmd.Wait()
}
