package util

import (
	"strconv"
	"time"
)

func GenerateNumber() int {
	millis := time.Now().UnixNano() / int64(time.Millisecond)
	millisStr := strconv.FormatInt(millis, 10)
	number, _ := strconv.Atoi(millisStr[len(millisStr)-6:])
	return number
}
