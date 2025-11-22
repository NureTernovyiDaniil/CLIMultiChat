package messengers

type Messenger interface {
	SendMessage(channel, message string) error
	GetName() string
}
