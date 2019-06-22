package ChatPlugin

type WebSocketRS interface {
	Start() error
	Send(*Message)
	Read()
	GetId() int
	GetName() string
	IsAlive() bool
}
