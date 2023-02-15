package events

import (
	"fmt"
	"strings"
)

type LogType string

const (
	// LogTypeGame is a log containing changes to the game state
	LogTypeGame LogType = "GAME"

	// LogTypeRCON is a log containing the command executed by RCON/DevBot
	LogTypeRCON LogType = "RCON"

	LogTypePlayer LogType = "PLAYER"

	LogTypeLog LogType = "LOG"

	LogTypeChat LogType = "CHAT"

	LogTypeMap LogType = "MAP"
)

type LogActivity string

const (
	ActivityCommand   LogActivity = "Command"
	ActivitySpawn     LogActivity = "Spawn"
	ActivityDeath     LogActivity = "Death"
	ActivityPurchase  LogActivity = "Purchase"
	ActivityCrate     LogActivity = "Crate"
	ActivityDestroyed LogActivity = "Destroyed"
	ActivitySay       LogActivity = "Say"
	// MapStart
	ActivityStart LogActivity = "Start"

	// Player Activities
	ActivityHWID     LogActivity = "HWID"
	ActivityEnter    LogActivity = "Enter"
	ActivityExit     LogActivity = "Exit"
	ActivityTeamJoin LogActivity = "TeamJoin"
	ActivityChangeID LogActivity = "ChangeID"
)

type LogMessage struct {
	Type     LogType
	Activity LogActivity
	Parts    []string

	// Message contains a human readable representation of this event
	// Currently only set when Type == LogTypeLog
	Message string
}

func ParseLog(msg string) (LogMessage, error) {
	var lm LogMessage

	headerIdx := strings.Index(msg, ";")
	if headerIdx == -1 {
		// Sometimes there wont be one of these, it's usually a log
		lm.Type = LogTypeLog
		lm.Message = strings.Trim(msg, "\n")

		return lm, nil
	}

	parts := strings.SplitN(msg[:headerIdx], string(Delimiter), 2)
	if len(parts) != 2 {
		return lm, fmt.Errorf("rcon/log: missing delimiter in header: %s", msg[:headerIdx])
	}

	switch LogType(parts[0]) {
	case LogTypeGame:
		lm.Type = LogTypeGame
	case LogTypeRCON:
		lm.Type = LogTypeRCON
	case LogTypeChat:
		lm.Type = LogTypeChat
	case LogTypeMap:
		lm.Type = LogTypeMap
	case LogTypePlayer:
		lm.Type = LogTypePlayer
	default:
		return lm, fmt.Errorf("rcon/log: unknown log type : %s", parts[0])
	}

	switch LogActivity(parts[1]) {
	case ActivityCommand:
		lm.Activity = ActivityCommand
	case ActivitySpawn:
		lm.Activity = ActivitySpawn
	case ActivityDeath:
		lm.Activity = ActivityDeath
	case ActivityPurchase:
		lm.Activity = ActivityPurchase
	case ActivityCrate:
		lm.Activity = ActivityCrate
	case ActivityDestroyed:
		lm.Activity = ActivityDestroyed
	case ActivitySay:
		lm.Activity = ActivitySay
	case ActivityStart:
		lm.Activity = ActivityStart
	case ActivityHWID:
		lm.Activity = ActivityHWID
	case ActivityEnter:
		lm.Activity = ActivityEnter
	case ActivityTeamJoin:
		lm.Activity = ActivityTeamJoin
	case ActivityChangeID:
		lm.Activity = ActivityChangeID
	case ActivityExit:
		lm.Activity = ActivityExit
	default:
		return lm, fmt.Errorf("rcon/log: unknown log activity : %s", parts[1])
	}

	lm.Parts = strings.Split(strings.Trim(msg[headerIdx+2:], "\n"), string(Delimiter))

	return lm, nil
}
