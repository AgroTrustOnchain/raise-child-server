package util

import "strings"

func StanderizeString(src string) string {
	return strings.TrimSpace(strings.ToLower(src))
}
