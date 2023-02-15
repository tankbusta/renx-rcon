package commands

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/tankbusta/renx-rcon/events"
	"github.com/tankbusta/renx-rcon/state"
)

type HandleCommandResp func(cmd ICommand, resp string)

type ICommand interface {
	Command() string
	MarshalRCON() []byte
	UnmarshalRCON(msg string, v any) error
	SkipFirstMsg() bool
}

var (
	botFields    = []string{"ID", "NAME", "TEAM", "SCORE", "CREDITS", "CHARACTER"}
	numBotFields = len(botFields)
)

type ListBotsCommand struct{}

func NewListBotsCommand() ListBotsCommand {
	return ListBotsCommand{}
}

func (s ListBotsCommand) SkipFirstMsg() bool {
	return true
}

func (s ListBotsCommand) Command() string {
	return "BotVarList"
}

func (s ListBotsCommand) MarshalRCON() []byte {
	return []byte("c" + s.Command() + " " + strings.Join(botFields, " ") + "\n")
}

func (s ListBotsCommand) UnmarshalRCON(msg string, v any) error {
	player, ok := v.(*state.Player)
	if !ok {
		return fmt.Errorf(
			"cannot UnmarshalRCON ListBotsCommand into %T. Expected *state.Player",
			v,
		)
	}

	parts := strings.Split(msg, string(events.Delimiter))
	if len(parts) != numBotFields {
		return fmt.Errorf(
			"unexpected number of fields in ListBotsCommand. expected %d got %d",
			numBotFields, len(parts),
		)
	}

	_ = parts[numBotFields-1] // Bounds check elimination

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
