package ChatPlugin

type WebSocketRS interface {
	Start()
	Send(Message)
	Read() Message
}
