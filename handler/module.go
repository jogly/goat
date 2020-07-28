// Package handler has all the application business logic handlers for the
// routes defined in `route`.
package handler

import (
	"go.uber.org/fx"

	"github.com/banditml/goat/route"
)

var Module = fx.Options(
	route.Fx(NewCampaignHandler),
)
