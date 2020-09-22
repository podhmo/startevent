package startevent

import (
	"log"
	"os"
)

type logger interface {
	Printf(string, ...interface{})
	Fatalf(string, ...interface{})
}

var (
	l *log.Logger
)

func getLogger() logger {
	if l != nil {
		return l
	}
	return log.New(os.Stderr, "startevent", log.Ldate)
}
