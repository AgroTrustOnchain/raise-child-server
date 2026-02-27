package util

import (
	"encoding/hex"
	"strings"
)

func IsValidSuiAddressStrict(input string) bool {
	// Kiểm tra prefix và độ dài tổng (0x + 64 hex = 66)
	if len(input) != 66 || !strings.HasPrefix(input, "0x") {
		return false
	}

	// Lấy phần thân (64 ký tự hex)
	addrBody := input[2:]

	// Kiểm tra xem có phải Hex hợp lệ không
	if _, err := hex.DecodeString(addrBody); err != nil {
		return false
	}

	return true
}
