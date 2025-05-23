// jules/ui/src/components/SiteMap.jsx
import React from 'react';

function SiteMap() {
  // TODO: Fetch and display sitemap data from the core/backend
  const discoveredUrls = [
    { id: 1, url: 'https://example.com', method: 'GET' },
    { id: 2, url: 'https://example.com/about', method: 'GET' },
    { id: 3, url: 'https://example.com/contact', method: 'POST' },
  ];

  return (
    <div className="sitemap-panel">
      <h2>Site Map</h2>
      <ul>
        {discoveredUrls.map(item => (
          <li key={item.id}>{item.method} - {item.url}</li>
        ))}
      </ul>
      {/* TODO: Add tree view for sitemap */}
    </div>
  );
}

export default SiteMap;
