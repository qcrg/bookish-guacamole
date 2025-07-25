package v0

import (
	"github.com/qcrg/bookish-guacamole/postgres"
	"github.com/rs/zerolog"
)

type Deps struct {
	Log zerolog.Logger
	DB  *postgres.DB
}
