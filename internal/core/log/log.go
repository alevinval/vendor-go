package log

import (
	"os"

	log "github.com/withmandala/go-log"
)

var logger = log.New(os.Stderr)

func init() {
	logger.WithoutTimestamp()
}

func GetLogger() *log.Logger {
	return logger
}

func EnableDebug() {
	logger.WithDebug()
}
