// Package route routes paths to handlers.
package route

import (
	"go.uber.org/fx"
)

var Module = fx.Invoke(Register)
