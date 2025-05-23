// jules/core/plugin_system.go
package core

import "fmt"

// ScanPlugin defines the interface for a scanner plugin.
type ScanPlugin interface {
	Name() string
	Description() string
	RunScan(targetURL string, client *JulesHTTPClient) ([]Issue, error)
}

// Issue represents a vulnerability found by a plugin.
type Issue struct {
	Type        string
	Severity    string
	URL         string
	Description string
	// TODO: Add more details like CWE, evidence, remediation advice
}

type PluginManager struct {
	Plugins []ScanPlugin
}

func NewPluginManager() *PluginManager {
	return &PluginManager{
		Plugins: []ScanPlugin{},
	}
}

func (pm *PluginManager) RegisterPlugin(plugin ScanPlugin) {
	pm.Plugins = append(pm.Plugins, plugin)
	fmt.Printf("Registered plugin: %s\n", plugin.Name())
}

func (pm *PluginManager) RunScans(targetURL string, client *JulesHTTPClient) []Issue {
	var allIssues []Issue
	for _, plugin := range pm.Plugins {
		fmt.Printf("Running scan with plugin: %s on %s\n", plugin.Name(), targetURL)
		issues, err := plugin.RunScan(targetURL, client)
		if err != nil {
			fmt.Printf("Error running plugin %s: %v\n", plugin.Name(), err)
			continue
		}
		allIssues = append(allIssues, issues...)
	}
	return allIssues
}
