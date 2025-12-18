package main

import (
	"fmt"
	"runtime"
)

// Version information - updated during release builds
var (
	Version   = "dev"     // Version number (e.g., "1.0.0")
	BuildDate = "unknown" // Build date
	GitCommit = "unknown" // Git commit hash
	GoVersion = runtime.Version()
)

// GetVersionString returns a formatted version string
func GetVersionString() string {
	return fmt.Sprintf("AlternateDNS v%s", Version)
}

// GetFullVersionInfo returns detailed version information
func GetFullVersionInfo() string {
	return fmt.Sprintf("**Version:** %s\n\n**Build Date:** %s\n\n**Git Commit:** %s\n\n**Go Version:** %s\n\n**Platform:** %s/%s",
		Version, BuildDate, GitCommit, GoVersion, runtime.GOOS, runtime.GOARCH)
}
