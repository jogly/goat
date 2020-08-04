package zapfx

import (
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/banditml/goat/envfx"
)

var Module = fx.Provide(func(e *envfx.Env) *zap.Logger {
	var l *zap.Logger
	var err error
	if e.IsRelease() {
		// Production is a pre-configured JSON-based logger capped to INFO level.
		l, err = zap.NewProduction()
		l = l.With(
			zap.String("dd.env", e.Env),
			zap.String("dd.version", e.Version),
		)
	} else {
		// Development is a pre-configured TEXT-based logger capped to DEBUG level.
		l, err = zap.NewDevelopment()
	}
	if err != nil {
		// Any issues is ridiculously improbable but means the app can't run or
		// log so bail hard.
		panic(err)
	}
	return l
})
