package container

import (
	"MCDaemon-go/lib"
	"sync"
)

//单例模式
var (
	ContainerInstance *Container
	once              sync.Once
)

//获取单例实例
func GetInstance() *Container {
	once.Do(func() {
		ContainerInstance = &Container{}
		ContainerInstance.Init()
	})
	return ContainerInstance
}

type Container struct {
	Servers map[string]lib.Server //所有已启动的服务器
	lock    sync.Mutex            //同步锁
	Group   sync.WaitGroup        //协程组同步
}

func (c *Container) Init() {
	c.Servers = make(map[string]lib.Server)
}

func (c *Container) Add(name string, ls lib.Server) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.Servers[name] = ls
}

func (c *Container) Del(name string) {
	c.lock.Lock()
	defer c.lock.Unlock()
	deleteServer := c.Servers[name]
	deleteServer.RunUniquePlugin(deleteServer.Close)
	delete(c.Servers, name)
}

func (c *Container) GetRuntimeServer() []string {
	c.lock.Lock()
	defer c.lock.Unlock()
	var res []string
	for k, _ := range c.Servers {
		res = append(res, k)
	}
	return res
}
