package container

import (
	"MCDaemon-go/config"
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

func (c *Container) Add(name string, workDir string, svr lib.Server) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.Servers[name] = svr
	//读取配置
	argv := config.GetStartConfig()
	c.Group.Add(1)
	go svr.Start(name, argv, workDir)
}

func (c *Container) Del(name string) {
	c.lock.Lock()
	defer c.lock.Unlock()
	if _, ok := c.Servers[name]; ok {
		deleteServer := c.Servers[name]
		deleteServer.Close()
		delete(c.Servers, name)
	}
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

func (c *Container) IsRuntime(name string) bool {
	c.lock.Lock()
	defer c.lock.Unlock()
	if _, ok := c.Servers[name]; ok {
		return true
	}
	return false
}
