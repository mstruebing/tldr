package cache

import (
	"strings"
	"testing"
)

func TestCacheDir(t *testing.T) {
	cacheDir, err := cacheDir()

	if err != nil || !strings.HasSuffix(cacheDir, ".tldr") {
		t.Error("Expected linux, got ", err, cacheDir)
	}
}
