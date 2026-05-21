import { Routes, Route } from 'react-router-dom';
import { AdminRouter } from './admin/router';
import { OpenDesignApp } from './pages/OpenDesignApp';

function App() {
  return (
    <div className="App">
      <Routes>
        <Route path="/admin/*" element={<AdminRouter />} />
        <Route path="/*" element={<OpenDesignApp />} />
      </Routes>
    </div>
  );
}

export default App;
