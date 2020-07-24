package model

import (
	"github.com/jinzhu/gorm"
)

func Migrate(db *gorm.DB) error {
	err := db.Debug().AutoMigrate(new(Campaign)).Error
	return err
}
