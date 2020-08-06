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

	Env       *envfx.Env
	Lifecycle fx.Lifecycle
}

func StartTracer(p Params) {
	enabled := p.Env.IsRelease()
	if !enabled {
		return
	}
	p.Lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			tracer.Start(
				tracer.WithEnv(p.Env.Env),
				tracer.WithServiceVersion(p.Env.Version),
			)
			return nil
		},
		OnStop: func(context.Context) error {
			tracer.Stop()
			return nil
		},
	})
}
