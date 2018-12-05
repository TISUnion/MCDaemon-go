package lib

import "sync"

type CommandList struct {
	Commands []*Command
	Lock     sync.Mutex
}

//单例
var (
	cl       *CommandList
	is_build bool
)

const listSize = 20

//添加一个待运行命令至队尾
func (cl *CommandList) Push(c *Command) bool {
	cl.Lock.Lock()
	defer cl.Lock.Unlock()
	clLen := len(cl.Commands)
	if clLen < listSize {
		cl.Commands[clLen] = c
		return false
	}
	return true
}

//出队一个命令，准备执行
func (cl *CommandList) Pop() *Command {
	cl.Lock.Lock()
	clLen := len(cl.Commands)
	c := cl.Commands[clLen-1]
	cl.Commands = cl.Commands[:clLen-2]
	cl.Lock.Unlock()
	return c
}

func init() {
	is_build = false
}

func GetInstance() *CommandList {
	if !is_build {
		cl = &CommandList{}
		cl.Commands = make([]*Command, 0, listSize)
	}
	return cl
}
