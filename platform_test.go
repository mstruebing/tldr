package tldr

import (
	"testing"
	"time"

	"github.com/mstruebing/tldr/cache"
	"github.com/stretchr/testify/require"
)

const (
	remoteURL = "https://tldr.sh/assets/tldr.zip"
	ttl       = time.Hour * 24 * 7
)

func TestCurrentPlattform(t *testing.T) {
	currentPlattform := CurrentPlatform("linux")

	if currentPlattform != "linux" {
		t.Error("Expected linux, got ", currentPlattform)
	}

	currentPlattform = CurrentPlatform("LINUX")
	if currentPlattform != "linux" {
		t.Error("Expected linux, got ", currentPlattform)
	}

	currentPlattform = CurrentPlatform("darwin")
	if currentPlattform != "osx" {
		t.Error("Expected osx, got ", currentPlattform)
	}

	currentPlattform = CurrentPlatform("DARWIN")
	if currentPlattform != "osx" {
		t.Error("Expected osx, got ", currentPlattform)
	}

	currentPlattform = CurrentPlatform("sunos")
	if currentPlattform != "sunos" {
		t.Error("Expected sunos, got ", currentPlattform)
	}

	currentPlattform = CurrentPlatform("SUNOS")
	if currentPlattform != "sunos" {
		t.Error("Expected sunos, got ", currentPlattform)
	}

	currentPlattform = CurrentPlatform("windows")
	if currentPlattform != "windows" {
		t.Error("Expected windows, got ", currentPlattform)
	}

	currentPlattform = CurrentPlatform("WINDOWS")
	if currentPlattform != "windows" {
		t.Error("Expected windows, got ", currentPlattform)
	}
}

func TestAvailablePlatforms(t *testing.T) {
	tests := []struct {
		name    string
		current string
		want    []string
		wantErr bool
	}{
		{
			name:    "linux",
			current: "linux",
			want:    []string{"android", "common", "linux", "osx", "sunos", "windows"},
		},
		{
			name:    "stuff",
			current: "stuff",
			want:    []string{"android", "common", "linux", "osx", "sunos", "windows", "stuff"},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repository, err := cache.NewRepository(remoteURL, ttl)
			require.NoError(t, err, "NewRepository() error %v", err)

			got, err := AvailablePlatforms(repository, tt.current)
			require.NoError(t, err, "AvailablePlatforms() error %v", err)
			require.ElementsMatch(t, tt.want, got, "expected available platforms %s, got %s", tt.want, got)
		})
	}
}
