# Changelog

All notable changes to this enhanced version of AlternateDNS will be documented in this file.

## Enhanced Features (This Fork)

### GUI and User Experience
- Added modern desktop GUI using Fyne framework
- Implemented tabbed interface (Status, DNS Servers, DNS Tester, Settings, Logs)
- Added system tray integration with GUI window management
- Improved log display with better readability
- Added real-time status updates and countdown timer

### Smart DNS Switching
- Implemented latency-based DNS comparison before automatic switching
- Added configurable test domains for DNS benchmarking
- Smart switching only applies to automatic timer-based changes
- Manual "Change DNS Now" button forces immediate change (bypasses latency testing)

### DNS Testing Features
- Built-in DNS Tester tab for benchmarking all configured DNS servers
- Real-time latency measurement and success rate calculation
- Automatic sorting by performance
- Support for custom test domains

### Build and Distribution
- Portable single-executable builds
- Embedded resources (icon, default config)
- Cross-platform build scripts (build.bat, build.sh)
- Windows GUI build without console window

### Code Improvements
- Thread-safe application state management
- Better error handling and logging
- Improved code organization (separated GUI, state, and DNS testing logic)
- Enhanced configuration management

## Original Features (by MaxIsJoe)

- Periodic DNS rotation
- Cross-platform support
- System tray integration
- Startup registration
- Notification support
- YAML configuration

---

**Note:** This changelog documents enhancements made in this fork. All original functionality and code structure are credited to [MaxIsJoe](https://github.com/MaxIsJoe).

