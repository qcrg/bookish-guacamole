package main

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/kataras/iris/v12"
	ictx "github.com/kataras/iris/v12/context"
	"github.com/kataras/iris/v12/core/host"
	"github.com/kataras/iris/v12/versioning"
	v0 "github.com/qcrg/bookish-guacamole/api/v0"
	"github.com/qcrg/bookish-guacamole/postgres"
	"github.com/qcrg/bookish-guacamole/utils/initiator"
	"github.com/rs/zerolog/log"
)

func init_v0(app *iris.Application, db *postgres.DB) {
	v0group := versioning.NewGroup(app, ">=0.0.0 <="+v0.Version)

	err := v0.InitV0(v0group, &v0.Deps{Log: log.Logger, DB: db})
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize API:v0")
	}
}

func init_default_routs(app *iris.Application) {
	app.Get("/health", func(ctx iris.Context) {
		ctx.StatusCode(iris.StatusOK)
	})

	app.Get("/version", func(ctx iris.Context) {
		ctx.JSON(iris.Map{"version": v0.Version})
	})
}

func log_request(ctx iris.Context) {
	version := versioning.GetVersion(ctx)
	start := time.Now()
	ctx.Next()

	log.Info().
		Str("version", version).
		Str("method", ctx.Method()).
		Str("path", ctx.Path()).
		Int("status_code", ctx.GetStatusCode()).
		Int64("latency_us", time.Since(start).Microseconds()).
		Send()
}

func make_runner() (iris.Runner, error) {
	port, exists := os.LookupEnv("BHGL_PORT")
	if !exists {
		log.Warn().Msg("BHGL_PORT is empty or not defined, set default port 8643")
		port = "8643"
	}
	addr := fmt.Sprintf(":%s", port)
	cert_path, exists := os.LookupEnv("BHGL_CERT_PATH")
	if !exists {
		return nil, errors.New("BHGL_CERT_PATH is empty or not defined")
	}
	key_path, exists := os.LookupEnv("BHGL_SKEY_PATH")
	if !exists {
		return nil, errors.New("BHGL_SKEY_PATH is empty or not defined")
	}
	return iris.TLS(addr, cert_path, key_path), nil
}

func init_version_failure() {
	versioning.NotFoundHandler = func(ctx *ictx.Context) {
		ctx.StopWithJSON(
			iris.StatusBadRequest,
			iris.Map{
				"error": iris.Map{
					"code":   -1,
					"reason": "Version is not found or not defined",
				},
			},
		)
	}
}

func startup_log(task host.TaskHost) {
	log.Info().Str("url", task.HostURL()).Msg("bookish-guacamole is running")
}

func main() {
	err := initiator.Preinit()
	if err != nil {
		log.Fatal().Err(err).Msg("Prenitialization failed")
	}
	err = postgres.Init()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize postgres")
	}

	db, err := postgres.NewDatabase(postgres.ConfigEnv{})
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create postgres database")
	}
	defer db.Close()

	app := iris.New()
	app.ConfigureHost(func(h *iris.Supervisor) {
		h.RegisterOnServe(startup_log)
	})
	app.Use(log_request)

	app.UseRouter(versioning.Aliases(versioning.AliasMap{
		"v0": v0.Version,
	}))

	init_version_failure()
	init_default_routs(app)
	init_v0(app, db)

	runner, err := make_runner()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create runner for application")
	}

	err = app.Run(
		runner,
		iris.WithoutStartupLog,
	)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to start application")
	}
}
