package util

import (
	"fmt"
	"log"
	"os"
)

func GetLogConfig(level string) *log.Logger {
	var standerizedLevel = fmt.Sprintf("[%s] ", level)
	var logger *log.Logger = log.New(os.Stdout, standerizedLevel, log.LstdFlags)
	logger.SetFlags(log.LstdFlags | log.Lshortfile)
	return logger
}
