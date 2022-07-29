package games

import "strings"

type Team uint8

const (
	TeamUnknown Team = iota
	TeamGDI
	TeamNOD
)

func (s Team) String() string {
	switch s {
	case TeamGDI:
		return "GDI"
	case TeamNOD:
		return "NOD"
	default:
		return "Unknown"
	}
}

func (s *Team) ParseString(input string) {
	switch strings.ToLower(input) {
	case "gdi":
		*s = TeamGDI
	case "nod":
		*s = TeamNOD
	}
}

type Game uint8

const (
	GameUnknown Game = iota
	GameRenegadeX
	GameFirestormAssault
	GameFirestorm
)

func (s Game) String() string {
	switch s {
	case GameRenegadeX:
		return "Renegade X"
	case GameFirestormAssault:
		return "Firestorm Assault"
	case GameFirestorm:
		return "Firestorm"
	default:
		return "Unknown"
	}
}

func (s *Game) ParseString(input string) {
	switch strings.ToLower(input) {
	case "renegade x", "renx", "renegadex", "renegade-x":
		*s = GameRenegadeX
	case "firestorm assault", "fa":
		*s = GameFirestormAssault
	case "firestorm", "fs":
		*s = GameFirestorm
	default:
		*s = GameUnknown
	}
}
