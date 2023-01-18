package config

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/traefik/paerser/env"
	"github.com/traefik/paerser/file"
)

// Parse rakes a path to a config file, an environment
// variable prefix and a default instance of T whichs
// values will be used when not set otherwise.
//
// Each step overwrites values set in a previous step.
// The config is loaded in the following priority:
//   - default state
//   - environment variables
//   - passed cfgFile
//
// The cfg file can be in the format YAML, TOML, INI
// or JSON.
func Parse[T any](cfgFile string, envPrefix string, def ...T) (cfg T, err error) {
	if len(def) != 0 {
		cfg = def[0]
	}

	godotenv.Load()
	err = env.Decode(os.Environ(), envPrefix, &cfg)

	if cfgFile != "" {
		err = file.Decode(cfgFile, &cfg)
		if err != nil && !os.IsNotExist(err) {
			return cfg, err
		}
	}

	return cfg, err
}
