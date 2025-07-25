package initiator

import (
	"os"
	"sync"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func get_default_logger() (zerolog.Logger, error) {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnixNano
	logger := log.Output(zerolog.ConsoleWriter{Out: os.Stderr}).
		With().
		Str("tag", "bhgl").
		Logger()

	level_str := os.Getenv("LOG_LEVEL")
	if len(level_str) == 0 {
		level_str = "info"
	}
	level, err := zerolog.ParseLevel(level_str)
	if err != nil {
		return logger, err
	}
	zerolog.SetGlobalLevel(level)
	return logger, nil
}

var GetDefaultLogger = sync.OnceValues(get_default_logger)
var InitLogger = sync.OnceValue(func() (err error) {
	log.Logger, err = GetDefaultLogger()
	return
})
