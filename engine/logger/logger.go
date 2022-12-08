package logger

import (
	"flag"
	"github.com/hashicorp/go-hclog"
	"os"
)

var (
	engineLogLevel = flag.String("log.level.engine", os.Getenv("LOG_LEVEL_ENGINE"), "log level")
	httpLogLevel   = flag.String("log.level.http", os.Getenv("LOG_LEVEL_HTTP"), "log level")
	Engine         hclog.Logger
	HttpLogger     hclog.Logger
)

func init() {

	Engine = hclog.New(&hclog.LoggerOptions{
		Name:                 "engine",
		Level:                hclog.LevelFromString(*engineLogLevel),
		JSONFormat:           false,
		TimeFormat:           "2006-01-02 15:04:05",
		Color:                100,
		ColorHeaderOnly:      true,
		ColorHeaderAndFields: true,
		IndependentLevels:    true,
	})

	Engine.Debug("Engine logger initialised")
	HttpLogger = Engine.Named("http")
	HttpLogger.SetLevel(hclog.LevelFromString(*httpLogLevel))
	HttpLogger.Debug("Http logger initialised")
}
