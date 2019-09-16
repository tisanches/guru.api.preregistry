package logger

import (
	"fmt"
	Logger "github.com/bestmethod/go-logger"
	"os"
	"strings"
)

var LOG = Logger.Logger{}

func InitLog(header string, service string, level string){
	err := LOG.Init(header, service, parseLogLevel(level), parseLogLevel(level), 0)
	if err != nil {
		fmt.Fprintf(os.Stderr, "CRITICAL Could not initialize logger. Quitting. Details: %s\n", err)
		os.Exit(1)
	}
}

func parseLogLevel(level string) int{
	switch strings.ToLower(level) {
	case "critical":
		return Logger.LEVEL_CRITICAL
	case "error":
		return Logger.LEVEL_ERROR
	case "warn":
		return Logger.LEVEL_WARN
	case "info":
		return Logger.LEVEL_INFO
	default:
		return Logger.LEVEL_DEBUG
	}
}
