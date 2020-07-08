// Package app collects all the dependencies of the main app.
package app

import (
	"context"
	"fmt"
	"net/http"

	"github.com/banditml/goat/envfx"
	"github.com/banditml/goat/ginfx"
	"github.com/banditml/goat/route"
	"github.com/banditml/goat/zapfx"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var Module = fx.Options(
	// Order of Modules does not matter! Thanks FX.
	ginfx.Module,
	zapfx.Module,
	envfx.Module,
	route.Module,
)

// Params is an FX-convention for declaring everything a function or struct
// depends on.  In this case, the `Start` function depends on the
// `Environment`, `Lifecycle`, `Engine`, and a `Logger`.
//
// `fx.In` is struct-embedding
// (https://golang.org/doc/effective_go.html#embedding) and it tells the FX dep
// system to populate everything in this struct with instances of the same
// type.
type Params struct {
	fx.In

	Environment *envfx.Env
	Lifecycle   fx.Lifecycle
	Engine      *gin.Engine

	Logger *zap.Logger
}

func Start(p Params) {
	// create a new http server that uses Gin as the request handler
	// on the port provided via the environment.
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", p.Environment.Port),
		Handler: p.Engine,
	}
	// Add some hooks to the app lifecycle. Hooks because that guarantees all
	// the dependencies have had a chance to start (namely the routes still
	// need to be defined) before we start accepting requests.
	p.Lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			// Kick off the main goroutine.  `ListenAndServe` is a blocking
			// function call.  It will never stop until we tell it to or it is
			// killed by the system.  We tell it to shut up in the Stop hook
			// below.
			go func() {
				// it is GENERALLY not a great idea to use closure variables
				// for goroutines due to race conditions.  It is OK here
				// because there is 1 and only ever 1 of this goroutine.
				if err := server.ListenAndServe(); err != nil {
					p.Logger.Fatal("listening failed", zap.Error(err))
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			// We got a signal to shut down so instead of just dying like an
			// asshole we'll let the http server finish wrapping up anything it
			// might be doing first.  FX has a timeout on shutting down so if
			// the server doesn't shut down nice, we'll still stop eventually.
			p.Logger.Info("gracefully shutting down")
			var err error
			if err = server.Shutdown(ctx); err != nil {
				p.Logger.Error("failed to shut down gracefully", zap.Error(err))
			}
			return err
		},
	})
}

func New() *fx.App {
	return fx.New(
		Module,
		// Invocations and Extractions always happen after everything is
		// provided, so even these could be in the beginning of the list but
		// are called in the order they are defined.
		fx.Invoke(Start),
	)
}
