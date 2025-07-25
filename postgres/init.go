package postgres

import (
	"sync"

	_ "github.com/lib/pq"
	"github.com/qcrg/bookish-guacamole/utils/initiator"
	"github.com/rs/zerolog"
)

var log zerolog.Logger

func init_func() error {
	var err error
	initiator.Preinit()
	dlog, err := initiator.GetDefaultLogger()
	if err != nil {
		return err
	}
	log = dlog.With().
		Str("tag", "database").
		Str("type", "postgres").
		Logger()
	return nil
}

var Init = sync.OnceValue(init_func)
