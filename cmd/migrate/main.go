package main

import (
	"fmt"
	"os"

	"github.com/banditml/goat/app"
	"github.com/banditml/goat/model"
	"github.com/jinzhu/gorm"
	"go.uber.org/fx"
)

func main() {
	db := new(gorm.DB)
	fx.New(fx.Options(
		app.BaseModule,
		app.DBModule,
		fx.Populate(&db),
	))

	if err := db.AutoMigrate(new(model.Campaign)).Error; err != nil {
		fmt.Fprintf(os.Stderr, "failed to migrate Campaign: %s", err.Error())
	}
	fmt.Println("success")
}
