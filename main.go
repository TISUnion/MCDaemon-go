package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
)

func main() {
	cmd := exec.Command("java", "-Xmx1024M", "-Xms1024M", "-jar", "minecraft/server.jar", "nogui")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}
	defer stdout.Close()
	res, _ := ioutil.ReadAll(stdout)
	fmt.Println(string(res))
}
