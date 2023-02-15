package events

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Version struct {
	RCONVersion       int
	GameVersion       int
	GameVersionPretty string
}

func (s Version) String() string {
	return fmt.Sprintf("RCON v%d on Game Version %d :: %s", s.RCONVersion, s.GameVersion, s.GameVersionPretty)
}

func (s *Version) Parse(input string) error {
	parts := strings.SplitN(input, string(Delimiter), 3)
	if len(parts) != 3 {
		return fmt.Errorf("Unknown RCON version format: %s", input)
	}
	_ = parts[2]

	rv, err := strconv.Atoi(parts[0])
	if err != nil {
		return fmt.Errorf("Unknown RCON version `%s`: %w", parts[0], err)
	}
	s.RCONVersion = rv

	gv, err := strconv.Atoi(parts[1])
	if err != nil {
		return fmt.Errorf("Unknown Game version `%s`: %w", parts[1], err)
	}

	s.GameVersion = gv
	s.GameVersionPretty = strings.Trim(parts[2], "\n")

	return nil
}

type ServerError struct {
	ErrorMsg string
}

func (s ServerError) String() string {
	return s.ErrorMsg
}

func (s ServerError) Error() string {
	return s.String()
}

func (s *ServerError) Parse(input string) error {
	s.ErrorMsg = input
	return nil
}

type (
	ServerInfo struct {
		Port             int    `rcon:"PORT"`
		Name             string `rcon:"SERVERNAME"`
		Map              string `rcon:"GETPACKAGENAME"`
		NumPlayers       int    `rcon:"NUMPLAYERS"`
		MaxPlayers       int    `rcon:"MAXPLAYERS"`
		VehicleLimit     int    `rcon:"VEHICLELIMIT"`
		MineLimit        int    `rcon:"MINELIMIT"`
		TimeLimit        int    `rcon:"TIMELIMIT"`
		RequiresPassword bool   `rcon:"REQUIRESPASSWORD"`
		RequiresSteam    bool   `rcon:"BSTEAMREQUIRED"`
		IsListed         bool   `rcon:"BLISTED"`
		GameVersion      string `rcon:"GAMEVERSION"`
		GameMode         string `rcon:"MODE"`
		InProgress       bool   `rcon:"MATCHINPROGRESS"`
		GameState        string `rcon:"GAMESTATE"`

		LastUpdated time.Time
	}

	PlayerBase struct {
		ID     int    `rcon:"ID"`
		Team   string `rcon:"TEAM"`
		Name   string `rcon:"NAME"`
		Score  int    `rcon:"SCORE"`
		Deaths int    `rcon:"DEATHS"`
		IsSpy  bool   `rcon:"SPY"`

		LastUpdated time.Time
	}

	Bot struct {
		PlayerBase
	}
)
