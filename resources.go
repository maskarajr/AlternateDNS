package main

import (
	_ "embed"
	"os"
	"path/filepath"
)

//go:embed icon.ico
var embeddedIcon []byte

//go:embed default_config.yaml
var embeddedDefaultConfig []byte

// getEmbeddedIcon returns the embedded icon bytes
func getEmbeddedIcon() []byte {
	return embeddedIcon
}

// getEmbeddedDefaultConfig returns the embedded default config
func getEmbeddedDefaultConfig() []byte {
	return embeddedDefaultConfig
}

// getConfigPath returns the path to config.yaml in the executable directory
func getConfigPath() (string, error) {
	exePath, err := os.Executable()
	if err != nil {
		return "", err
	}
	exeDir := filepath.Dir(exePath)
	return filepath.Join(exeDir, "config.yaml"), nil
}
