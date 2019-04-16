package ChatPlugin

type WebSocketRS interface {
	Start() error
	Send(*Message)
	Read() *Message
}
