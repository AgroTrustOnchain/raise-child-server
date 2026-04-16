package util

import (
	"fmt"
	"log"
	"os"
)

func GetLogConfig(level string) *log.Logger {
	var standardLevel = fmt.Sprintf("[%s] ", level)
	var logger *log.Logger = log.New(os.Stdout, standardLevel, log.LstdFlags)
	logger.SetFlags(log.LstdFlags | log.Lshortfile)
	return logger
}
