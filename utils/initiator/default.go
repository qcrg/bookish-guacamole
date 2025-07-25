package initiator

import (
	"sync"

	"github.com/joho/godotenv"
)

func init_all() (err error) {
	godotenv.Load()
	err = InitLogger()
	return err
}

var Preinit = sync.OnceValue(init_all)
