// Package envfx surfaces information that comes from an Environment in a typed and managed format.
package envfx

import (
	"strings"

	env "github.com/Netflix/go-env"
	"github.com/go-playground/validator/v10"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// Mode lets consumers know what mode the app in running in.
type Mode = string

const (
	// Test is the test Mode.
	Test Mode = "test"
	// Local is the local Mode.  There should be minimal
	// dependency on this.
	Local Mode = "local"
	// Production is a singular "released" Mode.  Staging is not
	// distinct from Production.
	Release Mode = "release"

	DefaultEnv = "local"
)

// Env is a collection of configuration sourced from environment variables.
type Env struct {
	Env  string `env:"ENV" validate:"required"`
	Mode Mode
	Port int `env:"PORT"`

	Postgres struct {
		User     string `env:"PGUSER"`
		Password string `env:"PGPASSWORD"`
		Host     string `env:"PGHOST"`
		Database string `env:"PGDATABASE"`
	}

	Extras env.EnvSet `json:"-"`
}

func (e *Env) IsRelease() bool {
	return e.Mode == Release
}

// Module is a FX Module
var Module = fx.Provide(NewEnv)

// NewEnv creates an Env from the OS environment.
func NewEnv() *Env {
	l, _ := zap.NewDevelopment()
	e := &Env{}
	es, err := env.UnmarshalFromEnviron(e)
	if err != nil {
		l.Fatal(err.Error())
	}
	e.Extras = es

	if e.Env == "" {
		e.Env = DefaultEnv
	}
	if err = validator.New().Struct(e); err != nil {
		l.Fatal(err.Error())
	}
	switch strings.ToLower(e.Env) {
	case "test":
		e.Mode = Test
	case "local":
		e.Mode = Local
	case "dev", "development", "stage", "staging", "prod", "production":
		e.Mode = Release
	}

	l.Debug("successfully loaded environment", zap.Any("context", e))
	return e
}
