package commands

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/tankbusta/renx-rcon/events"
	"github.com/tankbusta/renx-rcon/state"
)

type ICommand interface {
	Command() string
	MarshalRCON() string
	UnmarshalRCON(msg string, v any) error
}

var (
	botFields    = []string{"ID", "NAME", "TEAM", "SCORE", "CREDITS", "CHARACTER"}
	numBotFields = len(botFields)
)

type ListBotsCommand struct{}

func (s ListBotsCommand) Command() string {
	return "BotVarList"
}

func (s ListBotsCommand) MarshalRCON() string {
	return s.Command() + " " + strings.Join(botFields, events.Delimiter) + "\n"
}

func (s ListBotsCommand) UnmarshalRCON(msg string, v any) error {
	player, ok := v.(*state.Player)
	if !ok {
		return fmt.Errorf(
			"cannot UnmarshalRCON ListBotsCommand into %T. Expected *state.Player",
			v,
		)
	}

	parts := strings.Split(msg, events.Delimiter)
	if len(parts) != numBotFields {
		return fmt.Errorf(
			"unexpected number of fields in ListBotsCommand. expected %d got %d",
			numBotFields, len(parts),
		)
	}
	_ = parts[numBotFields] // Bounds check elimination

	playerID, err := strconv.Atoi(parts[0])
	if err != nil {
		return fmt.Errorf("failed to parse playerID %s: %w", parts[0], err)
	}

	player.ID = playerID
	player.Name = parts[1]
	player.Team.ParseString(parts[2])

	score, err := strconv.Atoi(parts[3])
	if err != nil {
		return fmt.Errorf("failed to parse playerID %s: %w", parts[3], err)
	}

	player.Score = score
	player.LastUpdated = time.Now()

	return nil
}
