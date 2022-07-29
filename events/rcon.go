package events

type (
	ClientType byte
	ServerType byte
)

type SupportedMsgTypes interface {
	ClientType | ServerType
}

type Header[T SupportedMsgTypes] struct {
	Type      T
	EventName string
}

type EventParser interface {
	Parse(string) error
	String() string
}

const (
	NewLine   = 0x0A // '\n'
	Delimiter = string(0x02)
)

const (
	// Client->Server
	Authenticate ClientType = 'a'
	Subscribe    ClientType = 's'
	UnSubscribe  ClientType = 'u'
	Command      ClientType = 'c'
	DevBot       ClientType = 'd'

	// Server -> Client
	RCONGameVersion          ServerType = 'v'
	AuthenticationSuccess    ServerType = 'a'
	Error                    ServerType = 'e'
	CommandResponse          ServerType = 'r'
	CommandExecutionFinished ServerType = 'c'
	GameLog                  ServerType = 'l'
	ServerDevBot             ServerType = 'd'
)
