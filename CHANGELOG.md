# Changelog

All notable changes to AlternateDNS will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Version information display in Settings tab
- Automated GitHub releases with multi-platform builds
- Version embedding in executables

## [1.0.0] - TBD

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

## Original Features

- Periodic DNS rotation
- Cross-platform support
- System tray integration
- Startup registration
- Notification support
- YAML configuration
