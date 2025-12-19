package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"github.com/gen2brain/beeep"
	"gopkg.in/yaml.v2"
)

type Config struct {
	DNSAddresses          []string `yaml:"dns_addresses"`
	RunOnStartup          bool     `yaml:"run_on_startup"`
	ChangeIntervalHours   int      `yaml:"change_interval_hours"`   // Deprecated: kept for backward compatibility
	ChangeIntervalMinutes int      `yaml:"change_interval_minutes"` // New: interval in minutes
	NotifyUser            bool     `yaml:"notify_user"`
	TestDomains           []string `yaml:"test_domains"` // Domains used for DNS latency testing
}

var config Config
var appIcon []byte

// Custom writer that redirects to appState logs
type logWriter struct{}

func (w logWriter) Write(p []byte) (n int, err error) {
	message := strings.TrimSpace(string(p))
	if message != "" {
		appState.AddLog(message)
		// Update GUI if available (safely check if GUI is initialized)
		if mainWindow != nil {
			fyne.Do(func() {
				updateLogsDisplay()
			})
		}
	}
	return len(p), nil
}

func main() {
	// Redirect standard log output to appState (all log.Printf, log.Println, etc. will go to Logs tab)
	log.SetOutput(logWriter{})

	// Check for debug flag
	if len(os.Args) > 1 && os.Args[1] == "--debug" {
		appState.SetDebugMode(true)
		appState.AddLog("Debug mode enabled via command line")
	}

	// Try to load icon from file, fallback to embedded
	appIcon = getIcon("icon.ico")
	if len(appIcon) == 0 {
		appIcon = getEmbeddedIcon()
	}

	// Read config
	err := readConfig()
	if err != nil {
		// Show error in GUI if possible, otherwise fatal
		log.Printf("Failed to read config: %v", err)
	}

	// Debug mode will be set via GUI or command line flag

	// Get active interfaces
	if runtime.GOOS == "windows" {
		ifaces, err := getActiveWindowsInterfaces()
		if err != nil {
			appState.AddLog(fmt.Sprintf("Warning: Failed to get active interfaces: %v", err))
		} else {
			appState.SetInterfaces(ifaces)
			appState.AddLog(fmt.Sprintf("Active interfaces: %v", ifaces))
		}
	}

	if runtime.GOOS == "darwin" {
		appState.AddLog("Note: macOS support is experimental")
	}

	// Set startup if configured
	if config.RunOnStartup {
		err = setRunOnStartup()
		if err != nil {
			appState.AddLog(fmt.Sprintf("Warning: Failed to set run on startup: %v", err))
		}
	}

	// Initialize GUI on main thread (Fyne requirement)
	// This will also set up the system tray using Fyne's native API
	setupGUI()

	// Run Fyne event loop (this blocks and handles GUI events)
	guiApp.Run()
}

// System tray menu items are now handled by Fyne's native desktop API in setupGUI()

func checkAdmin() error {
	if appState.GetDebugMode() {
		appState.AddLog(fmt.Sprintf("OS: %s", runtime.GOOS))
	}

	switch runtime.GOOS {
	case "windows":
		cmd := exec.Command("net", "session")
		err := cmd.Run()
		if err != nil {
			return fmt.Errorf("this program must be run as an administrator")
		}
	case "linux", "darwin":
		if os.Geteuid() != 0 {
			return fmt.Errorf("this program must be run as root")
		}
	default:
		return fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}
	return nil
}

func readConfig() error {
	configPath, err := getConfigPath()
	if err != nil {
		return fmt.Errorf("failed to get config path: %v", err)
	}

	data, err := os.ReadFile(configPath)
	if os.IsNotExist(err) {
		// Use embedded default config
		data = getEmbeddedDefaultConfig()
		if len(data) == 0 {
			// Fallback: generate default config
			err = generateDefaultConfig()
			if err != nil {
				return err
			}
			data, err = os.ReadFile(configPath)
			if err != nil {
				return err
			}
		} else {
			// Write embedded config to file
			err = os.WriteFile(configPath, data, 0644)
			if err != nil {
				return fmt.Errorf("failed to write default config: %v", err)
			}
		}
	} else if err != nil {
		return err
	}

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return err
	}

	// Handle backward compatibility: convert hours to minutes if needed
	if config.ChangeIntervalMinutes <= 0 {
		if config.ChangeIntervalHours > 0 {
			// Convert old hours format to minutes
			config.ChangeIntervalMinutes = config.ChangeIntervalHours * 60
			config.ChangeIntervalHours = 0 // Clear old value
		} else {
			// Default to 6 hours (360 minutes)
			config.ChangeIntervalMinutes = 360
		}
	}
	return nil
}

// changeDNS changes the DNS server. If forceChange is true, it switches to the next DNS
// without latency testing. If false, it tests all DNS servers and uses the one with lowest latency.
func changeDNS(forceChange bool) error {
	if len(config.DNSAddresses) == 0 {
		return fmt.Errorf("no DNS addresses specified in config")
	}

	_, currentIdx := appState.GetCurrentDNS()

	// If force change, skip latency testing and switch to next DNS immediately
	if forceChange {
		if currentIdx < 0 || currentIdx >= len(config.DNSAddresses) {
			currentIdx = 0
		}
		nextIndex := (currentIdx + 1) % len(config.DNSAddresses)
		nextDNS := config.DNSAddresses[nextIndex]
		appState.AddLog(fmt.Sprintf("Force changing DNS to %s", nextDNS))
		return applyDNS(nextDNS, nextIndex)
	}

	// Smart switching: Test all DNS servers and use the one with lowest latency
	// This applies to both service start and automatic timer-based changes
	testDomains := config.TestDomains
	if len(testDomains) == 0 {
		testDomains = defaultTestDomains
	}

	bestDNS, bestIdx := findBestDNS(config.DNSAddresses, testDomains)

	if bestIdx < 0 {
		// Fallback: use first DNS if all tests failed
		bestDNS = config.DNSAddresses[0]
		bestIdx = 0
		appState.AddLog(fmt.Sprintf("Using fallback DNS: %s", bestDNS))
	} else {
		// Check if we're switching from current DNS
		if currentIdx >= 0 && currentIdx < len(config.DNSAddresses) {
			currentDNS := config.DNSAddresses[currentIdx]
			if bestDNS != currentDNS {
				appState.AddLog(fmt.Sprintf("Switching from %s to %s (better performance)", currentDNS, bestDNS))
			} else {
				appState.AddLog(fmt.Sprintf("Keeping current DNS (%s) - it's the best performing", bestDNS))
			}
		} else {
			appState.AddLog(fmt.Sprintf("Setting DNS to %s (best performing)", bestDNS))
		}
	}

	return applyDNS(bestDNS, bestIdx)
}

// restoreDNS restores DNS settings to automatic/DHCP
func restoreDNS() error {
	var allErrors []string

	switch runtime.GOOS {
	case "windows":
		activeInterfaces, err := getActiveWindowsInterfaces()
		if err != nil {
			return err
		}

		targetInterfaces := activeInterfaces
		selected := appState.GetSelectedInterface()
		if selected != "" {
			found := false
			for _, iface := range activeInterfaces {
				if iface == selected {
					targetInterfaces = []string{selected}
					found = true
					break
				}
			}
			if !found {
				appState.AddLog(fmt.Sprintf("Selected interface %s not found, restoring DNS on all interfaces", selected))
			}
		}

		if len(targetInterfaces) == 0 {
			return fmt.Errorf("no active network interfaces available")
		}

		for _, iface := range targetInterfaces {
			// Reset DNS to automatic (DHCP)
			cmd := exec.Command("powershell", "Set-DnsClientServerAddress", "-InterfaceAlias", iface, "-ResetServerAddresses")
			output, err := cmd.CombinedOutput()
			if err != nil {
				errMsg := fmt.Sprintf("Error restoring DNS for interface %s: %v. Output: %s", iface, err, string(output))
				allErrors = append(allErrors, errMsg)
				appState.AddLog(errMsg)
			} else {
				appState.AddLog(fmt.Sprintf("Restored DNS to automatic (DHCP) for interface %s", iface))
			}
		}
	case "linux":
		// On Linux, we need to restore the original resolv.conf
		// This is complex as we'd need to backup the original, so for now we'll just log
		// The user may need to manually restore or restart network manager
		appState.AddLog("Note: On Linux, DNS restoration may require manual intervention or network manager restart")
		// Try to use systemd-resolved if available
		cmd := exec.Command("sh", "-c", "systemctl is-active --quiet systemd-resolved && sudo systemctl restart systemd-resolved || true")
		output, err := cmd.CombinedOutput()
		if err != nil {
			errMsg := fmt.Sprintf("Note: Could not automatically restore DNS on Linux: %v. Output: %s", err, string(output))
			appState.AddLog(errMsg)
		} else {
			appState.AddLog("Attempted to restore DNS via systemd-resolved")
		}
	case "darwin":
		// macOS: Reset DNS to automatic
		cmd := exec.Command("networksetup", "-setdnsservers", "Wi-Fi", "Empty")
		output, err := cmd.CombinedOutput()
		if err != nil {
			errMsg := fmt.Sprintf("Error restoring DNS on macOS: %v. Output: %s", err, string(output))
			allErrors = append(allErrors, errMsg)
			appState.AddLog(errMsg)
		} else {
			appState.AddLog("Restored DNS to automatic on macOS")
		}
	default:
		return fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}

	if len(allErrors) > 0 {
		errorMsg := strings.Join(allErrors, "\n")
		return fmt.Errorf("%s", errorMsg)
	}

	// Clear current DNS from state
	appState.SetCurrentDNS("", -1)

	// Update GUI if available
	if mainWindow != nil {
		fyne.Do(func() {
			updateStatusDisplay()
			updateLogsDisplay()
		})
	}

	return nil
}

// applyDNS applies the DNS settings to the system
func applyDNS(currentDNS string, currentIdx int) error {
	var allErrors []string

	switch runtime.GOOS {
	case "windows": // fuck you
		activeInterfaces, err := getActiveWindowsInterfaces()
		if err != nil {
			return err
		}
		// Keep app state in sync
		appState.SetInterfaces(activeInterfaces)

		targetInterfaces := activeInterfaces
		selected := appState.GetSelectedInterface()
		if selected != "" {
			found := false
			for _, iface := range activeInterfaces {
				if iface == selected {
					targetInterfaces = []string{selected}
					found = true
					break
				}
			}
			if !found {
				appState.AddLog(fmt.Sprintf("Selected interface %s not found, applying DNS to all interfaces", selected))
			}
		}

		if len(targetInterfaces) == 0 {
			return fmt.Errorf("no active network interfaces available")
		}

		for _, iface := range targetInterfaces {
			cmd := exec.Command("powershell", "Set-DnsClientServerAddress", "-InterfaceAlias", iface, "-ServerAddresses", currentDNS)
			output, err := cmd.CombinedOutput()
			if err != nil {
				errMsg := fmt.Sprintf("Error changing DNS for interface %s to %s: %v. Output: %s", iface, currentDNS, err, string(output))
				allErrors = append(allErrors, errMsg)
				appState.AddLog(errMsg)
			} else {
				if appState.GetDebugMode() {
					appState.AddLog(fmt.Sprintf("Changed DNS for interface %s to %s", iface, currentDNS))
				}
			}
		}
	case "linux": // THE GOAT
		cmd := exec.Command("sh", "-c", fmt.Sprintf("echo 'nameserver %s' | sudo tee /etc/resolv.conf", currentDNS))
		if appState.GetDebugMode() {
			appState.AddLog(fmt.Sprintf("Setting DNS on Linux to %s", currentDNS))
		}
		output, err := cmd.CombinedOutput()
		if err != nil {
			errMsg := fmt.Sprintf("Error setting DNS on Linux to %s: %v. Output: %s", currentDNS, err, string(output))
			allErrors = append(allErrors, errMsg)
			appState.AddLog(errMsg)
		}
	case "darwin": //shitos
		cmd := exec.Command("networksetup", "-setdnsservers", "Wi-Fi", currentDNS)
		output, err := cmd.CombinedOutput()
		if err != nil {
			errMsg := fmt.Sprintf("Error setting DNS on macOS to %s: %v. Output: %s", currentDNS, err, string(output))
			allErrors = append(allErrors, errMsg)
			appState.AddLog(errMsg)
		}
	default:
		return fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}

	if len(allErrors) > 0 {
		errorMsg := strings.Join(allErrors, "\n")
		return fmt.Errorf("%s", errorMsg)
	}

	if config.NotifyUser {
		err := beeep.Notify("DNS Change", fmt.Sprintf("DNS has been changed to %s", currentDNS), "")
		if err != nil {
			appState.AddLog(fmt.Sprintf("Warning: Failed to show notification: %v", err))
		}
	}

	// Update state with new DNS
	appState.SetCurrentDNS(currentDNS, currentIdx)

	// Calculate next change time
	var interval time.Duration
	if appState.GetDebugMode() {
		interval = 10 * time.Second
	} else {
		intervalMinutes := config.ChangeIntervalMinutes
		if intervalMinutes == 0 && config.ChangeIntervalHours > 0 {
			intervalMinutes = config.ChangeIntervalHours * 60
		}
		interval = time.Duration(intervalMinutes) * time.Minute
	}
	appState.SetNextChangeTime(time.Now().Add(interval))

	// Update GUI if available
	if mainWindow != nil {
		updateStatusDisplay()
		updateLogsDisplay()
	}

	return nil
}

func getActiveWindowsInterfaces() ([]string, error) {
	cmd := exec.Command("powershell", "-NoProfile", "-NonInteractive",
		"Get-NetAdapter | Where-Object { $_.Status -eq 'Up' } | Select-Object -ExpandProperty Name") // thanks ChatGPT, if I had to go through more powershell errors I would have gone insane.
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Error executing PowerShell command: %v, output: %s", err, output)
		return nil, fmt.Errorf("error executing PowerShell command: %v, output: %s", err, output)
	}
	text := strings.TrimSpace(string(output))
	if text == "" {
		return []string{}, nil
	}

	normalized := strings.ReplaceAll(text, "\r\n", "\n")
	lines := strings.Split(normalized, "\n")

	var interfaces []string
	for _, line := range lines {
		name := strings.TrimSpace(line)
		if name != "" {
			interfaces = append(interfaces, name)
		}
	}

	return interfaces, nil
}

func getIcon(s string) []byte {
	b, err := os.ReadFile(s)
	if err != nil {
		// Return empty if file doesn't exist, will use embedded
		return nil
	}
	return b
}

func generateDefaultConfig() error {
	configPath, err := getConfigPath()
	if err != nil {
		return fmt.Errorf("failed to get config path: %v", err)
	}

	defaultConfig := Config{
		DNSAddresses:          []string{"1.1.1.1", "1.0.0.1", "9.9.9.9"},
		RunOnStartup:          true,
		ChangeIntervalMinutes: 360, // 6 hours default
		NotifyUser:            true,
	}

	data, err := yaml.Marshal(&defaultConfig)
	if err != nil {
		return err
	}

	err = os.WriteFile(configPath, data, 0644)
	if err != nil {
		return err
	}

	return nil
}

func setRunOnStartup() error {
	exePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %v", err)
	}

	switch runtime.GOOS {
	case "windows":
		script := fmt.Sprintf(`
$path = 'HKCU:\Software\Microsoft\Windows\CurrentVersion\Run'
$name = 'DNSChanger'
$value = '%s'

if (Get-ItemProperty -Path $path -Name $name -ErrorAction SilentlyContinue) {
    Set-ItemProperty -Path $path -Name $name -Value $value
} else {
    New-ItemProperty -Path $path -Name $name -Value $value -PropertyType String
}
`, exePath) // why

		cmd := exec.Command("powershell", "-Command", script)
		output, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("failed to set run on startup in registry: %v. Output: %s", err, string(output))
		}
	case "linux", "darwin":
		// Create a .desktop file in ~/.config/autostart
		desktopFileContent := fmt.Sprintf(`[Desktop Entry]
Type=Application
Exec=%s
Hidden=false
NoDisplay=false
X-GNOME-Autostart-enabled=true
Name=DNSChanger
Comment=Start DNSChanger on startup
`, exePath)

		autostartDir := filepath.Join(os.Getenv("HOME"), ".config", "autostart")
		err := os.MkdirAll(autostartDir, 0755)
		if err != nil {
			return fmt.Errorf("failed to create autostart directory: %v", err)
		}

		desktopFilePath := filepath.Join(autostartDir, "DNSChanger.desktop")
		err = os.WriteFile(desktopFilePath, []byte(desktopFileContent), 0644)
		if err != nil {
			return fmt.Errorf("failed to write .desktop file: %v", err)
		}
	default:
		return fmt.Errorf("unsupported operating system: %s", runtime.GOOS) // can TempleOS even compile Go programs?
	}
	return nil
}
