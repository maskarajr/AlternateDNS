package main

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/gen2brain/beeep"
	"gopkg.in/yaml.v2"
)

var guiApp fyne.App
var mainWindow fyne.Window
var statusDNSLabel *widget.Label
var statusStatusLabel *widget.Label
var statusCountdownLabel *widget.Label
var statusInterfacesLabel *widget.Label
var statusInterfaceSelect *widget.Select
var statusStartStopBtn *widget.Button
var statusChangeNowBtn *widget.Button
var dnsList *widget.List
var dnsSelectedIndex int = -1
var dnsAddBtn *widget.Button
var dnsRemoveBtn *widget.Button
var dnsUpBtn *widget.Button
var dnsDownBtn *widget.Button
var settingsIntervalHoursEntry *widget.Entry
var settingsIntervalMinutesEntry *widget.Entry
var settingsStartupCheck *widget.Check
var settingsNotifyCheck *widget.Check
var settingsDebugCheck *widget.Check
var settingsSaveBtn *widget.Button
var logsText *widget.RichText
var logsClearBtn *widget.Button
var testerResultsList *widget.List
var testerTestBtn *widget.Button
var testerStatusLabel *widget.Label
var testerResults []DNSTestResult

var updateTimer *time.Ticker

// setupGUI creates and shows the main GUI window
func setupGUI() {
	guiApp = app.NewWithID("com.alternatedns.app")

	// Set icon
	iconResource := fyne.NewStaticResource("icon.ico", getEmbeddedIcon())
	guiApp.SetIcon(iconResource)

	mainWindow = guiApp.NewWindow(GetVersionString())
	mainWindow.Resize(fyne.NewSize(600, 500))
	mainWindow.CenterOnScreen()

	// Create tabs
	tabs := container.NewAppTabs(
		container.NewTabItem("Status", createStatusTab()),
		container.NewTabItem("DNS Servers", createDNSTab()),
		container.NewTabItem("DNS Tester", createDNSTesterTab()),
		container.NewTabItem("Settings", createSettingsTab()),
		container.NewTabItem("Logs", createLogsTab()),
	)

	mainWindow.SetContent(tabs)

	// Handle window close - minimize to tray instead
	mainWindow.SetCloseIntercept(func() {
		mainWindow.Hide()
	})

	// Start update timer
	startUpdateTimer()

	mainWindow.Show()
}

func createStatusTab() fyne.CanvasObject {
	// Current DNS display
	statusDNSLabel = widget.NewLabel("Not Set")
	statusDNSLabel.TextStyle = fyne.TextStyle{Bold: true}
	statusDNSLabel.Alignment = fyne.TextAlignCenter

	// Service status
	statusStatusLabel = widget.NewLabel("Stopped")
	statusStatusLabel.Alignment = fyne.TextAlignCenter

	// Countdown
	statusCountdownLabel = widget.NewLabel("--:--:--")
	statusCountdownLabel.Alignment = fyne.TextAlignCenter

	// Interfaces
	statusInterfacesLabel = widget.NewLabel("No interfaces detected")
	statusInterfacesLabel.Wrapping = fyne.TextWrapWord
	statusInterfaceSelect = widget.NewSelect([]string{}, func(value string) {
		if value == "" || value == appState.GetSelectedInterface() {
			return
		}
		appState.SetSelectedInterface(value)
		appState.AddLog(fmt.Sprintf("Active interface set to %s", value))
		updateLogsDisplay()
	})
	statusInterfaceSelect.PlaceHolder = "Select interface"
	statusInterfaceSelect.Disable()

	// Buttons
	statusStartStopBtn = widget.NewButton("Start Service", func() {
		if appState.IsRunning() {
			stopService()
		} else {
			startService()
		}
	})

	statusChangeNowBtn = widget.NewButton("Change DNS Now", func() {
		go func() {
			err := changeDNS(true) // Force change from manual button
			if err != nil {
				fyne.Do(func() {
					dialog.ShowError(err, mainWindow)
				})
				appState.AddLog(fmt.Sprintf("ERROR: %v", err))
				updateLogsDisplay()
			} else {
				dns, idx := appState.GetCurrentDNS()
				appState.AddLog(fmt.Sprintf("DNS changed to %s (index %d)", dns, idx))
				updateLogsDisplay()
				updateStatusDisplay()
			}
		}()
	})

	// Layout
	statusContainer := container.NewVBox(
		widget.NewCard("Current DNS", "", statusDNSLabel),
		widget.NewCard("Service Status", "", statusStatusLabel),
		widget.NewCard("Next Change In", "", statusCountdownLabel),
		widget.NewCard("Active Interfaces", "", container.NewVBox(
			statusInterfacesLabel,
			statusInterfaceSelect,
		)),
		container.NewHBox(statusStartStopBtn, statusChangeNowBtn),
	)

	return container.NewScroll(statusContainer)
}

func createDNSTab() fyne.CanvasObject {
	// DNS list
	dnsList = widget.NewList(
		func() int {
			return len(config.DNSAddresses)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			label := obj.(*widget.Label)
			if id < len(config.DNSAddresses) {
				currentDNS, currentIdx := appState.GetCurrentDNS()
				dns := config.DNSAddresses[id]
				marker := ""
				if id == currentIdx && appState.IsRunning() {
					marker = " → "
				}
				label.SetText(fmt.Sprintf("%d. %s%s", id+1, dns, marker))
				if dns == currentDNS {
					label.Importance = widget.HighImportance
				} else {
					label.Importance = widget.MediumImportance
				}
			}
		},
	)

	// Track selection
	dnsList.OnSelected = func(id widget.ListItemID) {
		dnsSelectedIndex = int(id)
	}

	// Buttons
	dnsAddBtn = widget.NewButton("Add DNS", func() {
		showAddDNSDialog()
	})

	dnsRemoveBtn = widget.NewButton("Remove Selected", func() {
		if dnsSelectedIndex >= 0 && dnsSelectedIndex < len(config.DNSAddresses) {
			config.DNSAddresses = append(config.DNSAddresses[:dnsSelectedIndex], config.DNSAddresses[dnsSelectedIndex+1:]...)
			saveConfig()
			dnsSelectedIndex = -1
			dnsList.Refresh()
			updateStatusDisplay()
		}
	})

	dnsUpBtn = widget.NewButton("Move Up", func() {
		if dnsSelectedIndex > 0 && dnsSelectedIndex < len(config.DNSAddresses) {
			config.DNSAddresses[dnsSelectedIndex], config.DNSAddresses[dnsSelectedIndex-1] = config.DNSAddresses[dnsSelectedIndex-1], config.DNSAddresses[dnsSelectedIndex]
			saveConfig()
			dnsSelectedIndex--
			dnsList.Select(dnsSelectedIndex)
			dnsList.Refresh()
		}
	})

	dnsDownBtn = widget.NewButton("Move Down", func() {
		if dnsSelectedIndex >= 0 && dnsSelectedIndex < len(config.DNSAddresses)-1 {
			config.DNSAddresses[dnsSelectedIndex], config.DNSAddresses[dnsSelectedIndex+1] = config.DNSAddresses[dnsSelectedIndex+1], config.DNSAddresses[dnsSelectedIndex]
			saveConfig()
			dnsSelectedIndex++
			dnsList.Select(dnsSelectedIndex)
			dnsList.Refresh()
		}
	})

	buttonContainer := container.NewGridWithColumns(2,
		dnsAddBtn,
		dnsRemoveBtn,
		dnsUpBtn,
		dnsDownBtn,
	)

	return container.NewBorder(nil, buttonContainer, nil, nil, dnsList)
}

func createSettingsTab() fyne.CanvasObject {
	// Interval - Hours and Minutes
	totalMinutes := config.ChangeIntervalMinutes
	if totalMinutes == 0 && config.ChangeIntervalHours > 0 {
		totalMinutes = config.ChangeIntervalHours * 60
	}
	hours := totalMinutes / 60
	minutes := totalMinutes % 60

	settingsIntervalHoursEntry = widget.NewEntry()
	settingsIntervalHoursEntry.SetText(fmt.Sprintf("%d", hours))
	settingsIntervalHoursEntry.SetPlaceHolder("Hours")

	settingsIntervalMinutesEntry = widget.NewEntry()
	settingsIntervalMinutesEntry.SetText(fmt.Sprintf("%d", minutes))
	settingsIntervalMinutesEntry.SetPlaceHolder("Minutes")

	// Checkboxes
	settingsStartupCheck = widget.NewCheck("Run on startup", nil)
	settingsStartupCheck.SetChecked(config.RunOnStartup)

	settingsNotifyCheck = widget.NewCheck("Notify on DNS change", nil)
	settingsNotifyCheck.SetChecked(config.NotifyUser)

	settingsDebugCheck = widget.NewCheck("Debug mode", nil)
	settingsDebugCheck.SetChecked(appState.GetDebugMode())
	settingsDebugCheck.OnChanged = func(checked bool) {
		appState.SetDebugMode(checked)
		appState.AddLog(fmt.Sprintf("Debug mode: %v", checked))
		updateLogsDisplay()
	}

	// Save button
	settingsSaveBtn = widget.NewButton("Save Settings", func() {
		saveSettings()
	})

	// Version info
	versionLabel := widget.NewRichText()
	versionLabel.ParseMarkdown(fmt.Sprintf("**Version Information**\n\n%s", GetFullVersionInfo()))
	versionLabel.Wrapping = fyne.TextWrapWord

	intervalContainer := container.NewGridWithColumns(2,
		container.NewVBox(
			widget.NewLabel("Hours:"),
			settingsIntervalHoursEntry,
		),
		container.NewVBox(
			widget.NewLabel("Minutes:"),
			settingsIntervalMinutesEntry,
		),
	)

	settingsContainer := container.NewVBox(
		widget.NewForm(
			widget.NewFormItem("Change Interval", intervalContainer),
		),
		settingsStartupCheck,
		settingsNotifyCheck,
		settingsDebugCheck,
		settingsSaveBtn,
		widget.NewSeparator(),
		versionLabel,
	)

	return container.NewScroll(settingsContainer)
}

func createLogsTab() fyne.CanvasObject {
	logsText = widget.NewRichText()
	logsText.Wrapping = fyne.TextWrapWord
	logsText.Scroll = container.ScrollBoth

	logsClearBtn = widget.NewButton("Clear Logs", func() {
		appState.ClearLogs()
		updateLogsDisplay()
	})

	return container.NewBorder(nil, logsClearBtn, nil, nil, container.NewScroll(logsText))
}

func showAddDNSDialog() {
	entry := widget.NewEntry()
	entry.SetPlaceHolder("Enter DNS address (e.g., 1.1.1.1)")

	dialog.ShowForm("Add DNS Server", "Add", "Cancel",
		[]*widget.FormItem{
			widget.NewFormItem("DNS Address", entry),
		},
		func(confirmed bool) {
			if confirmed {
				dns := strings.TrimSpace(entry.Text)
				if dns != "" {
					config.DNSAddresses = append(config.DNSAddresses, dns)
					saveConfig()
					dnsList.Refresh()
					appState.AddLog(fmt.Sprintf("Added DNS: %s", dns))
					updateLogsDisplay()
				}
			}
		},
		mainWindow,
	)
}

func saveSettings() {
	// Parse hours and minutes
	var hours, minutes int
	_, err := fmt.Sscanf(settingsIntervalHoursEntry.Text, "%d", &hours)
	if err != nil {
		dialog.ShowError(fmt.Errorf("invalid hours: must be a number"), mainWindow)
		return
	}

	_, err = fmt.Sscanf(settingsIntervalMinutesEntry.Text, "%d", &minutes)
	if err != nil {
		dialog.ShowError(fmt.Errorf("invalid minutes: must be a number"), mainWindow)
		return
	}

	if hours < 0 || minutes < 0 {
		dialog.ShowError(fmt.Errorf("hours and minutes must be non-negative"), mainWindow)
		return
	}

	if hours == 0 && minutes == 0 {
		dialog.ShowError(fmt.Errorf("interval must be greater than 0"), mainWindow)
		return
	}

	totalMinutes := hours*60 + minutes
	config.ChangeIntervalMinutes = totalMinutes
	config.ChangeIntervalHours = 0 // Clear old value
	config.RunOnStartup = settingsStartupCheck.Checked
	config.NotifyUser = settingsNotifyCheck.Checked

	saveConfig()

	// Update startup if needed
	if config.RunOnStartup {
		err = setRunOnStartup()
		if err != nil {
			appState.AddLog(fmt.Sprintf("Warning: Failed to set startup: %v", err))
		}
	}

	// Restart ticker if running
	if appState.IsRunning() {
		ticker := appState.GetTicker()
		if ticker != nil {
			ticker.Stop()
		}
		var newTicker *time.Ticker
		if appState.GetDebugMode() {
			newTicker = time.NewTicker(10 * time.Second)
		} else {
			newTicker = time.NewTicker(time.Duration(config.ChangeIntervalMinutes) * time.Minute)
		}
		appState.SetTicker(newTicker)
		go startTickerLoop(newTicker)
	}

	dialog.ShowInformation("Settings Saved", "Settings have been saved successfully.", mainWindow)
	appState.AddLog("Settings saved")
	updateLogsDisplay()
}

func saveConfig() {
	configPath, err := getConfigPath()
	if err != nil {
		appState.AddLog(fmt.Sprintf("ERROR: Failed to get config path: %v", err))
		updateLogsDisplay()
		return
	}

	data, err := yaml.Marshal(&config)
	if err != nil {
		appState.AddLog(fmt.Sprintf("ERROR: Failed to marshal config: %v", err))
		updateLogsDisplay()
		return
	}

	err = os.WriteFile(configPath, data, 0644)
	if err != nil {
		appState.AddLog(fmt.Sprintf("ERROR: Failed to write config: %v", err))
		updateLogsDisplay()
		return
	}
}

func startService() {
	if appState.IsRunning() {
		return
	}

	// Check admin
	err := checkAdmin()
	if err != nil {
		dialog.ShowError(err, mainWindow)
		appState.AddLog(fmt.Sprintf("ERROR: %v", err))
		updateLogsDisplay()
		return
	}

	appState.SetRunning(true)
	var ticker *time.Ticker
	if appState.GetDebugMode() {
		ticker = time.NewTicker(10 * time.Second)
	} else {
		intervalMinutes := config.ChangeIntervalMinutes
		if intervalMinutes == 0 && config.ChangeIntervalHours > 0 {
			intervalMinutes = config.ChangeIntervalHours * 60
		}
		ticker = time.NewTicker(time.Duration(intervalMinutes) * time.Minute)
	}
	appState.SetTicker(ticker)

	// Initial DNS change (use smart switching on service start)
	go func() {
		err := changeDNS(false)
		if err != nil {
			appState.AddLog(fmt.Sprintf("ERROR: %v", err))
			updateLogsDisplay()
			dialog.ShowError(err, mainWindow)
		} else {
			appState.AddLog("Service started")
			updateLogsDisplay()
		}
		updateStatusDisplay()
	}()

	// Start ticker loop
	go startTickerLoop(ticker)

	updateStatusDisplay()
}

func stopService() {
	if !appState.IsRunning() {
		return
	}

	appState.SetRunning(false)
	ticker := appState.GetTicker()
	if ticker != nil {
		ticker.Stop()
		appState.SetTicker(nil)
	}

	// Restore DNS to automatic/DHCP
	go func() {
		err := restoreDNS()
		if err != nil {
			appState.AddLog(fmt.Sprintf("ERROR: Failed to restore DNS: %v", err))
			fyne.Do(func() {
				dialog.ShowError(err, mainWindow)
			})
		} else {
			appState.AddLog("DNS restored to automatic (DHCP)")
		}
		fyne.Do(func() {
			updateLogsDisplay()
			updateStatusDisplay()
		})
	}()

	appState.AddLog("Service stopped")
	fyne.Do(func() {
		updateLogsDisplay()
		updateStatusDisplay()
	})
}

func startTickerLoop(ticker *time.Ticker) {
	for range ticker.C {
		if !appState.IsRunning() {
			break
		}
		err := changeDNS(false) // Use smart switching for automatic timer-based changes
		if err != nil {
			appState.AddLog(fmt.Sprintf("ERROR: %v", err))
			updateLogsDisplay()
			if config.NotifyUser {
				beeep.Alert("DNS Change Error", err.Error(), "")
			}
		} else {
			dns, _ := appState.GetCurrentDNS()
			appState.AddLog(fmt.Sprintf("DNS changed to %s", dns))
			updateLogsDisplay()
		}
		updateStatusDisplay()
	}
}

func updateStatusDisplay() {
	if mainWindow == nil {
		return
	}

	// All GUI updates must be on main thread
	fyne.Do(func() {
		// Update DNS
		currentDNS, currentIdx := appState.GetCurrentDNS()
		if currentDNS != "" {
			statusDNSLabel.SetText(fmt.Sprintf("%s\n(Index: %d)", currentDNS, currentIdx+1))
		} else {
			statusDNSLabel.SetText("Not Set")
		}

		// Update status
		if appState.IsRunning() {
			statusStatusLabel.SetText("Running")
			statusStatusLabel.Importance = widget.SuccessImportance
			statusStartStopBtn.SetText("Stop Service")
		} else {
			statusStatusLabel.SetText("Stopped")
			statusStatusLabel.Importance = widget.WarningImportance
			statusStartStopBtn.SetText("Start Service")
		}

		// Update countdown
		if appState.IsRunning() {
			nextChange := appState.GetNextChangeTime()
			if !nextChange.IsZero() {
				remaining := time.Until(nextChange)
				if remaining > 0 {
					hours := int(remaining.Hours())
					minutes := int(remaining.Minutes()) % 60
					seconds := int(remaining.Seconds()) % 60
					statusCountdownLabel.SetText(fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds))
				} else {
					statusCountdownLabel.SetText("Changing soon...")
				}
			} else {
				statusCountdownLabel.SetText("Calculating...")
			}
		} else {
			statusCountdownLabel.SetText("--:--:--")
		}

		// Update interfaces
		interfaces := appState.GetInterfaces()
		if statusInterfaceSelect != nil {
			selectedInterface := appState.GetSelectedInterface()
			switch {
			case len(interfaces) == 0:
				statusInterfacesLabel.SetText("No interfaces detected")
				statusInterfaceSelect.Options = []string{}
				statusInterfaceSelect.ClearSelected()
				statusInterfaceSelect.PlaceHolder = "No interfaces"
				statusInterfaceSelect.Disable()
			case len(interfaces) == 1:
				single := interfaces[0]
				statusInterfacesLabel.SetText(single)
				statusInterfaceSelect.Options = interfaces
				if statusInterfaceSelect.Selected != single {
					statusInterfaceSelect.SetSelected(single)
				}
				appState.SetSelectedInterface(single)
				statusInterfaceSelect.Disable()
			default:
				statusInterfacesLabel.SetText(strings.Join(interfaces, ", "))
				statusInterfaceSelect.Options = interfaces
				if selectedInterface == "" {
					selectedInterface = interfaces[0]
					appState.SetSelectedInterface(selectedInterface)
				}
				if statusInterfaceSelect.Selected != selectedInterface {
					statusInterfaceSelect.SetSelected(selectedInterface)
				}
				if statusInterfaceSelect.Disabled() {
					statusInterfaceSelect.Enable()
				}
			}
			statusInterfaceSelect.Refresh()
		} else {
			if len(interfaces) > 0 {
				statusInterfacesLabel.SetText(strings.Join(interfaces, ", "))
			} else {
				statusInterfacesLabel.SetText("No interfaces detected")
			}
		}

		// Refresh DNS list to show current marker
		if dnsList != nil {
			dnsList.Refresh()
		}
	})
}

func updateLogsDisplay() {
	if logsText == nil {
		return
	}

	// All GUI updates must be on main thread
	fyne.Do(func() {
		logs := appState.GetLogs()
		logText := strings.Join(logs, "\n")

		// Create rich text segments - use foreground color from theme (white on dark theme)
		if logText != "" {
			segments := []widget.RichTextSegment{
				&widget.TextSegment{
					Text: logText,
					Style: widget.RichTextStyle{
						ColorName: theme.ColorNameForeground,
					},
				},
			}
			logsText.Segments = segments
		} else {
			logsText.Segments = []widget.RichTextSegment{}
		}
		logsText.Refresh()
	})
}

func createDNSTesterTab() fyne.CanvasObject {
	// Status label
	testerStatusLabel = widget.NewLabel("Ready to test DNS servers")
	testerStatusLabel.Wrapping = fyne.TextWrapWord

	// Test button
	testerTestBtn = widget.NewButton("Test All DNS Servers", func() {
		go runDNSTests()
	})

	// Results list
	testerResults = []DNSTestResult{}
	testerResultsList = widget.NewList(
		func() int {
			return len(testerResults)
		},
		func() fyne.CanvasObject {
			dnsLabel := widget.NewLabel("")
			latencyLabel := widget.NewLabel("")
			successLabel := widget.NewLabel("")
			statusLabel := widget.NewLabel("")
			return container.NewHBox(dnsLabel, latencyLabel, successLabel, statusLabel)
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			if id >= len(testerResults) {
				return
			}
			result := testerResults[id]

			// Get the HBox container and its children
			box := obj.(*fyne.Container)
			labels := box.Objects

			// DNS server
			labels[0].(*widget.Label).SetText(result.DNS)

			// Latency
			if result.AvgLatency > 0 {
				labels[1].(*widget.Label).SetText(result.AvgLatency.Round(time.Millisecond).String())
			} else {
				labels[1].(*widget.Label).SetText("N/A")
			}

			// Success rate
			labels[2].(*widget.Label).SetText(fmt.Sprintf("%.1f%%", result.SuccessRate))

			// Status
			statusText := result.Status
			if result.Error != "" {
				statusText += " (" + result.Error + ")"
			}
			labels[3].(*widget.Label).SetText(statusText)
		},
	)

	// Layout
	top := container.NewVBox(
		testerStatusLabel,
		testerTestBtn,
	)

	return container.NewBorder(
		top,
		nil,
		nil,
		nil,
		testerResultsList,
	)
}

func runDNSTests() {
	fyne.Do(func() {
		testerStatusLabel.SetText("Testing DNS servers...")
		testerTestBtn.Disable()
		testerResults = []DNSTestResult{}
		testerResultsList.Refresh()
	})

	testDomains := config.TestDomains
	if len(testDomains) == 0 {
		testDomains = defaultTestDomains
	}

	results := make([]DNSTestResult, len(config.DNSAddresses))

	for i, dns := range config.DNSAddresses {
		appState.AddLog(fmt.Sprintf("Testing DNS server: %s", dns))
		result := testDNSLatency(dns, testDomains, 5*time.Second)
		results[i] = result

		appState.AddLog(fmt.Sprintf("DNS %s: Avg latency %v, Success rate %.1f%%, Status: %s",
			dns, result.AvgLatency, result.SuccessRate, result.Status))

		// Update GUI
		fyne.Do(func() {
			testerResults = results[:i+1]
			testerResultsList.Refresh()
		})
	}

	// Sort by latency (best first)
	sort.Slice(results, func(i, j int) bool {
		if results[i].Status == "error" {
			return false
		}
		if results[j].Status == "error" {
			return true
		}
		return results[i].AvgLatency < results[j].AvgLatency
	})

	fyne.Do(func() {
		testerResults = results
		testerResultsList.Refresh()
		testerStatusLabel.SetText(fmt.Sprintf("Testing complete. Tested %d DNS servers.", len(results)))
		testerTestBtn.Enable()
	})
}

func startUpdateTimer() {
	updateTimer = time.NewTicker(1 * time.Second)
	go func() {
		for range updateTimer.C {
			if mainWindow != nil {
				updateStatusDisplay()
			}
		}
	}()
}

func showWindow() {
	if mainWindow != nil {
		fyne.Do(func() {
			mainWindow.Show()
			mainWindow.RequestFocus()
		})
	}
}
