package handler

import (
	"github.com/banditml/goat/route"
	"go.uber.org/fx"
)

var Module = fx.Options(
	route.Fx(NewCampaignHandler),
)
