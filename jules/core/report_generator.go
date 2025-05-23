// jules/core/report_generator.go
package core

import (
	"encoding/json"
	"fmt"
	"os"
	// "html/template" // For HTML reports
)

// Report represents the structure of the scan report.
type Report struct {
	TargetURL string  `json:"target_url"`
	Issues    []Issue `json:"issues"`
	// TODO: Add summary statistics, scan duration, etc.
}

// GenerateJSONReport generates a JSON report of the issues found.
func GenerateJSONReport(targetURL string, issues []Issue, outputPath string) error {
	report := Report{
		TargetURL: targetURL,
		Issues:    issues,
	}

	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal report to JSON: %w", err)
	}

	err = os.WriteFile(outputPath, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write JSON report to file: %w", err)
	}
	fmt.Printf("JSON report generated: %s\n", outputPath)
	return nil
}

// GenerateHTMLReport generates an HTML report of the issues found.
// func GenerateHTMLReport(targetURL string, issues []Issue, outputPath string) error {
// 	// TODO: Implement HTML report generation using html/template
// 	// - Define an HTML template
// 	// - Parse the template
// 	// - Execute the template with report data
// 	fmt.Printf("HTML report generation for %s not yet implemented. Output path: %s\n", targetURL, outputPath)
// 	return nil
// }
