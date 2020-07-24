package model

import (
	"context"

	"github.com/jinzhu/gorm"
	"go.uber.org/fx"
)

var Module = fx.Invoke(func(lc fx.Lifecycle, db *gorm.DB) {
	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			err := db.AutoMigrate(Campaign{}).Error
			return err
		},
	})
})
