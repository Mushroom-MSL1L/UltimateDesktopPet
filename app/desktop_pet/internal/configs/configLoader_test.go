package configs

import (
	"os"
	"path/filepath"
	"testing"
)

type testConfig struct {
	A int    `yaml:"a"`
	B string `yaml:"b"`
}

type testProvider struct{}

func (testProvider) DefaultConfig() *testConfig {
	return &testConfig{
		A: 1,
		B: "default",
	}
}

func TestLoadConfig_CreatesMissingFileWithDefaults(t *testing.T) {
	cfgPath := filepath.Join(t.TempDir(), "system.yaml")
	cfg := LoadConfig(cfgPath, testProvider{})

	if cfg.A != 1 || cfg.B != "default" {
		t.Fatalf("LoadConfig = %#v, want A=1 B=%q", cfg, "default")
	}
	if _, err := os.Stat(cfgPath); err != nil {
		t.Fatalf("expected config file created: %v", err)
	}
}

func TestLoadConfig_MergesWithExistingFile(t *testing.T) {
	cfgPath := filepath.Join(t.TempDir(), "system.yaml")
	_ = LoadConfig(cfgPath, testProvider{})

	if err := os.WriteFile(cfgPath, []byte("b: override\n"), 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}

	cfg := LoadConfig(cfgPath, testProvider{})
	if cfg.A != 1 || cfg.B != "override" {
		t.Fatalf("LoadConfig = %#v, want A=1 B=%q", cfg, "override")
	}
}

