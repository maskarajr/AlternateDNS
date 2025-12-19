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

// findBestDNS tests all DNS servers and returns the one with the lowest latency
// Returns: (bestDNS, bestIndex)
func findBestDNS(dnsServers []string, testDomains []string) (string, int) {
	if len(dnsServers) == 0 {
		return "", -1
	}

	if len(testDomains) == 0 {
		testDomains = defaultTestDomains
	}

	timeout := 3 * time.Second
	appState.AddLog(fmt.Sprintf("Testing all %d DNS servers to find the best one...", len(dnsServers)))

	var bestDNS string
	bestIdx := -1
	bestLatency := time.Duration(0)
	bestSuccessRate := 0.0

	for idx, dns := range dnsServers {
		result := testDNSLatency(dns, testDomains, timeout)

		appState.AddLog(fmt.Sprintf("DNS %d/%d (%s): Avg latency %v, Success rate %.1f%%",
			idx+1, len(dnsServers), dns, result.AvgLatency, result.SuccessRate))

		// Skip DNS servers that completely failed
		if result.Status == "error" {
			appState.AddLog(fmt.Sprintf("  Skipping %s: Failed to resolve any domains", dns))
			continue
		}

		// If this is the first working DNS, use it as baseline
		if bestIdx == -1 {
			bestDNS = dns
			bestIdx = idx
			bestLatency = result.AvgLatency
			bestSuccessRate = result.SuccessRate
			continue
		}

		// Decision logic:
		// 1. Prefer DNS with higher success rate (if difference is significant >10%)
		// 2. If success rates are similar, prefer lower latency
		// 3. If current has very low success rate (<50%) and candidate is better, switch

		shouldUseThis := false

		if result.SuccessRate < 50 && bestSuccessRate >= 50 {
			// Current best is good, candidate is bad - keep best
			shouldUseThis = false
		} else if bestSuccessRate < 50 && result.SuccessRate >= 50 {
			// Current best is bad, candidate is good - switch
			shouldUseThis = true
		} else if result.SuccessRate > bestSuccessRate+10 {
			// Candidate has significantly better success rate (>10% difference)
			shouldUseThis = true
		} else if result.SuccessRate >= bestSuccessRate-10 && result.AvgLatency > 0 {
			// Success rates are similar (within 10%), compare latency
			if result.AvgLatency < bestLatency {
				shouldUseThis = true
			}
		}

		if shouldUseThis {
			bestDNS = dns
			bestIdx = idx
			bestLatency = result.AvgLatency
			bestSuccessRate = result.SuccessRate
		}
	}

	if bestIdx == -1 {
		// All DNS servers failed, fallback to first one
		appState.AddLog("Warning: All DNS servers failed testing, using first DNS as fallback")
		return dnsServers[0], 0
	}

	appState.AddLog(fmt.Sprintf("Best DNS selected: %s (latency: %v, success rate: %.1f%%)",
		bestDNS, bestLatency, bestSuccessRate))

	return bestDNS, bestIdx
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
