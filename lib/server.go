package lib

type Server interface {
	Say(string)
	Tell(string, string)
	Execute(string)
}
