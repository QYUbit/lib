package bufti

import (
	"strings"
)

func isListType(valType BuftiType) (BuftiType, bool) {
	s, found := strings.CutPrefix(string(valType), "[")
	if !found {
		return "", false
	}
	after, found := strings.CutSuffix(s, "]")
	return BuftiType(after), found
}

func isInRange(v float64, min float64, max float64) bool {
	return v >= min && v <= max
}
