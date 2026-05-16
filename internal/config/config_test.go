package config

import (
	"os"
	"path/filepath"
	"testing"
)

func useTempHome(t *testing.T) string {
	t.Helper()

	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home)
	return home
}

func TestLoadUsesHostedDefaultWhenConfigMissing(t *testing.T) {
	useTempHome(t)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.APIURL != "https://api.sendr.app" {
		t.Fatalf("APIURL = %q, want hosted default", cfg.APIURL)
	}
}

func TestClearRemovesTokenAndPreservesAPIKey(t *testing.T) {
	home := useTempHome(t)
	cfg := &Config{
		APIURL: "https://api.sendr.app",
		Token:  "jwt-token",
		APIKey: "mk_live_123456789",
	}

	if err := cfg.Clear(); err != nil {
		t.Fatalf("Clear() error = %v", err)
	}

	loaded, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if loaded.Token != "" {
		t.Fatalf("Token = %q, want empty", loaded.Token)
	}
	if loaded.APIKey != cfg.APIKey {
		t.Fatalf("APIKey = %q, want %q", loaded.APIKey, cfg.APIKey)
	}

	configPath := filepath.Join(home, ".config", "sendr", "config.json")
	info, err := os.Stat(configPath)
	if err != nil {
		t.Fatalf("Stat(%q) error = %v", configPath, err)
	}
	if info.IsDir() {
		t.Fatalf("%q is a directory, want file", configPath)
	}
}
