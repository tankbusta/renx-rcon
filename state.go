package rcon

// GameStateManager is tied to a server and updates the state of the game
// as RCON messages flow over the wire.
//
// It will optionally issue commands to the server if it's connected to ensure
// our state is as accurate as possible.
type GameStateManager struct {
	Players Players

	Map string
}

func NewGameState() *GameStateManager {
	return &GameStateManager{
		Players: make(Players, 0), // While players are mostly limited to 64, server admins might go crazy
	}
}
