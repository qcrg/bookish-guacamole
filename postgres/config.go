package postgres

import (
	"fmt"
	"os"

	"github.com/qcrg/bookish-guacamole/utils"
)

type Config interface {
	GetConnectionString() string

	GetTLSMod() string
}

type ConfigEnv struct{}

const (
	_KEY_PREFIX       = "BHGL_DATABASE_"
	TLS_MOD_KEY       = _KEY_PREFIX + "TLS_MOD"
	HOST_KEY          = _KEY_PREFIX + "HOST"
	PORT_KEY          = _KEY_PREFIX + "PORT"
	USERNAME_KEY      = _KEY_PREFIX + "USERNAME"
	PASSWD_KEY        = _KEY_PREFIX + "PASSWD"
	DATABASE_NAME_KEY = _KEY_PREFIX + "NAME"
)

func (t ConfigEnv) GetConnectionString() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		t.GetLogin(),
		t.GetPasswd(),
		t.GetHost(),
		t.GetPort(),
		t.GetName(),
		t.GetTLSMod(),
	)
}

func (ConfigEnv) GetLogin() string {
	return os.Getenv(USERNAME_KEY)
}

func (ConfigEnv) GetPasswd() string {
	return os.Getenv(PASSWD_KEY)
}

func (ConfigEnv) GetHost() string {
	return os.Getenv(HOST_KEY)
}

func (ConfigEnv) GetPort() string {
	return utils.GetEnv(PORT_KEY, "5432")
}

func (ConfigEnv) GetName() string {
	return os.Getenv(DATABASE_NAME_KEY)
}

func (ConfigEnv) GetTLSMod() string {
	return utils.GetEnv(TLS_MOD_KEY, "disable")
}
