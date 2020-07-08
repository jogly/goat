package zapfx

import (
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var Module = fx.Provide(func() *zap.Logger {
	l, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	return l
})
