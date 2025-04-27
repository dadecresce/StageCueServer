package config

import (
    "os"

    "github.com/BurntSushi/toml"
)

// Config holds server configuration loaded from TOML.
type Config struct {
    // Address the server listens on, e.g. ":8080"
    Address string `toml:"address"`
    // LogLevel can be "debug","info","warn","error"
    LogLevel string `toml:"log_level"`
}

// Parse reads the given path and unmarshals TOML into a Config.
func Parse(path string) (*Config, error) {
    cfg := &Config{
        Address: ":8080",
        LogLevel: "info",
    }
    f, err := os.Open(path)
    if err != nil {
        // not fatal: return defaults
        return cfg, nil
    }
    defer f.Close()
    if _, err := toml.NewDecoder(f).Decode(cfg); err != nil {
        return nil, err
    }
    return cfg, nil
}
