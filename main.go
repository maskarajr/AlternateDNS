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
	"github.com/getlantern/systray"
	"gopkg.in/yaml.v2"
)

type Config struct {
	DNSAddresses        []string `yaml:"dns_addresses"`
	RunOnStartup        bool     `yaml:"run_on_startup"`
	ChangeIntervalHours int      `yaml:"change_interval_hours"`
	NotifyUser          bool     `yaml:"notify_user"`
	TestDomains         []string `yaml:"test_domains"` // Domains used for DNS latency testing
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
	setupGUI()

	// Initialize systray in background goroutine
	go func() {
		systray.Run(onReady, onExit)
	}()

	// Run Fyne event loop (this blocks and handles GUI events)
	guiApp.Run()
}

func onReady() {
	systray.SetIcon(appIcon)
	systray.SetTitle("AlternateDNS")
	systray.SetTooltip("AlternateDNS - DNS Rotation Tool")

	mOpenWindow := systray.AddMenuItem("Open Window", "Show the main window")
	mChange := systray.AddMenuItem("Change DNS", "Cycle to the next DNS server")
	mQuit := systray.AddMenuItem("Quit", "Quit the application")

	// Auto-start if configured
	if config.RunOnStartup {
		// Don't auto-start, let user control via GUI
		appState.AddLog("Application started (auto-start disabled, use GUI to start service)")
	}

	for {
		select {
		case <-mOpenWindow.ClickedCh:
			showWindow()
		case <-mQuit.ClickedCh:
			// Stop service if running
			if appState.IsRunning() {
				stopService()
			}
			systray.Quit()
			os.Exit(0)
		case <-mChange.ClickedCh:
			go func() {
				err := changeDNS(true) // Force change from systray menu
				if err != nil {
					beeep.Alert("DNS Change Error", err.Error(), "")
					appState.AddLog(fmt.Sprintf("ERROR: %v", err))
					updateLogsDisplay()
				} else {
					dns, idx := appState.GetCurrentDNS()
					appState.AddLog(fmt.Sprintf("DNS changed to %s (index %d)", dns, idx))
					updateLogsDisplay()
					updateStatusDisplay()
				}
			}()
		}
	}
}

// Tick is now handled by startTickerLoop in gui.go

func onExit() {
	// Stop service if running
	if appState.IsRunning() {
		stopService()
	}
	if appState.GetDebugMode() {
		appState.AddLog("Exiting the application")
	}
}

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

	if config.ChangeIntervalHours <= 0 {
		config.ChangeIntervalHours = 6
	}
	return nil
}

// changeDNS changes the DNS server. If forceChange is true, it switches to the next DNS
// without latency testing. If false, it uses smart switching logic to compare latencies.
func changeDNS(forceChange bool) error {
	if len(config.DNSAddresses) == 0 {
		return fmt.Errorf("no DNS addresses specified in config")
	}

	_, currentIdx := appState.GetCurrentDNS()

	// If no DNS is set yet (first time), just set the first one without testing
	if currentIdx < 0 || currentIdx >= len(config.DNSAddresses) {
		currentIdx = 0
		currentDNS := config.DNSAddresses[currentIdx]
		// Apply DNS immediately without testing on first run
		return applyDNS(currentDNS, currentIdx)
	}

	currentDNS := config.DNSAddresses[currentIdx]
	nextIndex := (currentIdx + 1) % len(config.DNSAddresses)
	nextDNS := config.DNSAddresses[nextIndex]

	// If force change, skip latency testing and switch immediately
	if forceChange {
		appState.AddLog(fmt.Sprintf("Force changing DNS from %s to %s", currentDNS, nextDNS))
		return applyDNS(nextDNS, nextIndex)
	}

	// Smart switching: Test current vs next DNS before switching (only for automatic changes)
	testDomains := config.TestDomains
	if len(testDomains) == 0 {
		testDomains = defaultTestDomains
	}

	betterDNS, betterIdx, shouldSwitch := compareDNS(currentDNS, currentIdx, nextDNS, nextIndex, testDomains)

	if !shouldSwitch {
		appState.AddLog(fmt.Sprintf("Keeping current DNS (%s) - it performs better than next DNS (%s)", currentDNS, nextDNS))
		// Still apply the current DNS to ensure it's set
		return applyDNS(betterDNS, betterIdx)
	}

	// Switch to better DNS
	appState.AddLog(fmt.Sprintf("Switching from %s to %s (better performance)", currentDNS, betterDNS))
	currentDNS = betterDNS
	currentIdx = betterIdx

	return applyDNS(currentDNS, currentIdx)
}

// applyDNS applies the DNS settings to the system
func applyDNS(currentDNS string, currentIdx int) error {
	var allErrors []string

	switch runtime.GOOS {
	case "windows": // fuck you
		interfaces, err := getActiveWindowsInterfaces()
		if err != nil {
			return err
		}
		for _, iface := range interfaces {
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
		interval = time.Duration(config.ChangeIntervalHours) * time.Hour
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
	interfaces := strings.Fields(string(output))
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
		DNSAddresses:        []string{"1.1.1.1", "1.0.0.1", "9.9.9.9"},
		RunOnStartup:        true,
		ChangeIntervalHours: 6,
		NotifyUser:          true,
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
