package state

import (
	"sort"
	"time"

	"github.com/tankbusta/renx-rcon/games"
)

// Player represents an individual in the game whether they're real or a bot
type Player struct {
	// SpawnedAt is the time the player was spawned into the game
	//
	// It's imperative that `.IsZero()` is used because this may be a zero value
	// if we've recently connected to the game
	SpawnedAt time.Time

	// LastDeath is the time the player last died
	//
	// It's imperative that `.IsZero()` is used because this may be a zero value
	// if we've recently connected to the game
	LastDeath time.Time

	// LastUpdated is the time this player record was last updated by a game event
	LastUpdated time.Time

	// Name of the player
	Name string

	// ID of the player
	ID int

	Score int

	// Team the player is on (GDI or NOD)
	Team games.Team

	IsBot       bool
	IsDeveloper bool
	IsAdmin     bool

	HardwareID string
	SteamID    string
}

type Players []*Player

func (s Players) Len() int { return len(s) }

func (s Players) Less(i, j int) bool { return s[i].ID < s[j].ID }
func (s Players) Swap(i, j int)      { s[i].ID, s[j].ID = s[j].ID, s[i].ID }

func (s Players) LocatePlayer(id int) *Player {
	playerIdx := sort.Search(id, func(i int) bool {
		return s[i].ID == id
	})

	if found := playerIdx != -1; found {
		return s[playerIdx]
	}

	return nil
}

func (s Players) DeleteByID(id int) bool {
	playerIdx := sort.Search(id, func(i int) bool {
		return s[i].ID == id
	})

	if found := playerIdx != -1; found {
		s = append(s[:playerIdx], s[playerIdx+1:]...)
		sort.Sort(s) // double check the list is sorted again
		return true
	}

	return false
}
