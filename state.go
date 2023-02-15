package rcon

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/tankbusta/renx-rcon/commands"
	"github.com/tankbusta/renx-rcon/state"
)

var cmdUpdateBotState = commands.NewListBotsCommand()

type IServer interface {
	WriteMsg(msg commands.ICommand, cb commands.HandleCommandResp)
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

func (s *GameStateManager) onBotStateUpdate(cmd commands.ICommand, data string) {
	var p state.Player

	if err := cmd.UnmarshalRCON(data, &p); err != nil {
		log.Printf("[ !! ] Failed to unmarshal bot state: %s", err)
		return
	}

	log.Println("[ !! ] Received bot state update")
	fmt.Println(p)
}

// dispatchStateCheck sends several messages to the server to verify the game state matches
func (s *GameStateManager) dispatchStateCheck() error {
	s.parent.WriteMsg(cmdUpdateBotState, s.onBotStateUpdate)

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

			// We need to send a message to RCON at least once every 60 seconds otherwise the game server will disconnect us
			// So let's take this opportunity to update our state!
			if s.parent.Ready() {
				s.dispatchStateCheck()
			}
		}
	}
}
