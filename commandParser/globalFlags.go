package commandParser

import (
	"log"
	"os"
)

type GlobalFlags struct {
	Debug bool
	Color bool
}

func (g *GlobalFlags) GetDebug() func(format string, a ...interface{}) {

	debugFunc := func(format string, a ...interface{}) {}
	prefix := "debug: "
	if g.Color {
		prefix = "\033[32m" + prefix + "\033[0m"
	}
	if g.Debug {
		debugFunc = log.New(os.Stderr, prefix, 0).Printf
	}
	return debugFunc
}
