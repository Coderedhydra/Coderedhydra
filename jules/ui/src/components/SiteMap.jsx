// jules/ui/src/components/SiteMap.jsx
import React from 'react';

function SiteMap() {
  // TODO: Fetch and display sitemap data from the core/backend
  // TODO: Consider adding visual cues for pages with forms, inputs analyzed by LLM, or KB insights.
  const discoveredUrls = [
    { id: 1, url: 'https://example.com', method: 'GET', notes: "Has login form" },
    { id: 2, url: 'https://example.com/about', method: 'GET' },
    { id: 3, url: 'https://example.com/contact', method: 'POST', notes: "Contact form, LLM analyzed" },
  ];

  return (
    <div className="sitemap-panel" style={{ border: '1px solid #ccc', padding: '10px', margin: '10px' }}>
      <h2>Site Map</h2>
      <ul style={{ listStyleType: 'none', padding: 0 }}>
        {discoveredUrls.map(item => (
          <li key={item.id}>{item.method} - {item.url} {item.notes ? `(${item.notes})` : ''}</li>
        ))}
      </ul>
      {/* TODO: Add tree view for sitemap */}
    </div>
  );
}

export default SiteMap;
