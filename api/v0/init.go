package v0

import (
	"errors"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/versioning"
	"github.com/qcrg/bookish-guacamole/api/v0/token"
)

type handler = func(iris.Context, *Deps)

func make_handler(fn handler, deps *Deps) iris.Handler {
	return func(ctx iris.Context) {
		fn(ctx, deps)
	}
}

func InitV0(app *versioning.Group, deps *Deps) error {
	err := token.Init()
	if err != nil {
		return errors.Join(errors.New("Failed to initiate token"), err)
	}

	app.Post("/auth", make_handler(create_session, deps))
	app.Put("/sessions/current", make_handler(refresh_current_session, deps))

	gparty := app.Party("/")
	gparty.Use(make_handler(auth, deps))

	gparty.Get("/me", make_handler(me, deps))
	gparty.Delete("/sessions/current", make_handler(close_current_session, deps))
	return nil
}
