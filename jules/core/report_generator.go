// jules/core/report_generator.go
package core

import (
	"encoding/json"
	"fmt"
	"html/template" // For HTML reports (still placeholder)
	"os"
	"time" // For report timestamp
)

// Report represents the structure of the scan report.
// It now includes summary statistics about LLM and KB usage.
type Report struct {
	ScanID                   string    `json:"scan_id"`
	TargetURL                string    `json:"target_url"`
	ScanStartTime            time.Time `json:"scan_start_time"`
	ScanEndTime              time.Time `json:"scan_end_time"`
	TotalIssuesFound         int       `json:"total_issues_found"`
	TotalLLMAssistedFindings int       `json:"total_llm_assisted_findings"`
	TotalKBRefinedFindings   int       `json:"total_kb_refined_findings"`
	Issues                   []Issue   `json:"issues"` // Uses the updated Issue struct from plugin_system.go
	// TODO: Add summary of URLs scanned, forms found, inputs analyzed etc.
}

// GenerateJSONReport generates a JSON report of the issues found.
func GenerateJSONReport(scanID string, targetURL string, startTime time.Time, issues []Issue, outputPath string) error {
	llmAssistedCount := 0
	kbRefinedCount := 0
	for _, issue := range issues {
		if issue.LLMPayloadSource != "" || issue.LLMPrompt != "" { // Basic check for LLM assistance
			llmAssistedCount++
		}
		// A finding is KB-refined if its payload source explicitly states it OR if it has contributing KB articles.
		if issue.LLMPayloadSource == "refined_from_kb" || len(issue.ContributingKBArticles) > 0 {
			kbRefinedCount++
		}
	}

	report := Report{
		ScanID:                   scanID,
		TargetURL:                targetURL,
		ScanStartTime:            startTime,
		ScanEndTime:              time.Now(), // Set scan end time
		TotalIssuesFound:         len(issues),
		TotalLLMAssistedFindings: llmAssistedCount,
		TotalKBRefinedFindings:   kbRefinedCount,
		Issues:                   issues,
	}

	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal report to JSON: %w", err)
	}

	err = os.WriteFile(outputPath, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write JSON report to file: %w", err)
	}
	fmt.Printf("REPORT_GEN: JSON report generated: %s\n", outputPath)
	return nil
}

// GenerateHTMLReport generates an HTML report of the issues found. (Placeholder)
func GenerateHTMLReport(scanID string, targetURL string, startTime time.Time, issues []Issue, outputPath string) error {
	// TODO: Implement HTML report generation using html/template
	// - Define an HTML template that can display all fields in the Report and Issue structs.
	// - Include sections for LLM prompts, responses, confidence, KB articles.
	// - Parse the template.
	// - Execute the template with report data.
	
	// For now, just create a very basic placeholder file or log.
	// Calculate summary stats for HTML placeholder too
	llmAssistedCount := 0
	kbRefinedCount := 0
	for _, issue := range issues {
		if issue.LLMPayloadSource != "" || issue.LLMPrompt != "" {
			llmAssistedCount++
		}
		if issue.LLMPayloadSource == "refined_from_kb" || len(issue.ContributingKBArticles) > 0 {
			kbRefinedCount++
		}
	}
	
	htmlContent := `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Jules Security Scan Report (Placeholder)</title>
    <style>body { font-family: sans-serif; margin: 20px; } h1, h2 { color: #333; } li { margin-bottom: 5px; }</style>
</head>
<body>
    <h1>Jules Security Scan Report (Placeholder)</h1>
    <p><strong>Scan ID:</strong> ` + template.HTMLEscapeString(scanID) + `</p>
    <p><strong>Target URL:</strong> ` + template.HTMLEscapeString(targetURL) + `</p>
    <p><strong>Scan Start Time:</strong> ` + startTime.Format(time.RFC1123) + `</p>
    <p><strong>Scan End Time:</strong> ` + time.Now().Format(time.RFC1123) + `</p>
    <p><strong>Total Issues Found:</strong> ` + fmt.Sprintf("%d", len(issues)) + `</p>
    <p><strong>Total LLM-Assisted Findings:</strong> ` + fmt.Sprintf("%d", llmAssistedCount) + `</p>
    <p><strong>Total KB-Refined Findings:</strong> ` + fmt.Sprintf("%d", kbRefinedCount) + `</p>
    <p><em>This is a placeholder HTML report. Full details for each issue, including LLM prompts, LLM responses (snippets), confidence scores, payload sources, and contributing KB articles, will be displayed here in a structured format.</em></p>
    <h2>Issues Summary (Illustrative)</h2>
    <ul>`
	for _, issue := range issues {
		details := fmt.Sprintf("Severity: %s, Type: %s, URL: %s, Input: %s", issue.Severity, issue.Type, issue.URL, issue.InputName)
		if issue.LLMPayloadSource != "" {
			details += fmt.Sprintf(" (LLM Source: %s, Confidence: %.2f)", issue.LLMPayloadSource, issue.LLMConfidence)
		}
		if len(issue.ContributingKBArticles) > 0 {
			details += fmt.Sprintf(" (KB Articles: %s)", template.HTMLEscapeString(fmt.Sprintf("%v", issue.ContributingKBArticles)))
		}
		htmlContent += `<li>` + template.HTMLEscapeString(details) + `</li>`
	}
	htmlContent += `
    </ul>
</body>
</html>`

	err := os.WriteFile(outputPath, []byte(htmlContent), 0644)
	if err != nil {
		return fmt.Errorf("failed to write placeholder HTML report: %w", err)
	}
	fmt.Printf("REPORT_GEN: Placeholder HTML report generated: %s\n", outputPath)
	return nil
}
