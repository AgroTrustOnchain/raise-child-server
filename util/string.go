package util

import "strings"

func StanderizeString(src string) string {
	return strings.TrimSpace(strings.ToLower(src))
}

func FormatAddress(address string) string {
	return address[:8] + "..." + address[len(address)-6:]
}
