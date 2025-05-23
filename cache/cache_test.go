package cache_test

import (
	"os"
	"testing"

	"github.com/fischyn/omfetch/cache"
)

func TestSaveLoadCPUCache(t *testing.T) {
	const filename = "testcache.omfetch"
	data := []byte("test data")

	err := cache.SaveCPUCache(filename, data)
	if err != nil {
		t.Fatalf("SaveCPUCache failed: %v", err)
	}
	defer os.Remove(filename)

	loaded, err := cache.LoadCPUCache(filename)
	if err != nil {
		t.Fatalf("LoadCPUCache failed: %v", err)
	}

	if string(loaded) != string(data) {
		t.Errorf("Loaded data mismatch: got %q, want %q", loaded, data)
	}

	info, err := os.Stat(filename)
	if err != nil {
		t.Fatalf("Failed to stat file: %v", err)
	}
	mode := info.Mode().Perm()
	if mode != cache.ReadOnlyPerm {
		t.Errorf("File permission mismatch: got %v, want %v", mode, cache.ReadOnlyPerm)
	}
}
