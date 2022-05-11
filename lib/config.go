package lib

import (
	"encoding/json"
	"errors"
	"os"
	"path"

	"github.com/pwuersch/ggt/lib/globals"
)

const (
	AppName          = "ggt-cli"
	ProfilesDirMode  = os.FileMode(0700)
	ProfilesFileMode = os.FileMode(0600)
)

func ReadProfiles() ([]globals.Profile, error) {
	err := ensureDirCreated()
	if err != nil {
		return []globals.Profile{}, err
	}

	profilePath, err := GetProfilePath()
	if err != nil {
		return []globals.Profile{}, err
	}

	profiles := []globals.Profile{}

	content, err := os.ReadFile(profilePath)
	if err != nil {
		return []globals.Profile{}, err
	}

	return profiles, json.Unmarshal(content, &profiles)
}

func GetProfilePath() (string, error) {
	profileDir, err := getConfigDir()
	if err != nil {
		return "", err
	}

	return path.Join(profileDir, "profiles.json"), nil
}

func getConfigDir() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	return path.Join(configDir, AppName), nil
}

func storeProfiles(profiles []globals.Profile) error {
	err := ensureDirCreated()
	if err != nil {
		return err
	}

	content, err := json.MarshalIndent(profiles, "", "\t")
	if err != nil {
		return err
	}

	profilePath, err := GetProfilePath()
	if err != nil {
		return err
	}

	return os.WriteFile(profilePath, content, ProfilesFileMode)
}

func ensureDirCreated() error {
	profileDir, err := getConfigDir()
	if err != nil {
		return err
	}

	info, err := os.Stat(profileDir)
	if err != nil || !info.IsDir() {
		err = os.MkdirAll(profileDir, ProfilesDirMode)
		if err != nil {
			return err
		}
	}

	return nil
}

func AddProfile(newProfile globals.Profile) error {
	profiles, _ := ReadProfiles()

	// check for existing profile name
	for _, profile := range profiles {
		if profile.Name == newProfile.Name {
			return errors.New("a profile with this name already exists")
		}
	}

	profiles = append(profiles, newProfile)

	return storeProfiles(profiles)
}

func RemoveProfile(profiles []globals.Profile, index int) error {
	profiles = append(profiles[:index], profiles[index+1:]...)

	return storeProfiles(profiles)
}
