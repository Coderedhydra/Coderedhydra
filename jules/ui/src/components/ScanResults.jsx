// jules/ui/src/components/ScanResults.jsx
import React, { useState } from 'react';

function ScanResults() {
  // TODO: Fetch and display scan results from the core/backend
  const [filter, setFilter] = useState('all'); // 'all', 'high', 'medium', 'low'
  const issues = [
    { id: 1, severity: 'High', type: 'SQLi', url: 'https://example.com/search?q=1', details: '...' },
    { id: 2, severity: 'Medium', type: 'XSS', url: 'https://example.com/profile', details: '...' },
    { id: 3, severity: 'Low', type: 'Missing Header', url: 'https://example.com', details: '...' },
  ];

  const filteredIssues = issues.filter(issue => 
    filter === 'all' || issue.severity.toLowerCase() === filter
  );

  return (
    <div className="scan-results-panel">
      <h2>Scan Results</h2>
      <div>
        Filter by severity:
        <select onChange={(e) => setFilter(e.target.value)} value={filter}>
          <option value="all">All</option>
          <option value="high">High</option>
          <option value="medium">Medium</option>
          <option value="low">Low</option>
        </select>
      </div>
      <ul>
        {filteredIssues.map(issue => (
          <li key={issue.id}>
            <strong>{issue.severity}</strong>: {issue.type} at {issue.url}
            {/* TODO: Add more details and allow expansion */}
          </li>
        ))}
      </ul>
    </div>
  );
}

export default ScanResults;
