package colorlogger

import (
	"fmt"
	"io"
	"log"
)

type colorLogger struct {
	color    bool
	logger   *log.Logger
	logLevel int
}

const (
	NONE  = -1
	DEBUG = 0
)

var logLevelsNames = []string{"debug"}

func ConstructLogger(logLevel int, color bool, writer io.Writer) *colorLogger {
	logger := log.New(writer, "", 0)
	loggerPlus := colorLogger{color, logger, logLevel}
	return &loggerPlus
}

func (g *colorLogger) Logf(logLevel int, msg string) {
	if g.logLevel >= logLevel && g.logLevel != NONE {
		if g.color {
			g.logger.Printf("\033[32m%s:\033[0m %s", logLevelsNames[logLevel], msg)
		} else {
			g.logger.Printf("%s: %s", logLevelsNames[logLevel], msg)

		}
	}
}

// This function is here so that existing Debugf outputs will still work.
// Once we figure out how to properly deal with logging, this can be revisited
func (g *colorLogger) Debugf(format string, a ...interface{}) {
	g.Logf(DEBUG, fmt.Sprintf(format, a...))
}
