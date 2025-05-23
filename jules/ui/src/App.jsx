// jules/ui/src/App.jsx
import React from 'react';
import './App.css'; // Default Vite CSS, can be modified
import SiteMap from './components/SiteMap';
import RepeaterPanel from './components/RepeaterPanel';
import ScanResults from './components/ScanResults';

function App() {
  return (
    <div className="App">
      <header className="App-header">
        <h1>Jules Web Security Framework</h1>
      </header>
      <main>
        {/* Basic layout, can be improved with routing and tabs */}
        <SiteMap />
        <RepeaterPanel />
        <ScanResults />
      </main>
    </div>
  );
}

export default App;
