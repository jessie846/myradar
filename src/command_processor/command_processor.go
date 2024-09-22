package command_processor

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// CommandError holds error information
type CommandError struct {
	Message string
}

// Implement the error interface
func (e CommandError) Error() string {
	return e.Message
}

// Command represents different commands
type Command interface{}

// Define different types of commands
type ChangeDatablockLeaderLength struct {
	Length int
	Flid   string
}

type ChangeDatablockPosition struct {
	Position rune
	Flid     string
}

type ChangeSector struct {
	SectorID string
}

type ShowFlightPlan struct {
	Flid string
}

type ToggleFDB struct {
	Flid string
}

// Helper function to check if the string is a single digit
func isDigit(s string) bool {
	return len(s) == 1 && s[0] >= '0' && s[0] <= '9'
}

// Helper function to check if the string matches leader line length "/N"
func isLeaderLineLength(s string) bool {
	return len(s) == 2 && s[0] == '/' && s[1] >= '0' && s[1] <= '9'
}

// Helper function to check if a string could be a flight ID
func isMaybeFlid(val string) bool {
	return len(val) >= 2 && len(val) <= 10
}

// Helper function to check if the string is a valid CID

// ParseCommand parses the input string and returns the corresponding command or error
func ParseCommand(input string) (Command, error) {
	pieces := strings.Fields(input)
	if len(pieces) == 0 {
		return nil, errors.New("empty input")
	}

	if pieces[0] == "QF" && len(pieces) == 2 {
		return ShowFlightPlan{Flid: pieces[1]}, nil
	} else if pieces[0] == "SI" && len(pieces) == 2 {
		return ChangeSector{SectorID: pieces[1]}, nil
	} else if isDigit(pieces[0]) {
		return ChangeDatablockPosition{
			Position: rune(pieces[0][0]),
			Flid:     pieces[1],
		}, nil
	} else if isLeaderLineLength(pieces[0]) {
		length, _ := strconv.Atoi(string(pieces[0][1]))
		return ChangeDatablockLeaderLength{
			Length: length,
			Flid:   pieces[1],
		}, nil
	} else if len(pieces) == 1 && isMaybeFlid(pieces[0]) {
		return ToggleFDB{Flid: pieces[0]}, nil
	}

	return nil, CommandError{Message: fmt.Sprintf("%s FORMAT", pieces[0])}
}
