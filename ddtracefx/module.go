package ddtracefx

import (
	"context"

	"go.uber.org/fx"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"

	"github.com/banditml/goat/envfx"
)

var Module = fx.Invoke(StartTracer)

type Params struct {
	fx.In

	env *envfx.Env
	lc  fx.Lifecycle
}

func StartTracer(p Params) {
	enabled := !p.env.IsRelease()
	if !enabled {
		return
	}
	p.lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			tracer.Start(
				tracer.WithEnv(p.env.Env),
				tracer.WithServiceVersion(p.env.Version),
			)
			return nil
		},
		OnStop: func(context.Context) error {
			tracer.Stop()
			return nil
		},
	})
}
