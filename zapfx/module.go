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
		l, err = zap.NewProduction()
		l = l.With(
			zap.String("dd.env", e.Env),
			zap.String("dd.version", e.Version),
		)
	} else {
		l, err = zap.NewDevelopment()
	}
	if err != nil {
		panic(err)
	}
	return l
})
