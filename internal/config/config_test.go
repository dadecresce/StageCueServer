package config_test

import (
	"os"
	"testing"

	"github.com/StageCue/StageCueServer/internal/config"
)

const sample = `
address   = ":9999"
log_level = "debug"
`

func TestParse_Defaults(t *testing.T) {
	cfg, err := config.Parse("non-esiste.toml") // deve cadere in default
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Address != ":8080" || cfg.LogLevel != "info" {
		t.Errorf("defaults not applied: %+v", cfg)
	}
}

func TestParse_File(t *testing.T) {
	f, _ := os.CreateTemp(t.TempDir(), "cfg-*.toml")
	f.WriteString(sample)
	f.Close()

	cfg, err := config.Parse(f.Name())
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if cfg.Address != ":9999" || cfg.LogLevel != "debug" {
		t.Errorf("values not parsed: %+v", cfg)
	}
}
