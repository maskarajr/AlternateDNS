package main

import (
	"context"
	"fmt"
	"net"
	"strings"
	"time"
)

// DNSTestResult represents the result of testing a DNS server
type DNSTestResult struct {
	DNS          string
	AvgLatency   time.Duration
	SuccessRate  float64
	Status       string // "success", "error", "partial"
	Error        string
	TestCount    int
	SuccessCount int
}

// Default test domains for benchmarking
var defaultTestDomains = []string{
	"google.com",
	"cloudflare.com",
	"github.com",
	"microsoft.com",
	"amazon.com",
}

// testDNSLatency tests a DNS server by resolving multiple domains and returns average latency
func testDNSLatency(dnsServer string, testDomains []string, timeout time.Duration) DNSTestResult {
	result := DNSTestResult{
		DNS:          dnsServer,
		Status:       "success",
		TestCount:    len(testDomains),
		SuccessCount: 0,
	}

	if len(testDomains) == 0 {
		testDomains = defaultTestDomains
		result.TestCount = len(testDomains)
	}

	var latencies []time.Duration
	var errors []string

	for _, domain := range testDomains {
		start := time.Now()
		ctx, cancel := context.WithTimeout(context.Background(), timeout)

		// Create custom resolver for this DNS server
		resolver := &net.Resolver{
			PreferGo: true,
			Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
				d := net.Dialer{
					Timeout: timeout,
				}
				return d.DialContext(ctx, "udp", dnsServer+":53")
			},
		}

		// Test DNS resolution
		_, err := resolver.LookupIPAddr(ctx, domain)
		latency := time.Since(start)
		cancel()

		if err != nil {
			errors = append(errors, fmt.Sprintf("%s: %v", domain, err))
			result.Status = "partial"
		} else {
			latencies = append(latencies, latency)
			result.SuccessCount++
		}
	}

	// Calculate average latency
	if len(latencies) > 0 {
		var total time.Duration
		for _, lat := range latencies {
			total += lat
		}
		result.AvgLatency = total / time.Duration(len(latencies))
		result.SuccessRate = float64(result.SuccessCount) / float64(result.TestCount) * 100
	} else {
		result.Status = "error"
		result.SuccessRate = 0
		if len(errors) > 0 {
			result.Error = strings.Join(errors[:min(3, len(errors))], "; ")
		} else {
			result.Error = "Failed to resolve any domains"
		}
	}

	return result
}

// testDNSLatencyQuick tests a DNS server with a single domain for quick comparison
func testDNSLatencyQuick(dnsServer string, testDomain string, timeout time.Duration) (time.Duration, error) {
	if testDomain == "" {
		testDomain = "google.com"
	}

	start := time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	resolver := &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{
				Timeout: timeout,
			}
			return d.DialContext(ctx, "udp", dnsServer+":53")
		},
	}

	_, err := resolver.LookupIPAddr(ctx, testDomain)
	latency := time.Since(start)

	if err != nil {
		return 0, err
	}

	return latency, nil
}

// compareDNS compares two DNS servers and returns which is better
// Returns: (betterDNS, betterIndex, shouldSwitch)
func compareDNS(currentDNS string, currentIdx int, nextDNS string, nextIdx int, testDomains []string) (string, int, bool) {
	if len(testDomains) == 0 {
		testDomains = defaultTestDomains
	}

	timeout := 3 * time.Second

	appState.AddLog(fmt.Sprintf("Testing current DNS (%s) vs next DNS (%s)...", currentDNS, nextDNS))

	// Test current DNS
	currentResult := testDNSLatency(currentDNS, testDomains, timeout)
	appState.AddLog(fmt.Sprintf("Current DNS (%s): Avg latency %v, Success rate %.1f%%",
		currentDNS, currentResult.AvgLatency, currentResult.SuccessRate))

	// Test next DNS
	nextResult := testDNSLatency(nextDNS, testDomains, timeout)
	appState.AddLog(fmt.Sprintf("Next DNS (%s): Avg latency %v, Success rate %.1f%%",
		nextDNS, nextResult.AvgLatency, nextResult.SuccessRate))

	// Decision logic:
	// 1. If current has error status and next doesn't, switch
	// 2. If next has significantly better latency (>20% improvement), switch
	// 3. If current has low success rate (<50%) and next is better, switch
	// 4. Otherwise, keep current

	shouldSwitch := false

	if currentResult.Status == "error" && nextResult.Status != "error" {
		shouldSwitch = true
		appState.AddLog("Switching: Current DNS failed, next DNS works")
	} else if nextResult.Status == "error" {
		shouldSwitch = false
		appState.AddLog("Keeping current: Next DNS failed")
	} else if currentResult.SuccessRate < 50 && nextResult.SuccessRate > currentResult.SuccessRate {
		shouldSwitch = true
		appState.AddLog(fmt.Sprintf("Switching: Current success rate (%.1f%%) is low, next is better (%.1f%%)",
			currentResult.SuccessRate, nextResult.SuccessRate))
	} else if currentResult.AvgLatency > 0 && nextResult.AvgLatency > 0 {
		// Check if next is at least 20% faster
		improvement := float64(currentResult.AvgLatency-nextResult.AvgLatency) / float64(currentResult.AvgLatency) * 100
		if improvement > 20 {
			shouldSwitch = true
			appState.AddLog(fmt.Sprintf("Switching: Next DNS is %.1f%% faster", improvement))
		} else {
			appState.AddLog(fmt.Sprintf("Keeping current: Next DNS is only %.1f%% faster (need >20%%)", improvement))
		}
	} else {
		shouldSwitch = false
		appState.AddLog("Keeping current: Unable to determine better DNS")
	}

	if shouldSwitch {
		return nextDNS, nextIdx, true
	}
	return currentDNS, currentIdx, false
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
