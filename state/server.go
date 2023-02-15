package state

import (
	"github.com/tankbusta/renx-rcon/events"
	"github.com/tankbusta/renx-rcon/games"
)

type (
	Server struct {
		// Game indicates which Totem Arts game this server is running
		Game games.Game

		// Address of the Game Server
		Address string

		// ConnectionID of the RCON connection
		ConnectionID string

		// Version of the UDK Game Server
		Version events.Version

		// IsAuthenticated marks if the server has accepted our RCON credentials
		IsAuthenticated bool

		// IsConnected marks if the server has been connected to (but not necessairly authenticated)
		IsConnected bool

		// HasStreamEnabled indicates if we've enabled the event stream from the game server
		HasStreamEnabled bool

		LatestInfo events.ServerInfo
	}
)
