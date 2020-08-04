// Package envfx surfaces information that comes from an Environment in a typed and managed format.
package envfx

import (
	"os"
	"strings"

	env "github.com/Netflix/go-env"
	"github.com/go-playground/validator/v10"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// Mode indicates what mode the app in running in. This should only be used
// when setting up integration stubs/replacements.
type Mode = string

const (
	// Local mode is set when running the application locally.
	Local Mode = "local"
	// Test mode is set when running application tests.
	Test Mode = "test"
	// Release mode is set when running the application remotely, including
	// `dev` and `prod` environments.  `Env.Env` indicates which _environment_,
	// such as `dev` or `prod`.
	Release Mode = "release"

	DefaultEnv = "local"
)

// Env is a collection of configuration sourced from environment variables.
type Env struct {
	// Env can be local, test, dev, or prod.
	Env  string `env:"ENV" validate:"required"`
	Mode Mode

	// Port is the port the application is listening on.
	Port int `env:"PORT"`

	// Version is the (likely) the git SHA.
	Version string `env:"VERSION"`

	// Postgres has the configuration for connecting to a PG db.
	// Not used in local and test modes.
	Postgres struct {
		User     string `env:"PGUSER"`
		Password string `env:"PGPASSWORD"`
		Host     string `env:"PGHOST"`
		Database string `env:"PGDATABASE"`
	}

	// Everything else for debugging purposes.  `json:"-"` to never serialize,
	// it has secrets.
	Extras env.EnvSet `json:"-"`
}

// IsRelease is true if the application is running in the cloud.
func (e *Env) IsRelease() bool {
	return e.Mode == Release
}

// Module is a FX Module
var Module = fx.Provide(NewEnv)

// NewEnv creates an Env from the OS environment.  `overrides` allow testing
// settings without modifying the global OS environment.
func NewEnv(overrides ...env.ChangeSet) *Env {
	// create a temporary console logger to use because the application logger
	// is dependent on the struct we're creating here.
	l, _ := zap.NewDevelopment()
	// Load an Env struct with variables from the OS environment.
	e := &Env{}
	envSet, err := env.EnvironToEnvSet(os.Environ())
	if err != nil {
		// Any errors here mean we're unstable, log and kill the app,
		// preventing startup.  The app won't start listening until all
		// dependencies report no errors so it is safe to do so immediately.
		l.Fatal(err.Error())
	}
	// Overrides allow testing environment settings without global impact.
	for _, override := range overrides {
		envSet.Apply(override)
	}
	err = env.Unmarshal(envSet, e)
	if err != nil {
		l.Fatal(err.Error())
	}
	e.Extras = envSet

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
	case "dev", "prod":
		e.Mode = Release
	}

	l.Debug("successfully loaded environment", zap.Any("context", e))
	return e
}
