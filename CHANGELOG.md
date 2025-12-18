# Changelog

All notable changes to AlternateDNS will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Hours and minutes interval configuration (replaces hours-only)
- DNS restoration to automatic/DHCP when service is stopped
- Version information display in Settings tab with formatted details
- Automated GitHub releases with multi-platform builds
- Version embedding in executables

### Fixed
- Interface selector now correctly handles multi-word interface names (e.g., "Radmin VPN")
- Interface selector remains enabled after stopping service when multiple interfaces are available

## [1.0.0] - 2025-12-18

### 1. User Interface Improvements

#### 1.1 **Modern Desktop GUI**

* **Description**: Complete GUI implementation using Fyne framework with tabbed interface for easy navigation
* **Implementation**: Added Status, DNS Servers, DNS Tester, Settings, and Logs tabs with modern UI components
* **Impact**: Significantly improved user experience with intuitive graphical interface instead of command-line only

#### 1.2 **System Tray Integration**

* **Description**: Seamless integration between GUI window and system tray for background operation
* **Implementation**: System tray menu with quick access to window, DNS change, and quit options
* **Impact**: Better user control with ability to minimize to tray while keeping service running

#### 1.3 **Real-time Status Display**

* **Description**: Live updates showing current DNS, service status, and countdown timer
* **Implementation**: Automatic refresh of status information with countdown to next DNS change
* **Impact**: Users always know current state and when next change will occur

#### 1.4 **Hours and Minutes Interval Configuration**

* **Description**: Granular interval control with separate hours and minutes inputs
* **Implementation**: Two separate input fields for hours and minutes, stored as total minutes in config with backward compatibility for old hours-only configs
* **Impact**: Users can set precise intervals like 30 minutes or 2 hours 15 minutes, enabling better control for shorter testing periods

#### 1.5 **DNS Restoration on Service Stop**

* **Description**: Automatic restoration of DNS settings to automatic/DHCP when service is stopped
* **Implementation**: `restoreDNS()` function that resets DNS to automatic on Windows, Linux, and macOS
* **Impact**: Prevents custom DNS from persisting after stopping the service, ensuring clean state restoration

#### 1.6 **Version Information Display**

* **Description**: Comprehensive version information displayed in Settings tab
* **Implementation**: Version, build date, git commit, Go version, and platform displayed with formatted layout (one detail per line)
* **Impact**: Users can easily identify the exact build version and system information for troubleshooting

### 2. Smart DNS Switching

#### 2.1 **Latency-Based DNS Comparison**

* **Description**: Automatic testing and comparison of all DNS servers before switching
* **Implementation**: Tests all configured DNS servers and selects the one with lowest latency and highest success rate
* **Impact**: Ensures optimal DNS performance by automatically selecting the best available server

#### 2.2 **Configurable Test Domains**

* **Description**: Customizable list of domains for DNS benchmarking
* **Implementation**: Users can specify test domains in config.yaml for more accurate latency testing
* **Impact**: Better DNS performance measurement tailored to user's actual usage patterns

#### 2.3 **Manual vs Automatic Switching**

* **Description**: Different behavior for manual DNS changes vs automatic timer-based changes
* **Implementation**: Manual "Change DNS Now" forces immediate switch, automatic changes use smart switching logic
* **Impact**: Users have control when needed while benefiting from smart switching automatically

### 3. DNS Testing Features

#### 3.1 **Built-in DNS Tester**

* **Description**: Dedicated tab for benchmarking all configured DNS servers
* **Implementation**: Real-time latency measurement and success rate calculation with automatic sorting
* **Impact**: Users can easily compare DNS server performance before configuring rotation

#### 3.2 **Performance Metrics**

* **Description**: Comprehensive performance data including average latency and success rates
* **Implementation**: Tests multiple domains and calculates statistics for each DNS server
* **Impact**: Data-driven DNS server selection based on actual performance metrics

### 4. Build and Distribution

#### 4.1 **Portable Single-Executable Builds**

* **Description**: Self-contained executables with embedded resources
* **Implementation**: Icon and default config embedded using Go embed, no external dependencies
* **Impact**: Easy distribution and deployment - just copy the executable

#### 4.2 **Cross-Platform Support**

* **Description**: Native builds for Windows, Linux, and macOS
* **Implementation**: Automated builds for multiple architectures (amd64, arm64) with platform-specific optimizations
* **Impact**: Works seamlessly across all major operating systems

#### 4.3 **Windows GUI Without Console**

* **Description**: Clean Windows application without console window
* **Implementation**: Windows-specific build flags to hide console, proper icon embedding
* **Impact**: Professional Windows application appearance

### 5. Code Improvements

#### 5.1 **Thread-Safe State Management**

* **Description**: Proper synchronization for concurrent access to application state
* **Implementation**: Mutex-protected state structure with safe getter/setter methods
* **Impact**: Reliable operation in multi-threaded GUI environment

#### 5.2 **Enhanced Error Handling**

* **Description**: Comprehensive error handling with user-friendly messages
* **Implementation**: Error logging to GUI Logs tab, non-blocking error dialogs
* **Impact**: Better user experience with clear error messages and debugging information

#### 5.3 **Improved Code Organization**

* **Description**: Separation of concerns with dedicated files for GUI, state, and DNS testing
* **Implementation**: Modular structure with clear responsibilities for each component
* **Impact**: Easier maintenance and future development

#### 5.4 **Version System Integration**

* **Description**: Comprehensive version tracking and display system
* **Implementation**: Version constants embedded via build flags, displayed in window title and Settings tab
* **Impact**: Easy identification of build version, commit, and build date for debugging and support

### 6. Bug Fixes

#### 6.1 **Interface Selector Improvements**

* **Description**: Fixed handling of multi-word network interface names
* **Implementation**: Changed from `strings.Fields()` to `strings.Split()` to preserve interface names with spaces
* **Impact**: Interfaces like "Radmin VPN" and "vEthernet (Default Switch)" are now correctly displayed as single options

#### 6.2 **Interface Selector State Management**

* **Description**: Interface selector remains enabled after stopping service when multiple interfaces are available
* **Implementation**: Updated `updateStatusDisplay()` to re-enable selector based on interface count
* **Impact**: Better user experience with consistent interface selection availability

## Original Features

- Periodic DNS rotation
- Cross-platform support
- System tray integration
- Startup registration
- Notification support
- YAML configuration
