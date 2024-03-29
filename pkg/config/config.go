package config

import (
	"github.com/joho/godotenv"
	"github.com/traefik/paerser/env"
	"os"

	"github.com/traefik/paerser/file"
)

// Parse takes pathes to a config files, an environment
// variable prefix and a default instance of T whichs
// values will be used when not set otherwise.
//
// Each step overwrites values set in a previous step.
// The config is loaded in the following priority:
//   - default state
//   - environment variables
//   - passed cfgFiles
//
// The cfg file can be in the format YAML, TOML, INI
// or JSON.
func Parse[T any](cfgFiles []string, envPrefix string, def ...T) (cfg T, err error) {
	if len(def) != 0 {
		cfg = def[0]
	}

	if envPrefix != "" {
		godotenv.Load()
		err = env.Decode(os.Environ(), envPrefix, &cfg)
		if err != nil {
			return cfg, err
		}
	}

	for _, cfgFile := range cfgFiles {
		err = file.Decode(cfgFile, &cfg)
		if err != nil && !os.IsNotExist(err) {
			return cfg, err
		}
	}

	return cfg, err
}
