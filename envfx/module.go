// Package envfx surfaces information that comes from an Environment in a typed and managed format.
package envfx

import (
	"strings"

	env "github.com/Netflix/go-env"
	"github.com/go-playground/validator/v10"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// Environment types Env.Tag
type Environment = string

const (
	// Test is the test environment.
	Test Environment = "test"
	// Development is the local environment.  There should be minimal
	// dependency on this.
	Development Environment = "development"
	// Production is a singular "released" environment.  Staging is not
	// distinct from Production.
	Production Environment = "production"
)

// Env is a collection of configuration sourced from environment variables.
type Env struct {
	Tag  Environment `env:"ENV" validate:"gt=0"`
	Port int         `env:"PORT" validate:"gt=0"`

	Extras env.EnvSet
}

func (e *Env) IsProduction() bool {
	return e.Tag == Production
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
	if err = validator.New().Struct(e); err != nil {
		l.Fatal(err.Error())
	}
	switch strings.ToLower(e.Tag) {
	case "test":
		e.Tag = Test
	case "dev", "devel", "development":
		e.Tag = Development
	case "prod", "production":
		e.Tag = Production
	}

	l.Debug("successfully loaded environment", zap.Any("context", e))
	return e
}
