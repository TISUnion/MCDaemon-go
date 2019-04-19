package ChatPlugin

type WebSocketRS interface {
	Start() error
	Send(*Message)
	Read(chan *msgPackage)
	GetId() int
	GetName() string
}
