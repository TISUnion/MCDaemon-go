/**
 * 热插件栗子
 * 这是一个实现复读姬的栗子
 */
package main

import (
	"fmt"
	"os"
)

func main() {
	//获取指令参数
	//注意 args 为数组，第0个参数是命令本身，是不需要的
	args := os.Args
	args = args[1:]
	//返回调用MCD的信息
	//注意，不能换行;要做没有传入参数的判断
	if len(args) != 0 {
		fmt.Print("say ", args[0])
	}
}
