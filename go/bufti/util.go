package bufti

import (
	"strings"
)

func isListType(valType BuftiType) (BuftiType, bool) {
	s, found := strings.CutPrefix(string(valType), "list:")
	return BuftiType(s), found
}

func isMapType(valType BuftiType) (BuftiType, BuftiType, bool) {
	s, found := strings.CutPrefix(string(valType), "map:")
	if !found {
		return "", "", false
	}
	parts := strings.Split(s, ":")
	if len(parts) != 2 {
		return "", "", false
	}
	return BuftiType(parts[0]), BuftiType(parts[1]), true
}

func isModelType(valType BuftiType) (string, bool) {
	return strings.CutPrefix(string(valType), "model:")
}

func isInRange(v float64, min float64, max float64) bool {
	return v >= min && v <= max
}
