// Package gormfx initializes and automigrates Gorm DB connections.
package gormfx

import (
	"context"
	"fmt"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/banditml/goat/envfx"
)

var Module = fx.Options(
	fx.Provide(NewDB),
)

type Params struct {
	fx.In

	Env       *envfx.Env
	Logger    *zap.Logger
	Lifecycle fx.Lifecycle
}

func NewDB(p Params) (*gorm.DB, error) {
	var db *gorm.DB
	var err error
	if p.Env.IsRelease() {
		username := os.Getenv("PGUSER")
		password := os.Getenv("PGPASSWORD")
		dbName := os.Getenv("PGDATABASE")
		dbHost := os.Getenv("PGHOST")

		dbURI := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=require password=%s", dbHost, username, dbName, password)
		db, err = gorm.Open("postgres", dbURI)
	} else {
		db, err = gorm.Open("sqlite3", "file::memory:?cache=shared")
	}
	if err != nil {
		return nil, err
	}
	p.Lifecycle.Append(fx.Hook{
		OnStop: func(c context.Context) error {
			return db.Close()
		},
	})
	return db, nil
}
