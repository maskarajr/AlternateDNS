package main

import (
	"sync"
	"time"
)

// AppState manages the application state thread-safely
type AppState struct {
	mu              sync.RWMutex
	isRunning       bool
	currentDNS      string
	currentDNSIndex int
	nextChangeTime  time.Time
	debugMode       bool
	ticker          *time.Ticker
	interfaces      []string
	logs            []string
	maxLogs         int
}

var appState = &AppState{
	maxLogs: 1000,
}

func (s *AppState) SetRunning(running bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.isRunning = running
}

func (s *AppState) IsRunning() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.isRunning
}

func (s *AppState) SetCurrentDNS(dns string, index int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.currentDNS = dns
	s.currentDNSIndex = index
}

func (s *AppState) GetCurrentDNS() (string, int) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.currentDNS, s.currentDNSIndex
}

func (s *AppState) SetNextChangeTime(t time.Time) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.nextChangeTime = t
}

func (s *AppState) GetNextChangeTime() time.Time {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.nextChangeTime
}

func (s *AppState) SetDebugMode(debug bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.debugMode = debug
}

func (s *AppState) GetDebugMode() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.debugMode
}

func (s *AppState) SetTicker(t *time.Ticker) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.ticker = t
}

func (s *AppState) GetTicker() *time.Ticker {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.ticker
}

func (s *AppState) SetInterfaces(ifaces []string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.interfaces = ifaces
}

func (s *AppState) GetInterfaces() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.interfaces
}

func (s *AppState) AddLog(message string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.logs = append(s.logs, message)
	if len(s.logs) > s.maxLogs {
		s.logs = s.logs[len(s.logs)-s.maxLogs:]
	}
}

func (s *AppState) GetLogs() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	logs := make([]string, len(s.logs))
	copy(logs, s.logs)
	return logs
}

func (s *AppState) ClearLogs() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.logs = []string{}
}

