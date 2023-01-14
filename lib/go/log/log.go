package log

import (
	"flag"
	"os"

	log "github.com/sirupsen/logrus"
)

var logLevelFlag = flag.String("log-level", log.InfoLevel.String(), "Level of logging")

func Init() *log.Logger {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)

	level, err := log.ParseLevel(*logLevelFlag)
	if err != nil {
		log.Panicf("Unknown debug level \"%s\"", *logLevelFlag)
	}
	log.SetLevel(level)

	return log.StandardLogger()
}

func Logger() *log.Logger {
	return log.StandardLogger()
}
