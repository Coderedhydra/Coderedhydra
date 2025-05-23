// jules/ui/src/components/ScanResults.jsx
import React, { useState }
from 'react';

function ScanResults() {
    // TODO: Fetch and display scan results from the core/backend
    // Mock data now includes LLM and KB related fields
    const initialIssues = [
        {
            id: 1,
            severity: 'High',
            type: 'SQLi (LLM Generated)',
            url: 'https://example.com/search?q=1',
            inputName: 'q',
            details: 'SQL Injection vulnerability found in parameter q.',
            payloadUsed: "' OR 1=1 --",
            evidence: "SQL syntax error near '1=1'",
            llmPrompt: "Generate SQLi for input 'q' in URL ...",
            llmRawResponse: "[' OR 1=1 --, ...]", // Snippet
            llmConfidence: 0.92,
            llmPayloadSource: "llm_generated",
            contributingKBArticles: ["CVE-2022-XXXX"],
        },
        {
            id: 2,
            severity: 'Medium',
            type: 'XSS (LLM Refined from KB)',
            url: 'https://example.com/profile',
            inputName: 'displayName',
            details: 'Reflected XSS in displayName field.',
            payloadUsed: "<img src=x onerror=alert('jules_xss')>",
            evidence: "Payload reflected in response: ...<img src=x onerror=alert...",
            llmPrompt: "Refine XSS for 'displayName' using KB-001 insights...",
            llmRawResponse: "[<img src=x onerror=alert('jules_xss')>]",
            llmConfidence: 0.85,
            llmPayloadSource: "llm_refined_from_kb",
            contributingKBArticles: ["OWASP-XSS", "KB-001-BypassTechnique"],
        },
        {
            id: 3,
            severity: 'Low',
            type: 'Missing Security Header',
            url: 'https://example.com',
            inputName: 'N/A',
            details: 'Content-Security-Policy header not set.',
            payloadUsed: "N/A",
            evidence: "Headers: ...",
            // No LLM fields for non-LLM specific findings, or they can be null/empty
        },
    ];

    const [issues, setIssues] = useState(initialIssues);
    const [severityFilter, setSeverityFilter] = useState('all');
    const [sourceFilter, setSourceFilter] = useState('all'); // 'all', 'llm_generated', 'llm_refined_from_kb', 'manual'

    const filteredIssues = issues.filter(issue => {
        const severityMatch = severityFilter === 'all' || issue.severity.toLowerCase() === severityFilter;
        const sourceMatch = sourceFilter === 'all' || (issue.llmPayloadSource && issue.llmPayloadSource === sourceFilter) || (sourceFilter === 'manual' && !issue.llmPayloadSource);
        return severityMatch && sourceMatch;
    });

    return (
        <div className="scan-results-panel" style={{ border: '1px solid #ccc', padding: '10px', margin: '10px' }}>
            <h2>Scan Results</h2>
            <div>
                <strong>Filter by Severity:</strong>
                <select onChange={(e) => setSeverityFilter(e.target.value)} value={severityFilter} style={{margin: '5px'}}>
                    <option value="all">All</option>
                    <option value="critical">Critical</option>
                    <option value="high">High</option>
                    <option value="medium">Medium</option>
                    <option value="low">Low</option>
                    <option value="info">Info</option>
                </select>
                <strong>Filter by Source:</strong>
                <select onChange={(e) => setSourceFilter(e.target.value)} value={sourceFilter} style={{margin: '5px'}}>
                    <option value="all">All Sources</option>
                    <option value="llm_generated">LLM Generated</option>
                    <option value="llm_refined_from_kb">LLM Refined (KB)</option>
                    <option value="manual">Manual/Static</option> {/* For non-LLM findings */}
                </select>
            </div>
            <ul style={{ listStyleType: 'none', padding: 0 }}>
                {filteredIssues.map(issue => (
                    <li key={issue.id} style={{ border: '1px solid #eee', padding: '10px', marginBottom: '10px' }}>
                        <h3>{issue.type} - {issue.severity}</h3>
                        <p><strong>URL:</strong> {issue.url}</p>
                        {issue.inputName && issue.inputName !== "N/A" && <p><strong>Input:</strong> {issue.inputName}</p>}
                        <p><strong>Details:</strong> {issue.details}</p>
                        {issue.payloadUsed && issue.payloadUsed !== "N/A" && <p><strong>Payload:</strong> <code>{issue.payloadUsed}</code></p>}
                        {issue.evidence && <p><strong>Evidence:</strong> <pre style={{whiteSpace: 'pre-wrap', backgroundColor: '#f5f5f5', padding: '5px'}}>{issue.evidence}</pre></p>}
                        
                        {issue.llmPayloadSource && (
                            <div className="llm-details" style={{marginTop: '5px', borderTop: '1px dashed #ddd', paddingTop: '5px'}}>
                                <h4>LLM & KB Insights:</h4>
                                <p><strong>Source:</strong> {issue.llmPayloadSource}</p>
                                {issue.llmConfidence && <p><strong>Confidence:</strong> {issue.llmConfidence.toFixed(2)}</p>}
                                {issue.llmPrompt && <p><strong>LLM Prompt (summary):</strong> <pre style={{whiteSpace: 'pre-wrap', backgroundColor: '#f9f9f9', padding: '3px', maxHeight: '50px', overflowY: 'auto'}}>{issue.llmPrompt.substring(0,100)}...</pre></p>}
                                {issue.llmRawResponse && <p><strong>LLM Response (summary):</strong> <pre style={{whiteSpace: 'pre-wrap', backgroundColor: '#f9f9f9', padding: '3px', maxHeight: '50px', overflowY: 'auto'}}>{issue.llmRawResponse.substring(0,100)}...</pre></p>}
                                {issue.contributingKBArticles && issue.contributingKBArticles.length > 0 && (
                                    <p><strong>Contributing KB Articles:</strong> {issue.contributingKBArticles.join(', ')}</p>
                                )}
                                {/* TODO: Add button to show full LLM prompt/response in a modal */}
                            </div>
                        )}
                    </li>
                ))}
            </ul>
            {filteredIssues.length === 0 && <p>No issues match the current filters.</p>}
        </div>
    );
}

export default ScanResults;
