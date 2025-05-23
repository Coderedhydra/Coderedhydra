// jules/ui/src/components/RepeaterPanel.jsx
import React, { useState } from 'react';

function RepeaterPanel() {
  const [request, setRequest] = useState('GET / HTTP/1.1\nHost: example.com\n...');
  const [response, setResponse] = useState('');

  const handleSendRequest = () => {
    // TODO: Implement sending the request through the core/backend
    // For now, simulate a response
    setResponse('HTTP/1.1 200 OK\nContent-Type: text/html\n\n<html><body>Hello!</body></html>');
  };

  return (
    <div className="repeater-panel">
      <h2>Repeater</h2>
      <textarea
        value={request}
        onChange={(e) => setRequest(e.target.value)}
        rows={10}
        cols={80}
      />
      <button onClick={handleSendRequest}>Send</button>
      <h3>Response</h3>
      <pre>{response}</pre>
    </div>
  );
}

export default RepeaterPanel;
