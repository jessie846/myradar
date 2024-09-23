package utils

import (
	"strings"
)

// FacilityChar returns the character associated with the facility code.
func FacilityChar(facility string) rune {
	switch facility {
	case "ZNY":
		return 'N'
	case "ZBW":
		return 'B'
	case "ZDC":
		return 'W'
	case "ZOB":
		return 'C'
	case "ACY":
		return 'Y'
	case "BBB":
		return 'G'
	case "JST":
		return 'F'
	case "NVF":
		return 'K'
	case "ROC":
		return 'T'
	case "RRR":
		return 'D'
	case "WWW":
		return 'V'
	case "ZJX":
		return 'J'
	case "ZTL":
		return 'T'
	case "ZMA":
		return 'Z'
	case "ZHU":
		return 'H'
	case "CLT":
		return 'E'
	case "TPA":
		return 'D'
	default:
		if len(facility) > 0 {
			return rune(facility[0])
		}
		return 'X'
	}
}

// isWhitespace checks if a character is whitespace.
func isWhitespace(c rune) bool {
	return c == ' ' || c == '\t'
}

// startNewToken determines if a new token should be started based on the current character.
func startNewToken(token string, c rune) bool {
	if isWhitespace(c) {
		return true
	}
	if len(token) > 0 {
		lastChar := rune(token[len(token)-1])
		return lastChar == '.' && c != '.'
	}
	return false
}

// WrapQFOutput wraps text to the specified width in characters.
func WrapQFOutput(text string, widthInChars int) string {
	lines := []string{}
	currentLine := ""
	currentToken := ""

	for _, c := range text {
		if c == '\n' {
			currentLine += currentToken
			lines = append(lines, currentLine)
			currentToken = ""
			currentLine = ""
		} else {
			shouldStartNewToken := startNewToken(currentToken, c)

			if !shouldStartNewToken {
				currentToken += string(c)
			}

			currentLineLength := len([]rune(currentLine))
			currentTokenLength := len([]rune(currentToken))

			if currentLineLength+currentTokenLength > widthInChars {
				lines = append(lines, currentLine)
				if shouldStartNewToken {
					currentLine = currentToken
					currentToken = string(c)
				} else {
					currentLine = ""
				}
			} else if currentLineLength+currentTokenLength == widthInChars && shouldStartNewToken {
				currentLine += currentToken
				lines = append(lines, currentLine)
				currentToken = string(c)
				currentLine = ""
			} else if shouldStartNewToken {
				currentLine += currentToken
				currentToken = string(c)
			}
		}
	}

	if len(currentToken) > 0 {
		currentLine += currentToken
	}
	if len(currentLine) > 0 {
		lines = append(lines, currentLine)
	}

	return strings.Join(lines, "\n")
}
