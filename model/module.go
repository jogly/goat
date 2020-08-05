package model

import (
	"context"
	"fmt"

	"github.com/jinzhu/gorm"
	"go.uber.org/fx"
)

var Module = fx.Invoke(func(lc fx.Lifecycle, db *gorm.DB) {
	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			err := db.AutoMigrate(Campaign{}).Error
			if err != nil {
				return fmt.Errorf("failed to migrate db: %s", err.Error())
			}
			return nil
		},
	})
})
