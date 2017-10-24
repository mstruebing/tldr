package tldr

import (
	"strings"
)

// CommonPlatform is the common platform.
const CommonPlatform = "common"

// CurrentPlatform returns the platform name of the current system.
func CurrentPlatform(current string) string {
	var os string = current
	if strings.ToLower(current) == "darwin" {
		os = "osx"
	}
	return strings.ToLower(os)
}

// AvailablePlatforms returns all the available platforms that are supported.
func AvailablePlatforms(r Repository, current string) ([]string, error) {
	platforms, err := r.AvailablePlatforms()
	if err != nil {
		return nil, err
	}

	currentPlattform := CurrentPlatform(current)
	var currentFound, commonFound bool
	for _, p := range platforms {
		if p == currentPlattform {
			currentFound = true
		} else if p == CommonPlatform {
			commonFound = true
		}
	}

	if !currentFound {
		platforms = append(platforms, CurrentPlatform(current))
	}
	if !commonFound {
		platforms = append(platforms, CommonPlatform)
	}
	return platforms, nil
}
