package commands

import (
	"strings"

	"github.com/tankbusta/renx-rcon/events"
)

var serverInfoFields = events.GetAllColumns(&events.ServerInfo{})
var serverInfoHeader = serverInfoFields.GenerateDefaultHeader()

type ServerInfoCommand struct{}

func NewServerInfoCommand() ServerInfoCommand {
	return ServerInfoCommand{}
}

func (s ServerInfoCommand) SkipFirstMsg() bool {
	return true
}

func (s ServerInfoCommand) Command() string {
	return "ServerInfo"
}

func (s ServerInfoCommand) MarshalRCON() []byte {
	return []byte("c" + s.Command() + " " + strings.Join(serverInfoFields, " ") + "\n")
}

func (s ServerInfoCommand) UnmarshalRCON(msg string, v any) error {
	return events.UnmarshalRCON(serverInfoHeader, msg, v)
}
