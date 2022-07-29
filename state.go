package rcon

import (
	"context"
	"log"
	"time"

	"github.com/tankbusta/renx-rcon/state"
)

const (
	playerStateMsg = "cClientVarList NAME ID IP HWID PING TEAM STEAM ADMIN SCORE CREDITS CHARACTER\n"
	botStateMsg    = "cBotVarList ID NAME TEAM SCORE CREDITS CHARACTER\n"
)

type IServer interface {
	WriteMsg(msg string)
	Ready() bool
}

// GameStateManager is tied to a server and updates the state of the game
// as RCON messages flow over the wire.
//
// It will optionally issue commands to the server if it's connected to ensure
// our state is as accurate as possible.
type GameStateManager struct {
	Players state.Players

	Map string

	LastUpdated time.Time

	// unexported fields below
	parent IServer
}

func NewGameState(parent IServer) *GameStateManager {
	gsm := &GameStateManager{
		Players: make(state.Players, 0), // While players are mostly limited to 64, server admins might go crazy
		parent:  parent,
	}

	return gsm
}

// dispatchStateCheck sends several messages to the server to verify the game state matches
func (s *GameStateManager) dispatchStateCheck() error {
	// s.parent.WriteMsg(playerStateMsg)
	s.parent.WriteMsg(botStateMsg)

	s.LastUpdated = time.Now()
	return nil
}

func (s *GameStateManager) Start(ctx context.Context) {
	ticker := time.NewTicker(time.Second * 5)
	defer ticker.Stop()

StateLoop:
	for {
		select {
		case <-ctx.Done():
			break StateLoop
		case <-ticker.C:
			log.Println("[ !! ] Dispatching state check")
			if s.parent.Ready() {
				s.dispatchStateCheck()
			}
		}
	}
}
