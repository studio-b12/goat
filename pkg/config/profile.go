package config

import (
	"fmt"
	"github.com/studio-b12/goat/pkg/engine"
	"github.com/zekrotja/rogu/log"
	"os"
	"path/filepath"
	"strings"
)

// LoadProfiles tries to find 'profile.*' configuration files in the 'goat' directory in the
// users home configuration directory, if existent. They will be parsed using config.Parse
// to a map of engine.State. After that, the 'default' profile is taken from the map, if
// existent, and the defined state is merged with the given state. After that, all profiles
// defined in profileNames is taken from the map and defined states are merged with the passed
// profile as well.
//
// If profileNames contains values, this function will return an error if one of the given
// profiles is not specified or if there is no goat configuration path.
func LoadProfiles(profileNames []string, state engine.State) (err error) {
	if state == nil {
		return nil
	}

	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		log.Warn().Tag("Profiles").Err(err).Msg("Failed getting user config directory")
		return nil
	}

	pth := filepath.Join(userConfigDir, "goat")
	_, err = os.Stat(pth)
	if os.IsNotExist(err) {
		if len(profileNames) > 0 {
			return fmt.Errorf("no profiles(.*) file can be found in your config directory (%s)", pth)
		}
		return nil
	}
	if err != nil {
		return err
	}

	dirEntries, err := os.ReadDir(pth)
	if err != nil {
		return err
	}

	var profileDirs []string
	for _, dirEntry := range dirEntries {
		if !dirEntry.IsDir() && strings.HasPrefix(dirEntry.Name(), "profiles") {
			profileDirs = append(profileDirs, filepath.Join(pth, dirEntry.Name()))
		}
	}

	profiles, err := Parse[map[string]any](profileDirs, "")
	if err != nil {
		return err
	}

	if err = getAndMergeProfile(profiles, state, "default", true); err != nil {
		return err
	}

	for _, profileName := range profileNames {
		if err = getAndMergeProfile(profiles, state, profileName, false); err != nil {
			return err
		}
	}

	return nil
}

func getAndMergeProfile(profiles map[string]any, state engine.State, profileName string, optional bool) error {
	v, ok := profiles[profileName]
	if !ok {
		if !optional {
			return fmt.Errorf("no profile found with name '%s'", profileName)
		}
		return nil
	}

	s, ok := v.(map[string]any)
	if !ok {
		return fmt.Errorf("value of profile field '%s' is not a state map", profileName)
	}

	state.Merge(s)

	return nil
}
